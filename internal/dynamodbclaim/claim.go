// Package dynamodbclaim is a loopback-only DynamoDB conditional claim experiment.
package dynamodbclaim

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Outcome string

const (
	Claimed    Outcome = "CLAIMED"
	NotClaimed Outcome = "NOT_CLAIMED"
	Unknown    Outcome = "ERROR/UNKNOWN"
)

type Client struct {
	dynamo *dynamodb.Client
	table  string
}

func NewLocal(table string) *Client {
	return newLocal(table, &http.Transport{})
}

func newLocal(table string, transport http.RoundTripper) *Client {
	config := aws.Config{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "test", SecretAccessKey: "test"}, nil
		}),
		// Replaying an ambiguous write could turn its success into a conditional rejection.
		Retryer: func() aws.Retryer { return aws.NopRetryer{} },
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	dynamo := dynamodb.NewFromConfig(config, func(options *dynamodb.Options) {
		options.BaseEndpoint = aws.String("http://127.0.0.1:4566")
	})
	return &Client{dynamo: dynamo, table: table}
}

// Claim attempts a single atomic update of an existing READY/version-1 item.
// Only Claimed confirms ownership. Unknown can include a committed update.
func (c *Client) Claim(ctx context.Context, workflowID, owner, attemptID string) (Outcome, error) {
	if workflowID == "" || owner == "" || attemptID == "" {
		return Unknown, errors.New("workflow ID, owner and attempt ID are required")
	}
	result, err := c.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           aws.String(c.table),
		Key:                 map[string]types.AttributeValue{"workflow_id": &types.AttributeValueMemberS{Value: workflowID}},
		ConditionExpression: aws.String("attribute_exists(workflow_id) AND #status = :ready AND #version = :one"),
		UpdateExpression:    aws.String("SET #status = :running, #owner = :owner, attempt_id = :attempt ADD #version :one, attempts :one"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status", "#version": "version", "#owner": "owner",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ready":   &types.AttributeValueMemberS{Value: "READY"},
			":running": &types.AttributeValueMemberS{Value: "RUNNING"},
			":owner":   &types.AttributeValueMemberS{Value: owner},
			":attempt": &types.AttributeValueMemberS{Value: attemptID},
			":one":     &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	var rejected *types.ConditionalCheckFailedException
	switch {
	case errors.As(err, &rejected):
		return NotClaimed, nil
	case err != nil:
		return Unknown, fmt.Errorf("claim workflow %q: %w", workflowID, err)
	}

	attrs := result.Attributes
	attemptCount, countErr := strconv.ParseInt(numberAttribute(attrs["attempts"]), 10, 64)
	if stringAttribute(attrs["workflow_id"]) != workflowID || stringAttribute(attrs["status"]) != "RUNNING" ||
		stringAttribute(attrs["owner"]) != owner || stringAttribute(attrs["attempt_id"]) != attemptID ||
		numberAttribute(attrs["version"]) != "2" || countErr != nil || attemptCount < 1 {
		return Unknown, fmt.Errorf("claim workflow %q: update response did not confirm the claim", workflowID)
	}
	return Claimed, nil
}

func stringAttribute(attribute types.AttributeValue) string {
	if value, ok := attribute.(*types.AttributeValueMemberS); ok {
		return value.Value
	}
	return ""
}

func numberAttribute(attribute types.AttributeValue) string {
	if value, ok := attribute.(*types.AttributeValueMemberN); ok {
		return value.Value
	}
	return ""
}
