package main

import (
	"encoding/json"
	"github.com/petrulis/abn-amro-assignment/model"
)

// Output represents operation response
type Output struct {
	Items     model.MessageRequestList
	NextToken *string
	Count     *int64
}

// Marshal implements Marshaler
func (o *Output) Marshal() []byte {
	b, _ := json.Marshal(o)
	return b
}

// String implement Stringer
func (o *Output) String() string {
	return string(o.Marshal())
}
