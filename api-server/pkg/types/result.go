package types

import (
	"context"
	"errors"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"kubeall.io/api-server/pkg/infra/constants"
	"net/http"
)

// Result represents the result of a request
type Result struct {
	ErrorCode   constants.ErrorCode `json:"code"`
	Payload     any                 `json:"payload,omitempty"`
	Message     string              `json:"message,omitempty"`
	FieldErrors map[string]string   `json:"fieldErrors,omitempty"`
	StatusCode  int                 `json:"statusCode"`
}

func (r Result) Error() string {
	return r.Message
}

func Fail(err error) *Result {
	var r *Result
	if errors.As(err, &r) {
		return r
	}
	return &Result{
		ErrorCode:  "UNEXPECTED_ERROR",
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}

func FailWithPayLoad(payload any, err error) *Result {
	result := Fail(err)
	result.Payload = payload
	return result
}

// FailWithErrorCode returns the error code associated with the given error code
func FailWithErrorCode(ctx context.Context, errorCode constants.ErrorCode, params map[string]string) *Result {
	var content = GetI18nMessage(ctx, errorCode, params)
	return &Result{
		ErrorCode:  errorCode,
		Message:    content,
		StatusCode: http.StatusBadRequest,
	}
}

func FailWithStatusCode(statusCode int) *Result {
	return &Result{StatusCode: statusCode}
}

func IsResult(err error) (bool, *Result) {
	var r *Result

	//为什么 target 是 nil 也能取地址？：
	//
	//在 Go 中，任何变量（包括 nil 指针）都有一个内存地址。target 是一个变量，存储在栈上，即使它的值是 nil，&target 仍然会返回 target 变量本身的地址。
	//errors.As 并不关心 target 的值是否为 nil，它只关心 target 的类型是否匹配（*Result）以及 &target 是否是一个有效的指针地址（**Result）。
	//当 errors.As 找到匹配的错误类型时，它会通过 &target 修改 target 的值，将其从 nil 更新为具体的 *Result 值。
	return errors.As(err, &r), r
}

func GetI18nMessage(ctx context.Context, errorCode constants.ErrorCode, params map[string]string) string {
	ginCtx := ctx.(*gin.Context)
	var content string
	if params != nil {
		content = ginI18n.MustGetMessage(ginCtx, &i18n.LocalizeConfig{
			MessageID:    string(errorCode),
			TemplateData: params,
		})
	} else {
		content = ginI18n.MustGetMessage(ginCtx, errorCode)
	}
	return content
}
