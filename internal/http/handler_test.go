/*  routes.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/10/18 13:26
 */

package http

import (
    "bytes"
    "encoding/json"
    "fmt"
    goHttp "net/http"
    "net/http/httptest"
    "testing"

    . "github.com/smartystreets/goconvey/convey"
    "github.com/vaksi/user_management/internal/healths"
    userSvcMocks "github.com/vaksi/user_management/internal/users/service/mocks"
    "github.com/vaksi/user_management/pkg/apiResponse"
)

func requestWithBody(method, path string, body interface{}) *goHttp.Request {
    b, err := json.Marshal(body)
    if err != nil {
        panic(err)
    }

    buf := bytes.NewBuffer(b)

    req := httptest.NewRequest(method, path, buf)
    req.Header.Add("Content-Type", "application/json")
    return req
}

func TestHandler_healthCheck(t *testing.T) {
    Convey("Health Check Status Endpoint Should Be Success", t, func() {
        // create mocks in services
        svcHealth := &healths.HealthService{}
        h := Handler{HealthService: svcHealth}
        route := h.MakeHandler()

        req := httptest.NewRequest("GET", PathPrefix+"/health-check", nil)
        w := httptest.NewRecorder()
        route.ServeHTTP(w, req)

        resp := apiResponse.APIOk
        resp.Message = h.HealthService.HealthCheckStatus()
        exp, _ := json.Marshal(resp)

        So(w.Code, ShouldEqual, goHttp.StatusOK)
        So(w.Body.String(), ShouldEqual, string(exp))
    })

    Convey("Health Check Status Endpoint Should Be NotFound", t, func() {
        // create mocks in services
        svcHealth := &healths.HealthService{}
        h := Handler{HealthService: svcHealth}
        route := h.MakeHandler()

        req := httptest.NewRequest("GET", PathPrefix+"/health-checks", nil)
        w := httptest.NewRecorder()
        route.ServeHTTP(w, req)

        So(w.Code, ShouldEqual, goHttp.StatusNotFound)
    })
}

func TestHandler_createUser(t *testing.T) {
    Convey("createUser endpoint should be success", t, func() {
        // param
        email := "test@email.com"
        fullName := "Test Name"
        gender := int8(1)

        // create mocks in user service
        svcUser := new(userSvcMocks.UserServiceFactory)
        svcUser.On("CreateUser", email, fullName, gender).Return(nil)
        h := Handler{UserService: svcUser}

        // Create Request
        route := h.MakeHandler()
        req := requestWithBody("POST",PathPrefix+"/users", map[string]interface{}{
            "email": email,
            "full_name": fullName,
            "gender": gender,
        })
        w := httptest.NewRecorder()
        route.ServeHTTP(w, req)

        // Response and expected
        resp := apiResponse.APICreated
        exp, _ := json.Marshal(resp)

        So(w.Body.String(), ShouldEqual, string(exp))
    })

    Convey("createUser endpoint should be invalid data error", t, func() {
        // param
        email := "test@email.com"
        fullName := "Test Name"
        gender := "data false"

        // create mocks in user service
        svcUser := new(userSvcMocks.UserServiceFactory)
        h := Handler{UserService: svcUser}

        // Create Request
        route := h.MakeHandler()
        req := requestWithBody("POST",PathPrefix+"/users", map[string]interface{}{
            "email": email,
            "full_name": fullName,
            "gender": gender,
        })
        w := httptest.NewRecorder()
        route.ServeHTTP(w, req)

        // Response and expected
        resp := apiResponse.APIErrorInvalidData
        exp, _ := json.Marshal(resp)

        So(w.Body.String(), ShouldEqual, string(exp))
    })

    Convey("createUser endpoint should be internal error", t, func() {
        // param
        email := "test@email.com"
        fullName := "Test Name"
        gender := int8(1)

        // create mocks in user service
        svcUser := new(userSvcMocks.UserServiceFactory)
        svcUser.On("CreateUser", email, fullName, gender).Return(fmt.Errorf("something error"))
        h := Handler{UserService: svcUser}

        // Create Request
        route := h.MakeHandler()
        req := requestWithBody("POST",PathPrefix+"/users", map[string]interface{}{
            "email": email,
            "full_name": fullName,
            "gender": gender,
        })
        w := httptest.NewRecorder()
        route.ServeHTTP(w, req)

        // Response and expected
        resp := apiResponse.APIErrorUnknown
        exp, _ := json.Marshal(resp)

        So(w.Body.String(), ShouldEqual, string(exp))
    })
}
