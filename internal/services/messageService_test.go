package services

import (
	"fmt"
	"testing"

	"github.com/vaksi/messaging/internal/models"
	"github.com/vaksi/messaging/internal/repositories"
	"gotest.tools/assert"
)

const (
	createMessageStoreMessageError = 1
	errorStoreMessage = "error Store Message"

	createMessagePublishMessageError = 2
	errorPublishMessage = "error Publish Message"

	createMessagePublishMessageErrorButErrorUpdateMessage = 3
	errorUpdateMessageOnErrorPublishMessage = "error Update On Error Publish Message"

	createMessageUpdateMessageErrorOnPublishSuccess = 4
	errorMessageUpdateMessageErrorOnPublishSuccess = "error Update On Publish Message Success"

	createMessageHasSent = 5
)

type MessageRepositoryFake struct {
}

func (mr *MessageRepositoryFake) GetMessageByUserIDStatusSent(userID uint64, status uint8) (messages []models.Message, err error) {
	return
}

func (mr *MessageRepositoryFake) GetMessageByUserIDInbox(userID uint64) (messages []models.Message, err error) {
	return
}

func (mr *MessageRepositoryFake) Store(msg *models.Message) (lastID uint64, err error) {
	switch msg.UserID {
	case createMessageStoreMessageError:
		return 0, fmt.Errorf(errorStoreMessage)
	case createMessagePublishMessageError:
		return createMessagePublishMessageError, nil
	case createMessagePublishMessageErrorButErrorUpdateMessage:
		return createMessagePublishMessageErrorButErrorUpdateMessage, nil
	case createMessageUpdateMessageErrorOnPublishSuccess:
		return createMessageUpdateMessageErrorOnPublishSuccess, nil
	case createMessageHasSent:
		return createMessageHasSent,nil
	}
	return
}
func (mr *MessageRepositoryFake) Publish(req repositories.PublishRequest) (err error) {
	switch req.MessageID {
	case createMessagePublishMessageError, createMessagePublishMessageErrorButErrorUpdateMessage:
		return fmt.Errorf(errorPublishMessage)
	}
	return
}
func (mr *MessageRepositoryFake) Update(msg *models.Message) (err error) {
	switch msg.ID {
	case createMessagePublishMessageErrorButErrorUpdateMessage:
		return fmt.Errorf(errorUpdateMessageOnErrorPublishMessage)
	case createMessageUpdateMessageErrorOnPublishSuccess:
		return fmt.Errorf(errorUpdateMessageOnErrorPublishMessage)
	}
	return
}

func TestMessageService_CreateMessage(t *testing.T) {
	messageRepo := &MessageRepositoryFake{}
	msgSvc := NewMessageService(messageRepo)
	t.Run("#TestCase1 Store of message should be error", func(t *testing.T) {
		param := CreateMessageParam{
			UserID:   1,
			Text:     "bla bla",
			Subject:  "okem",
			ToUserID: 2,
		}
		act := msgSvc.CreateMessage(param)
		assert.Error(t, act, errorStoreMessage)
	})
	t.Run("#TestCase2 Publish of message should be error", func(t *testing.T) {
		param := CreateMessageParam{
			UserID:   2,
			Text:     "bla bla",
			Subject:  "okem",
			ToUserID: 3,
		}
		act := msgSvc.CreateMessage(param)
		assert.NilError(t, act, errorPublishMessage)
	})
	t.Run("#TestCase3 Update should Be Error on Publish Message is error", func(t *testing.T) {
		param := CreateMessageParam{
			UserID:   3,
			Text:     "bla bla",
			Subject:  "okem",
			ToUserID: 4,
		}
		act := msgSvc.CreateMessage(param)
		assert.NilError(t, act, errorUpdateMessageOnErrorPublishMessage)
	})
	t.Run("#TestCase4 Update should Be Error on Publish Message is success", func(t *testing.T) {
		param := CreateMessageParam{
			UserID:   4,
			Text:     "bla bla",
			Subject:  "okem",
			ToUserID: 5,
		}
		act := msgSvc.CreateMessage(param)
		assert.NilError(t, act, errorMessageUpdateMessageErrorOnPublishSuccess)
	})
	t.Run("#TestCase5 Create Message has sent", func(t *testing.T) {
		param := CreateMessageParam{
			UserID:   5,
			Text:     "bla bla",
			Subject:  "okem",
			ToUserID: 6,
		}
		act := msgSvc.CreateMessage(param)
		assert.NilError(t, act, errorUpdateMessageOnErrorPublishMessage)
	})
}
