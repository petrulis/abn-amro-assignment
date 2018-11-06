package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/petrulis/abn-amro-assignment/model"
)

// Response represents operation response
type response struct {
	// Items represents a list of MessageRequestItems
	Items model.MessageRequestList

	// NextToken represents base64 encoded last evaluated key.
	NextToken *string

	// Count represents number of items in a response.
	Count int
}

// newResponse creates new response instance with items, count and nextToken.
func newResponse(list model.MessageRequestList, lastEvaluatedKey *model.Key) *response {
	resp := &response{
		Items: list,
		Count: len(list),
	}
	if lastEvaluatedKey != nil && !lastEvaluatedKey.IsEmpty() {
		nextToken := lastEvaluatedKey.EncodeBase64()
		resp.NextToken = aws.String(nextToken)
	}
	return resp
}

// Marshal implements Marshal interface
func (r *response) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}
