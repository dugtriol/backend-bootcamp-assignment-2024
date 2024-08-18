package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message   string `json:"message,omitempty"`
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
}

//type Response struct {
//	Status string `json:"status"`
//	Error  string `json:"error"`
//}
//
//const (
//	StatusOK    = "OK"
//	StatusError = "Message"
//)
//
//func OK() Response {
//	return Response{
//		Status: StatusOK,
//
//	}
//}
//
//func Error(msg string) Response {
//	return Response{
//		Status: StatusError,
//		Error:  msg,
//	}
//}

func MakeResponse(message, requestId string, code int) Response {
	return Response{
		Message:   message,
		RequestId: requestId,
		Code:      code,
	}

}

func ValidationError(errs validator.ValidationErrors, requestId string) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Code:      http.StatusBadRequest,
		Message:   strings.Join(errMsgs, ", "),
		RequestId: requestId,
	}

}
