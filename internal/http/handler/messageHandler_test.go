package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vaksi/messaging/internal/services"
	"github.com/vaksi/messaging/pkg/responses"
	"gotest.tools/assert"
)

type MessageServiceFake struct {}

func (msf *MessageServiceFake) CreateMessage(param services.CreateMessageParam) (err error) {
	return
}

func (msf *MessageServiceFake) GetMessagesRealtime(userID uint64, status uint8) (response []services.GetMessagesRealtimeResponse, err error) {
	return
}

func (msf *MessageServiceFake) GetMessagesInbox(userID uint64) (response []services.GetMessagesInboxResponse, err error) {
	return
}

func (msf *MessageServiceFake) ReceiveMessage(key []byte, msg []byte) (err error) {
	return
}

func TestMessageHandler_CreateMessage(t *testing.T) {
	msgSvc := &MessageServiceFake{}
	msgHandler := NewMessageHandler(msgSvc)

	// define router
	router := mux.NewRouter()
	router.HandleFunc("/messages", msgHandler.CreateMessage).Methods("POST")

	t.Run("TestCase1 Create Message is Succes", func(t *testing.T) {
		// payload
		payload := `{"user_id":1, "subject":"okem", "to_user_id":2, "text":"blablabla"}`
		// endpoint
		req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(payload))
		rw := httptest.NewRecorder()
		router.ServeHTTP(rw, req)

		// expectation
		act, _ := json.Marshal(responses.APICreated)

		// Assertion
		assert.Equal(t, rw.Code, http.StatusCreated)
		assert.Equal(t, rw.Body.String(), string(act))
	})
}
