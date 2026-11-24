package dynamodbclaim

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestLocalConditionalClaim(t *testing.T) {
	if os.Getenv("P1_DYNAMODB_CLAIM_TEST") != "1" {
		t.Skip("set P1_DYNAMODB_CLAIM_TEST=1 to use LocalStack on 127.0.0.1:4566")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := newClaimFixture(t, ctx, "p1-day1-claim")
	table := client.table

	for _, fixture := range []struct {
		id      string
		version string
	}{
		{id: "accepted", version: "1"},
		{id: "stale", version: "2"},
		{id: "ambiguous", version: "1"},
	} {
		item := readyItem(fixture.id, fixture.version)
		if _, err := client.dynamo.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(table), Item: item}); err != nil {
			t.Fatalf("seed %q: %v", fixture.id, err)
		}
		assertItem(t, ctx, client, "before", fixture.id, item)
	}

	outcome, err := client.Claim(ctx, "accepted", "worker-A", "attempt-A")
	if outcome != Claimed || err != nil {
		t.Fatalf("first claim = %s, %v", outcome, err)
	}
	accepted := runningItem("accepted", "worker-A", "attempt-A")
	assertItem(t, ctx, client, string(outcome), "accepted", accepted)

	for _, rejected := range []struct {
		id   string
		want map[string]types.AttributeValue
	}{
		{id: "accepted", want: accepted},
		{id: "stale", want: readyItem("stale", "2")},
		{id: "missing"},
	} {
		outcome, err = client.Claim(ctx, rejected.id, "worker-B", "attempt-B")
		if outcome != NotClaimed || err != nil {
			t.Fatalf("rejected claim for %q = %s, %v", rejected.id, outcome, err)
		}
		assertItem(t, ctx, client, string(outcome), rejected.id, rejected.want)
	}

	outcome, err = NewLocal(table+"-missing").Claim(ctx, "accepted", "worker-B", "attempt-B")
	var missingTable *types.ResourceNotFoundException
	if outcome != Unknown || !errors.As(err, &missingTable) {
		t.Fatalf("missing table claim = %s, %v", outcome, err)
	}
	t.Logf("missing table outcome=%s error=%v", outcome, err)

	transport := &http.Transport{}
	defer transport.CloseIdleConnections()
	lostResponse := errors.New("fixture dropped successful update response")
	requests := 0
	faulty := newLocal(table, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests++
		response, err := transport.RoundTrip(request)
		if err != nil || response.StatusCode != http.StatusOK {
			return response, err
		}
		response.Body.Close()
		return nil, lostResponse
	}))
	outcome, err = faulty.Claim(ctx, "ambiguous", "worker-C", "attempt-C")
	if outcome != Unknown || !errors.Is(err, lostResponse) || requests != 1 {
		t.Fatalf("lost response claim = %s, %v, requests=%d", outcome, err, requests)
	}
	t.Logf("lost response outcome=%s requests=%d error=%v", outcome, requests, err)
	assertItem(t, ctx, client, "after lost response", "ambiguous", runningItem("ambiguous", "worker-C", "attempt-C"))
}

type claimProcessResult struct {
	PID, ParentPID   int
	Owner, AttemptID string
	Outcome          Outcome
	Error            string
}

