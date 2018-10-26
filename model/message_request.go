package model

// MessageRequest represents MessageRequest table
type MessageRequest struct {
	RecipientIdentifier string
	RequestID           string
}

// MessageRequestList represents a list of MessageRequest items
type MessageRequestList []*MessageRequest
