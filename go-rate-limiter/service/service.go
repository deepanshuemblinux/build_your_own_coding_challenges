package service

import (
	"github.com/deepanshuemblinux/go-rate-limiter/types"
)

type MessageService interface {
	GetMessage(message string) *types.APIResponse
}

type textMessageService struct{}

func NewTextMessageService() *textMessageService {
	return &textMessageService{}
}

func (s *textMessageService) GetMessage(message string) *types.APIResponse {
	return &types.APIResponse{
		Message: message,
	}
}
