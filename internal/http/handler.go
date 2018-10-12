/*  routes.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/10/18 13:26 
 */

package http

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    config "github.com/spf13/viper"
    "github.com/urfave/negroni"
    "github.com/vaksi/user_management/internal/healths"
    "github.com/vaksi/user_management/internal/users/service"
    "github.com/vaksi/user_management/pkg/apiResponse"
)

// Handler define struct for handler
type Handler struct {
    HealthService *healths.HealthService
    UserService   service.UserServiceFactory
}

// PathPrefix defined for name spaces
const PathPrefix = "/user-management/v1"

type userParams struct {
    Email    string `json:"email"`
    FullName string `json:"full_name"`
    Gender   int8   `json:"gender"`
}

func (h *Handler) healthCheck(res http.ResponseWriter, req *http.Request) {
    resp := apiResponse.APIOk
    resp.Message = h.HealthService.HealthCheckStatus()
    apiResponse.Write(res, resp)
}

func (h *Handler) createUser(res http.ResponseWriter, req *http.Request) {
    var (
        body userParams
    )

    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        apiResponse.Write(res, apiResponse.APIErrorInvalidData)
        return
    }

    if err := h.UserService.CreateUser(body.Email, body.FullName, body.Gender); err != nil {
        apiResponse.Write(res, apiResponse.APIErrorUnknown)
        return
    }
    apiResponse.Write(res, apiResponse.APICreated)
}

// MakeHandler handler for internal route
func (h *Handler) MakeHandler() http.Handler {
    // define route
    router := mux.NewRouter().StrictSlash(false)
    route := router.PathPrefix(PathPrefix).Subrouter()

    // Health
    route.HandleFunc("/health-check", h.healthCheck).Methods("GET")

    // User
    route.HandleFunc("/users", h.createUser).Methods("POST")

    n := negroni.Classic()
    recovery := negroni.NewRecovery() // Panic handler
    if config.GetBool("app.debug") == false {
        recovery.PrintStack = false
    }
    n.Use(recovery)
    n.UseHandler(router)
    return n
}
