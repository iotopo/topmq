package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/web/errcodes"
	"net/http"
)

var (
	ErrInvalidToken = errcodes.NewHttpError(http.StatusUnauthorized, "invalid auth token")
)

type LoginUser = auth.LoginUser

type contextKeyLoginUser struct{}

func GetCtxLoginUser(ctx context.Context) *LoginUser {
	if user, err := GetContextValue[*LoginUser](ctx, contextKeyLoginUser{}); err == nil {
		return user
	}
	return nil
}

func SetLoginUser(c *gin.Context, user *LoginUser) {
	SetRequestContextValue(c, contextKeyLoginUser{}, user)
	//c.Set(userKey, user)
}

func GetLoginUser(c *gin.Context) *LoginUser {
	if user, err := GetRequestContextValue[*LoginUser](c, contextKeyLoginUser{}); err == nil {
		return user
	}
	//if v, ok := c.Get(userKey); ok {
	//	return v.(*LoginUser)
	//}
	return nil
}
