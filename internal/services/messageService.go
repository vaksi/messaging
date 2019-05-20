package services

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vaksi/messaging/internal/models"
	"github.com/vaksi/messaging/internal/repositories"
)

// MessageService of struct message
type MessageService struct {
	messageRepo repositories.IMessageRepository
}

// IMessageService of interface message service
type IMessageService interface {
	CreateMessage(CreateMessageParam) error
	GetMessagesRealtime(userID uint64, status uint8) (response []GetMessagesRealtimeResponse, err error)
	GetMessagesInbox(userID uint64) (response []GetMessagesInboxResponse, err error)
	ReceiveMessage(key []byte, msg []byte) error
}

func NewMessageService(messageRepo repositories.IMessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

// Message Status
const (
	MessageStatusSending     = 1
	MessageStatusSent        = 2
	MessageStatusSendFailure = 3
	MessageStatusReceived    = 4
)

var (
	MessageStatusSendingStr     = "sending"
	MessageStatusSentStr        = "sent_out"
	MessageStatusSendFailureStr = "failure"
	MessageStatusReceivedStr    = "received"
)

func statusMessageToString(statusID int8) string {
	switch statusID {
	case MessageStatusSending:
		return MessageStatusSendingStr
	case MessageStatusSent:
		return MessageStatusSentStr
	case MessageStatusSendFailure:
		return MessageStatusSendFailureStr
	case MessageStatusReceived:
		return MessageStatusReceivedStr
	}
	return ""
}

// CreateMessageParam for parameter CreateMessage
type CreateMessageParam struct {
	UserID   uint64
	Subject  string
	Text     string
	ToUserID uint64
}

// CreateMessage this service for created message and send
func (ms *MessageService) CreateMessage(param CreateMessageParam) (err error) {
	// store message
	messageID, err := ms.messageRepo.Store(&models.Message{
		UserID:   param.UserID,
		ToUserID: param.ToUserID,
		Subject:  param.Subject,
		Text:     param.Text,
		Status:   MessageStatusSending,
	})
	if err != nil {
		return err
	}

	// send message
	go ms.sendMessage(messageID)
	return
}

func (ms *MessageService) sendMessage(messageID uint64) {
	// publish message to kafka
	errSend := ms.messageRepo.Publish(repositories.PublishRequest{MessageID: messageID})
	if errSend != nil {
		logrus.Error(errSend)
		err := ms.messageRepo.Update(&models.Message{ID: messageID, Status: MessageStatusSendFailure})
		if err != nil {
			return
		}
		return
	}
	err := ms.messageRepo.Update(&models.Message{ID: messageID, Status: MessageStatusSent})
	if err != nil {
		return
	}
	logrus.Info("Message Has Been Sent")
	return
}

type GetMessagesRealtimeResponse struct {
	MessageID   uint64 `json:"message_id"`
	ToUserID    uint64 `json:"to_user_id"`
	Subject     string `json:"subject"`
	Text        string `json:"text"`
	Status      string `json:"status"`
	CreatedDate string `json:"created_date"`
	LastUpdated string `json:"last_updated"`
}

func (ms *MessageService) GetMessagesRealtime(userID uint64, status uint8) (response []GetMessagesRealtimeResponse, err error) {
	messages, err := ms.messageRepo.GetMessageByUserIDStatusSent(userID, status)
	if err != err {
		logrus.Error(err)
	}
	for _, message := range messages {
		response = append(response, GetMessagesRealtimeResponse{
			MessageID:   message.ID,
			ToUserID:    message.ToUserID,
			Subject: message.Subject,
			Text:        message.Text,
			Status:      statusMessageToString(message.Status),
			CreatedDate: message.CreatedAt.String(),
			LastUpdated: message.UpdatedAt.String(),
		})
	}
	return
}

type GetMessagesInboxResponse struct {
	MessageID   uint64 `json:"message_id"`
	FromUserID  uint64 `json:"from_user_id"`
	Subject     string `json:"subject"`
	Text        string `json:"text"`
	Status      string `json:"status"`
	CreatedDate string `json:"created_date"`
	LastUpdated string `json:"last_updated"`
}

func (ms *MessageService) GetMessagesInbox(userID uint64) (response []GetMessagesInboxResponse, err error) {
	messages, err := ms.messageRepo.GetMessageByUserIDInbox(userID)
	if err != err {
		logrus.Error(err)
	}
	for _, message := range messages {
		response = append(response, GetMessagesInboxResponse{
			MessageID:   message.ID,
			FromUserID:  message.UserID,
			Subject: message.Subject,
			Text:        message.Text,
			Status:      statusMessageToString(message.Status),
			CreatedDate: message.CreatedAt.String(),
			LastUpdated: message.UpdatedAt.String(),
		})
	}
	return
}

type SendMessagePayload struct {
	MessageID uint64 `json:"message_id"`
}
func (ms *MessageService) ReceiveMessage(key []byte, msg []byte) (err error) {
	time.Sleep(5 * time.Second) // special case
	// Get payload
	logrus.Info("receive message")
	payload := SendMessagePayload{}
	err = json.Unmarshal(msg, &payload)
	if err != nil {
		return err
	}
	// update status message
	err = ms.messageRepo.Update(&models.Message{ID: payload.MessageID, Status:MessageStatusReceived})

	return
}