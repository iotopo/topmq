package errcodes

import (
	"fmt"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	StatusCode int
	ErrorCode  string
	ErrorMsg   string
	ErrorArgs  any
}

func (e ErrorResponse) Is(another error) bool {
	if e2, ok := another.(ErrorResponse); ok {
		if e.StatusCode != e2.StatusCode {
			return false
		}
		if e.ErrorCode != e2.ErrorCode {
			return false
		}
		return true
		//if e.ErrorMsg != e2.ErrorMsg {
		//	return false
		//}
		//return cmp.Equal(e.ErrorArgs, e2.ErrorArgs)
	}
	return false
}

func (e ErrorResponse) Error() string {
	if e.ErrorMsg != "" {
		if e.ErrorArgs == nil {
			return e.ErrorMsg
		}
		return fmt.Sprintf("%s(%v)", e.ErrorMsg, e.ErrorArgs)
	}
	return e.ErrorCode
}

//goland:noinspection GoUnusedExportedFunction
func NewHttpError(statusCode int, errorMessage string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		ErrorCode:  strconv.Itoa(statusCode),
		ErrorMsg:   errorMessage,
		//ErrorArgs:  args,
	}
}

func NewFailure(errorCode, errorMessage string) ErrorResponse {
	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  errorCode,
		ErrorMsg:   errorMessage,
		//ErrorArgs:  args,
	}
}
func NewFailureWithArgs(errorCode, errorMessage string, args any) ErrorResponse {
	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  errorCode,
		ErrorMsg:   errorMessage,
		ErrorArgs:  args,
	}
}

// 简化形式：errorMessage 为空
//
//goland:noinspection GoUnusedExportedFunction
func NewSimpleFailure(errorCode string) ErrorResponse {
	return NewFailure(errorCode, "")
}

type FailureTemplate struct {
	ErrorCode string
	ErrorMsg  string
}

func (t FailureTemplate) Build(args any) ErrorResponse {
	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  t.ErrorCode,
		ErrorMsg:   t.ErrorMsg,
		ErrorArgs:  args,
	}
}

//goland:noinspection GoUnusedExportedFunction
func NewFailureTemplate(errorCode, format string) FailureTemplate {
	return FailureTemplate{
		ErrorCode: errorCode,
		ErrorMsg:  format,
	}
}
