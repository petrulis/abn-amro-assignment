package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/petrulis/abn-amro-assignment/model"
)

// Response represents operation response
type response struct {
	Items     model.MessageRequestList
	NextToken *string
	Count     int
}

func newResponse(list model.MessageRequestList, lastEvaluatedKey *model.Key) *response {
	resp := &response{
		Items: list,
		Count: len(list),
	}
	if lastEvaluatedKey != nil {
		nextToken := lastEvaluatedKey.EncodeBase64()
		resp.NextToken = aws.String(nextToken)
	}
	return resp
}

// Marshal implements Marshaler
func (r *response) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

// String implement Stringer
func (r *response) String() string {
	return string(r.Marshal())
}
