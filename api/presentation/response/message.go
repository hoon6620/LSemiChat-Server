package response

import (
	"app/api/domain/entity"
	"time"
)

type MessageResponse struct {
	ID        string        `json:"id"`
	Message   string        `json:"message"`
	Grade     int           `json:"grade"`
	CreatedAt *time.Time    `json:"created_at"`
	Author    *UserResponse `json:"author"`
}

type MessagesResponse struct {
	Messages []*MessageResponse `json:"messages"`
}

func ConvertToMessageResponse(msg *entity.Message) *MessageResponse {
	return &MessageResponse{
		ID:        msg.ID,
		Message:   msg.Message,
		Grade:     msg.Grade,
		CreatedAt: msg.CreatedAt,
		Author:    ConvertToUserResponse(msg.Author),
	}
}

func ConvertToMessagesResponse(messages []*entity.Message) *MessagesResponse {
	result := make([]*MessageResponse, 0, len(messages))
	for _, message := range messages {
		result = append(result, ConvertToMessageResponse(message))
	}
	return &MessagesResponse{
		Messages: result,
	}
}
