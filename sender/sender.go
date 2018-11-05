package sender

import "github.com/petrulis/abn-amro-assignment/model"

// Sender interface represents MessageRequest delivery channel.
type Sender interface {
	Send(request *model.MessageRequest) error
}
