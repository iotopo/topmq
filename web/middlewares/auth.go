package middlewares

import (
	"context"
	"fmt"
	auth2 "github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/web/errcodes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Authorize 统一身份认证拦截器，身份认证通过后，自动检查是否指定了项目
func Authorize(c *gin.Context) {
	if u := GetLoginUser(c); u != nil { // 避免重复验证
		return
	}
	//// 限速检测
	//if IsRateLimited() {
	//	c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "rate limit exceeded"})
	//	return
	//}

	var user *LoginUser
	var err error

	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" {
		csrfToken = c.Query("_csrf_token_")
		if csrfToken == "" {
			csrfToken = c.Query("_temp_token_")
		}
	}
	ctx := context.Background()
	if csrfToken != "" { // 临时 token
		user, err = auth2.ValidateUserToken(ctx, csrfToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "bad csrf token", "err": err.Error()})
			return
		}
		SetLoginUser(c, user)
	} else { // 授权 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authHeader = c.Query("_auth_token_")
		}
		if authHeader != "" { // 常规登录身份认证
			// SetAuthToken(c, token)
			//c.Set("auth_token", token)
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				var session auth2.UserSession
				session, err = auth2.GetSession(ctx, token)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "validate bearer token failed", "error": err.Error()})
					return
				}

				if session.UserID != "" { // 用户
					clientIP := c.ClientIP()
					if config.Conf.Web.CheckIPChange {
						if session.ClientIP == "" || session.ClientIP != clientIP {
							logrus.Debugf("session ip changed: %s -> %s", session.ClientIP, clientIP)
							err = errcodes.SessionIPChanged
							c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session IP changed"})
							return
						}
					}

					var userInfo *auth2.UserForAuthenticate
					userInfo, err = auth2.GetForAuthenticate(ctx, session.UserID)
					if err != nil {
						err = fmt.Errorf("failed verifying user from token(%w)", err)
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
						return
					}
					if userInfo == nil {
						err = fmt.Errorf("user %s not found", session.UserID)
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
						return
					}
					if userInfo.Locked {
						err = fmt.Errorf("user %s is locked", session.UserID)
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
						return
					}
					if !userInfo.UpdatedAt.IsZero() {
						// 如果用户的 UpdatedAt 比会话中的新，强制用户下线
						// 用户信息已更新，强制让会话失效
						if session.UpdatedAt.IsZero() || session.UpdatedAt.Before(userInfo.UpdatedAt) {
							auth2.DelToken(ctx, token)
							err = errcodes.SessionExpired
							c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
							return
						}
					}

					// 保活
					auth2.KeepAlive(token)
					user = &auth2.LoginUser{
						UID:      session.UserID,
						Name:     userInfo.Name,
						Username: userInfo.Username,
					}
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "bad auth authHeader"})
					return
				}
				SetAuthToken(c, token)
				SetLoginUser(c, user)
				auth2.KeepSession(token)
			} else {
				//web.HandleError(c, ErrBadToken) // 401
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "bad auth authHeader"})
				return
			}
		} else { // 免登录访问，身份认证
			//apiKey := c.GetHeader(ApiKeyHeaderKey)
			//if apiKey != "" {
			//	conf, err := service.OpenApi.GetConfig()
			//	if err != nil {
			//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "failed to load open api settings", "err": err.Error()})
			//		return
			//	}
			//	if conf == nil {
			//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "open api not enabled"})
			//		return
			//	}
			//	if !isOpenApiAllow(conf) {
			//		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"msg": "request rate limited"})
			//		return
			//	}
			//	user = service.NewAnonymousAdminUser()
			//	SetLoginUser(c, user)
			//}
		}
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "auth token expected"})
		return
	}
}
