/*  apiResponses.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/10/18 13:56 
 */

package apiResponse

import (
"encoding/json"
"net/http"

"bitbucket.org/kudoindonesia/microservice_user_management/constants"
)

// APIResponse defines attributes for api Response
type APIResponse struct {
    HTTPCode int         `json:"-"`
    Code     int         `json:"code"`
    Message  interface{} `json:"message"`
    Data     interface{} `json:"data,omitempty"`
}

// Write writes the data to http response writer
func Write(w http.ResponseWriter, response APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(response.HTTPCode)
    js, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
   if _, err := w.Write(js); err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   }
}

// WithMessage modifies api response's message
func (a *APIResponse) WithMessage(message interface{}) APIResponse {
    a.Message = message
    return *a
}

// WithHTTPCode modifies api response's http code
func (a *APIResponse) WithHTTPCode(httpCode int) APIResponse {
    a.HTTPCode = httpCode
    return *a
}

// Defaults API Response with standard HTTP Status Code. The default value can
// be changed either using `ModifyMessage` or `ModifyHTTPCode`. You can call it
// directly by :
//
// response.Write(res, response.APIErrorUnknown)
// return
var (
    // A generic error message, given when an unexpected condition was encountered.
    APIErrorUnknown = APIResponse{
        HTTPCode: http.StatusInternalServerError,
        Code:     constants.ErrorUnknownCode,
        Message:  "Internal Server Error",
    }

    // ApiOk Standard response for successful HTTP requests.
    APIOk = APIResponse{
        HTTPCode: http.StatusOK,
        Code:     constants.GeneralSuccessCode,
        Message:  "Success",
    }

    // The request has been fulfilled, resulting in the creation of a new resource
    APICreated = APIResponse{
        HTTPCode: http.StatusCreated,
        Code:     constants.GeneralSuccessCode,
        Message:  "Success",
    }

    // The request has been accepted for processing, but the processing has not been completed.
    APIAccepted = APIResponse{
        HTTPCode: http.StatusAccepted,
        Code:     constants.GeneralSuccessCode,
        Message:  "Success",
    }

    // The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
    // , size too large, invalid request message framing, or deceptive request routing).
    APIErrorValidation = APIResponse{
        HTTPCode: http.StatusBadRequest,
        Code:     constants.InvalidKeyFieldValue,
        Message:  "Data tidak sesuai",
    }

    // The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
    // , size too large, invalid request message framing, or deceptive request routing).
    APIErrorInvalidPassword = APIResponse{
        HTTPCode: http.StatusBadRequest,
        Code:     constants.ErrorInvalidPassword,
        Message:  "Invalid Password or Key",
    }

    APIErrorInvalidData = APIResponse{
        HTTPCode: http.StatusBadRequest,
        Code:     constants.ErrorInvalidData,
        Message:  "Invalid Input Data",
    }

    // The request was valid, but the server is refusing Action. The user might not have the necessary permissions for
    // a resource, or may need an account of some sort.
    APIErrorForbidden = APIResponse{
        HTTPCode: http.StatusForbidden,
        Code:     constants.ErrorForbiddenCode,
        Message:  "Action forbidden",
    }

    // The request was valid, but the server is refusing Action. The user might not have the necessary permissions for
    // a resource, or may need an account of some sort.
    APIErrorUnauthorized = APIResponse{
        HTTPCode: http.StatusUnauthorized,
        Code:     constants.ErrorUnauthorizedCode,
        Message:  "Unauthorized",
    }
)

// DataErrors This function to return standard data errors
func DataErrors(data interface{}) map[string]interface{} {
    return map[string]interface{}{
        "errors": data,
    }
}
