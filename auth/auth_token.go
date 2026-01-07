package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/web/errcodes"
	"net/http"
)

var (
	ErrAuthTokenNotSpecified = errcodes.NewHttpError(http.StatusUnauthorized, "auth token is not specified")
	ErrBadToken              = errcodes.NewHttpError(http.StatusUnauthorized, "bad auth token")
)

type contextKeyAuthToken struct{}

//goland:noinspection GoUnusedExportedFunction
func GetAuthToken(c *gin.Context) string {
	if token, err := GetRequestContextValue[string](c, contextKeyAuthToken{}); err == nil {
		return token
	}
	return ""
}

func SetAuthToken(c *gin.Context, token string) {
	SetRequestContextValue[string](c, contextKeyAuthToken{}, token)
}