func TestLocalClaimProcesses(t *testing.T) {
	if os.Getenv("P1_DYNAMODB_CLAIM_TEST") != "1" {
		t.Skip("set P1_DYNAMODB_CLAIM_TEST=1 to use LocalStack on 127.0.0.1:4566")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := newClaimFixture(t, ctx, "p1-day2-claim")
	const workflowID = "shared"
	item := readyItem(workflowID, "1")
	if _, err := client.dynamo.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(client.table), Item: item}); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}
	assertItem(t, ctx, client, "before contention", workflowID, item)

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	workDir := t.TempDir()
	startRead, startWrite := claimPipe(t)
	readyRead, readyWrite := claimPipe(t)
	deadline, _ := ctx.Deadline()
	if err := readyRead.SetReadDeadline(deadline); err != nil {
		t.Fatal(err)
	}

	const workerCount = 10
	type worker struct {
		owner, attemptID string
		command          *exec.Cmd
		stdout, stderr   bytes.Buffer
	}
	var workers [workerCount]worker
	defer func() {
		cancel()
		for i := range workers {
			if cmd := workers[i].command; cmd != nil && cmd.Process != nil && cmd.ProcessState == nil {
				_ = cmd.Wait()
			}
		}
	}()
	for i := range workers {
		worker := &workers[i]
		worker.owner = fmt.Sprintf("worker-%02d", i+1)
		worker.attemptID = fmt.Sprintf("attempt-%02d", i+1)
		cmd := exec.CommandContext(ctx, executable, "-test.run=^TestLocalClaimWorker$", "--",
			client.table, workflowID, worker.owner, worker.attemptID, workDir)
		cmd.Env = []string{"P1_DYNAMODB_CLAIM_WORKER=1"}
		cmd.Stdin, cmd.Stdout, cmd.Stderr = startRead, &worker.stdout, &worker.stderr
		cmd.ExtraFiles = []*os.File{readyWrite} // Readiness pipe is descriptor 3 in the child.
		worker.command = cmd
		if err := cmd.Start(); err != nil {
			t.Fatalf("start %s: %v", worker.owner, err)
		}
	}
	_ = readyWrite.Close()
	ready := make([]byte, workerCount)
	if n, err := io.ReadFull(readyRead, ready); err != nil {
		t.Fatalf("%d/%d children reached the gate: %v", n, workerCount, err)
	}
	t.Logf("parent_pid=%d releasing %d ready children", os.Getpid(), len(ready))
	if err := startWrite.Close(); err != nil {
		t.Fatalf("release children: %v", err)
	}

	var winner claimProcessResult
	counts := make(map[Outcome]int)
	pids := make(map[int]bool)
	for i := range workers {
		worker := &workers[i]
		cmd := worker.command
		if err := cmd.Wait(); err != nil {
			t.Errorf("wait %s pid=%d: %v\n%s%s", worker.owner, cmd.Process.Pid, err, &worker.stdout, &worker.stderr)
			continue
		}
		var result claimProcessResult
		if err := json.NewDecoder(&worker.stdout).Decode(&result); err != nil {
			t.Errorf("decode %s: %v", worker.owner, err)
			continue
		}
		t.Logf("pid=%d parent_pid=%d owner=%s attempt_id=%s outcome=%s error=%q", result.PID, result.ParentPID, result.Owner, result.AttemptID, result.Outcome, result.Error)
		if result.PID != cmd.Process.Pid || result.ParentPID != os.Getpid() || result.PID == os.Getpid() || pids[result.PID] ||
			result.Owner != worker.owner || result.AttemptID != worker.attemptID {
			t.Errorf("invalid process identity for %s: %+v", worker.owner, result)
			continue
		}
		pids[result.PID] = true
		counts[result.Outcome]++
		if result.Error != "" {
			t.Errorf("%s claim error: %s", result.Owner, result.Error)
		}
		if result.Outcome == Claimed {
			winner = result
		}
	}
	t.Logf("results=%d CLAIMED=%d NOT_CLAIMED=%d ERROR/UNKNOWN=%d", len(pids), counts[Claimed], counts[NotClaimed], counts[Unknown])
	if len(pids) != workerCount || counts[Claimed] != 1 || counts[NotClaimed] != workerCount-1 || counts[Unknown] != 0 {
		t.Fatalf("expected ten accounted results and exactly one confirmed winner: %v", counts)
	}
	markers, err := os.ReadDir(workDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(markers) != 1 || markers[0].Name() != winner.Owner {
		t.Fatalf("fixture work markers = %v, want only %s", markers, winner.Owner)
	}
	attempt, err := os.ReadFile(filepath.Join(workDir, winner.Owner))
	if err != nil || string(attempt) != winner.AttemptID {
		t.Fatalf("fixture work attempt = %q, %v, want %q", attempt, err, winner.AttemptID)
	}
	t.Logf("only fixture work marker: owner=%s attempt_id=%s", winner.Owner, attempt)
	assertItem(t, ctx, client, "after contention", workflowID, runningItem(workflowID, winner.Owner, winner.AttemptID))
}

