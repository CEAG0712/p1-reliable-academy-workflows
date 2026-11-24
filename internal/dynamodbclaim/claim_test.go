package dynamodbclaim

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	table := fmt.Sprintf("p1-day1-claim-%d", time.Now().UnixNano())
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
