package errcode

import (
	"fmt"
	"net/http"
)

type ErrorInfo struct {
	HttpCode int    `json:"httpCode"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

func New(errorInfo *ErrorInfo, params ...interface{}) *ErrorInfo {
	return &ErrorInfo{
		Code:    errorInfo.Code,
		Message: fmt.Sprintf(errorInfo.Message, params...),
	}
}

var (
	InvalidParamValue     = &ErrorInfo{http.StatusBadRequest, 900, "Param value invalid."}
	InternalServerError   = &ErrorInfo{http.StatusInternalServerError, 901, "Server error"}
	ResourceNotFound      = &ErrorInfo{http.StatusBadRequest, 404, "resource not found"}
	ResourceAlreadyExists = &ErrorInfo{http.StatusBadRequest, 408, "resource already exists"}
)
