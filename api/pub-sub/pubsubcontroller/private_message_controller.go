package pubsubcontroller

import (
	"github.com/karlosdaniel451/message-chat/usecase"
)

type PrivateMessageController struct {
	taskUseCase usecase.PrivateMessageUseCase
}

func NewPrivateMessageController(taskUseCase usecase.PrivateMessageUseCase) PrivateMessageController {
	return PrivateMessageController{taskUseCase: taskUseCase}
}