func TestLocalClaimWorker(t *testing.T) {
	if os.Getenv("P1_DYNAMODB_CLAIM_WORKER") != "1" {
		t.Skip("subprocess entry point for TestLocalClaimProcesses")
	}
	args := os.Args
	if len(args) != 8 || args[2] != "--" {
		t.Fatal("expected fixture table, workflow ID, owner, attempt ID and work directory")
	}
	table, workflowID, owner, attemptID, workDir := args[3], args[4], args[5], args[6], args[7]
	client := NewLocal(table)
	ready := os.NewFile(3, "claim-ready")
	if _, err := ready.Write([]byte{1}); err != nil {
		t.Fatalf("signal readiness: %v", err)
	}
	_ = ready.Close()
	if _, err := io.Copy(io.Discard, os.Stdin); err != nil {
		t.Fatalf("wait for start signal: %v", err)
	}

	outcome, err := client.Claim(context.Background(), workflowID, owner, attemptID)
	result := claimProcessResult{PID: os.Getpid(), ParentPID: os.Getppid(), Owner: owner, AttemptID: attemptID, Outcome: outcome}
	if err != nil {
		result.Error = err.Error()
	}
	// The marker is an observable fixture action, independent of the reported outcome.
	if outcome == Claimed && err == nil {
		if err := os.WriteFile(filepath.Join(workDir, owner), []byte(attemptID), 0600); err != nil {
			t.Fatalf("execute fixture work: %v", err)
		}
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		t.Fatalf("write claim result: %v", err)
	}
}

func claimPipe(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = read.Close()
		_ = write.Close()
	})
	return read, write
}

func newClaimFixture(t *testing.T, ctx context.Context, prefix string) *Client {
	t.Helper()
	table := fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	client := NewLocal(table)
	t.Logf("table=%s endpoint=http://127.0.0.1:4566 region=us-east-1", table)
	_, err := client.dynamo.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("workflow_id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{{AttributeName: aws.String("workflow_id"), KeyType: types.KeyTypeHash}},
	})
	if err != nil {
		t.Fatalf("create fixture table: %v", err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := client.dynamo.DeleteTable(cleanup, &dynamodb.DeleteTableInput{TableName: aws.String(table)}); err != nil {
			t.Errorf("delete fixture table %s: %v", table, err)
		}
	})
	waiter := dynamodb.NewTableExistsWaiter(client.dynamo)
	if err := waiter.Wait(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(table)}, 10*time.Second); err != nil {
		t.Fatalf("wait for fixture table: %v", err)
	}
	return client
}

func readyItem(id, version string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"workflow_id": &types.AttributeValueMemberS{Value: id},
		"status":      &types.AttributeValueMemberS{Value: "READY"},
		"version":     &types.AttributeValueMemberN{Value: version},
		"attempts":    &types.AttributeValueMemberN{Value: "0"},
	}
}

func runningItem(id, owner, attemptID string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"workflow_id": &types.AttributeValueMemberS{Value: id},
		"status":      &types.AttributeValueMemberS{Value: "RUNNING"},
		"version":     &types.AttributeValueMemberN{Value: "2"},
		"attempts":    &types.AttributeValueMemberN{Value: "1"},
		"owner":       &types.AttributeValueMemberS{Value: owner},
		"attempt_id":  &types.AttributeValueMemberS{Value: attemptID},
	}
}

func assertItem(t *testing.T, ctx context.Context, client *Client, phase, id string, want map[string]types.AttributeValue) {
	t.Helper()
	result, err := client.dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      aws.String(client.table),
		ConsistentRead: aws.Bool(true),
		Key:            map[string]types.AttributeValue{"workflow_id": &types.AttributeValueMemberS{Value: id}},
	})
	if err != nil {
		t.Fatalf("get item %q: %v", id, err)
	}
	gotJSON, wantJSON := itemJSON(result.Item), itemJSON(want)
	if !reflect.DeepEqual(result.Item, want) && !(len(result.Item) == 0 && len(want) == 0) {
		t.Fatalf("%s item %q = %s, want %s", phase, id, gotJSON, wantJSON)
	}
	t.Logf("%s item %q = %s", phase, id, gotJSON)
}

func itemJSON(item map[string]types.AttributeValue) string {
	if len(item) == 0 {
		return "<absent>"
	}
	values := make(map[string]map[string]string, len(item))
	for name, value := range item {
		switch attribute := value.(type) {
		case *types.AttributeValueMemberS:
			values[name] = map[string]string{"S": attribute.Value}
		case *types.AttributeValueMemberN:
			values[name] = map[string]string{"N": attribute.Value}
		default:
			values[name] = map[string]string{"unexpected type": fmt.Sprintf("%T", value)}
		}
	}
	encoded, _ := json.Marshal(values)
	return string(encoded)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
