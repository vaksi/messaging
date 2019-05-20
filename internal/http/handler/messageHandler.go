package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/vaksi/messaging/internal/services"
	"github.com/vaksi/messaging/pkg/logrus"
	"github.com/vaksi/messaging/pkg/responses"
)

type MessageHandler struct {
	messageSvc services.IMessageService
}

type IMessageHandler interface {
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetMessageRealtime(w http.ResponseWriter, r *http.Request)
	GetMessageInbox(w http.ResponseWriter, r *http.Request)
}

func NewMessageHandler(msgSvc services.IMessageService) *MessageHandler{
	return &MessageHandler{
		messageSvc: msgSvc,
	}
}

type createMessageParam struct {
	UserID   uint   `json:"user_id" valid:"required"`
	Subject  string `json:"subject" valid:"required"`
	Text     string `json:"text" valid:"required"`
	ToUserID uint   `json:"to_user_id"  valid:"required"`
}

func (mh *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var body createMessageParam
	// decode body params to struct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Error(err)
		responses.Write(w, responses.APIErrorValidation.WithMessage(err))
		return
	}
	isValid, errValidation := govalidator.ValidateStruct(body)
	if !isValid {
		logrus.Debug(errValidation)
		responses.Write(w, responses.APIErrorValidation.WithMessage(errValidation.Error()))
	}
	// save data to service
	err := mh.messageSvc.CreateMessage(services.CreateMessageParam{
		UserID:   uint64(body.UserID),
		Subject:  body.Subject,
		ToUserID: uint64(body.ToUserID),
		Text:     body.Text,
	})
	if err != nil {
		logrus.Error(err)
		responses.Write(w, responses.APIErrorUnknown)
		return
	}

	responses.Write(w, responses.APICreated)
}

func (mh *MessageHandler) GetMessageRealtime(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	status, _ := strconv.ParseUint(r.URL.Query().Get("status"), 10, 8)
	// save data to service
	resp, err := mh.messageSvc.GetMessagesRealtime(userID,uint8(status))
	if err != nil {
		logrus.Error(err)
		responses.Write(w, responses.APIErrorUnknown)
		return
	}

	responses.Write(w, responses.APIOK.WithData(resp))
}

func (mh *MessageHandler) GetMessageInbox(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	// save data to service
	resp, err := mh.messageSvc.GetMessagesInbox(userID)
	if err != nil {
		logrus.Error(err)
		responses.Write(w, responses.APIErrorUnknown)
		return
	}

	responses.Write(w, responses.APIOK.WithData(resp))
}