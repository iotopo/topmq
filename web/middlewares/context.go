package middlewares

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrContextNotSet        = errors.New("context not set")
	ErrContextNotRecognized = errors.New("context not recognized")
)

type ContextValue[T any] struct {
	Value T
}

type ContextError struct {
	Error error
}

func SetContextValue[T any](ctx context.Context, key any, value T) context.Context {
	return context.WithValue(ctx, key, ContextValue[T]{Value: value})
}

func SetContextError(ctx context.Context, key any, err error) context.Context {
	return context.WithValue(ctx, key, ContextError{Error: err})
}

func GetContextValue[T any](ctx context.Context, key any) (value T, err error) {
	item := ctx.Value(key)
	if item == nil {
		err = ErrContextNotSet
		return
	}
	if ctxValue, ok := item.(ContextValue[T]); ok {
		value = ctxValue.Value
		return
	}
	if ctxError, ok := item.(ContextError); ok {
		err = ctxError.Error
		return
	}
	err = ErrContextNotRecognized
	return
}

func SetRequestContextValue[T any](c *gin.Context, key any, value T) {
	c.Request = c.Request.WithContext(SetContextValue(c.Request.Context(), key, value))
}

func SetRequestContextError(c *gin.Context, key any, err error) {
	c.Request = c.Request.WithContext(SetContextError(c.Request.Context(), key, err))
}

func GetRequestContextValue[T any](c *gin.Context, key any) (value T, err error) {
	return GetContextValue[T](c.Request.Context(), key)
}
