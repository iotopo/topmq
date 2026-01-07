package web

import (
	"errors"
	"fmt"
	"github.com/iotopo/topmq/web/errcodes"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type APIResponse struct {
	ErrorCode    string      `json:"code,omitempty"`
	ErrorMessage string      `json:"msg,omitempty"`
	ErrArgs      any         `json:"errArgs,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Success      bool        `json:"success,omitempty"` // 兼容旧版本
}

func RequiredResponse(c *gin.Context, name string) {
	BadRequestResponse(c, fmt.Sprintf("%s is required", name))
}

func BadRequestResponse(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
		ErrorCode:    "400",
		ErrorMessage: msg,
	})
}

func InternalErrorResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
		ErrorCode:    "500",
		ErrorMessage: err.Error(),
	})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Data: data, Success: true})
}

func FailureResponse(c *gin.Context, code, msg string) {
	c.JSON(http.StatusOK, APIResponse{
		ErrorCode:    code,
		ErrorMessage: msg,
		//ErrArgs:      args,
	})
}

func FailureResponseWithArgs(c *gin.Context, code, msg string, args any) {
	c.JSON(http.StatusOK, APIResponse{
		ErrorCode:    code,
		ErrorMessage: msg,
		ErrArgs:      args,
	})
}

// HandleError 统一使用 errcodes 包下面的错误码实现
func HandleError(c *gin.Context, e error) {
	var errorResponse errcodes.ErrorResponse
	if !errors.As(e, &errorResponse) {
		InternalErrorResponse(c, e)
	} else {
		c.AbortWithStatusJSON(errorResponse.StatusCode, APIResponse{
			ErrorCode:    errorResponse.ErrorCode,
			ErrorMessage: errorResponse.ErrorMsg,
			ErrArgs:      errorResponse.ErrorArgs,
		})
	}
}

func ExcelFile(c *gin.Context, workbook *excelize.File, filename string) (err error) {
	c.Status(http.StatusOK)
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition#syntax
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s.xlsx", url.QueryEscape(filename)))
	_, err = workbook.WriteTo(c.Writer)
	return
}
