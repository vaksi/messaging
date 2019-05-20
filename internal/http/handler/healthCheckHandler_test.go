package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/vaksi/messaging/pkg/responses"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
)

func TestGetHealthCheck(t *testing.T) {
	// define router
	router := mux.NewRouter()
	router.HandleFunc("/health-check", GetHealthCheck).Methods("GET")
	// endpoint
	w := httptest.NewRequest("GET", "/health-check", nil)
	req := httptest.NewRecorder()
	router.ServeHTTP(req, w)

	// expectation
	act, _ := json.Marshal(responses.APIOK)

	// Assertion
	assert.Equal(t, http.StatusOK, req.Code)
	assert.Equal(t, string(act), req.Body.String())
}
