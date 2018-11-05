package dynamodbdriver

// DriverConfig represents DynamoDbDriver configuration
type DriverConfig struct {
	_ struct{}

	// MessageReqyestTable represents Amazon DynamoDB table for
	// model.MessageRequest
	MessageRequestTable *string
}
