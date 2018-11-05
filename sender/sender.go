package sender

import "github.com/petrulis/abn-amro-assignment/model"

type Sender interface {
	Send(request *model.MessageRequest) error
}
