package handler

import (
	"github.com/vaksi/messaging/pkg/responses"
	"net/http"
)

// GetHealthCheck this function for get health check
func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := responses.APIOK
	responses.Write(w, resp)
}