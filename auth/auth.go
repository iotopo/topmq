package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web/errcodes"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"net/http"
	"strings"
	"time"
)

var SessionExpiration = time.Hour

func GetSession(ctx context.Context, token string) (session UserSession, err error) {
	sessionKey := SessionKey(token)

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("get session panic: %v", e)
			cache.GetClient().Del(ctx, sessionKey)
		}
	}()

	var sessionBytes []byte
	sessionBytes, err = cache.GetClient().Get(ctx, sessionKey).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			err = fmt.Errorf("get session failed: %w", err)
		}
		return
	}

	if err = msgpack.Unmarshal(sessionBytes, &session); err != nil {
		err = fmt.Errorf("failed unmarshal session: %w", err)
		return
	}
	return
}

func KeepAlive(token string) {
	sessionKey := SessionKey(token)
	if err := cache.GetConn().Expire(context.Background(), sessionKey, SessionExpiration).Err(); err != nil {
		logrus.WithError(err).Error("extends session token expire error")
	}
}

func KeepSession(token string) {
	err := cache.GetConn().Expire(context.Background(), SessionKey(token), SessionExpiration).Err()
	if err != nil {
		logrus.WithError(err).Error("extends session token expire error")
	}
}
func NewAppSession(ctx context.Context, session UserSession) (string, error) {
	token := utils.UUIDWithoutDash()
	err := saveSession(ctx, session, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewUserSession(ctx context.Context, session UserSession) (string, error) {
	token := utils.UUIDWithoutDash()
	err := saveSession(ctx, session, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func saveSession(ctx context.Context, session UserSession, token string) error {
	sessionBytes, err := msgpack.Marshal(session)
	if err != nil {
		return err
	}
	sessionKey := SessionKey(token)
	return cache.GetClient().Set(ctx, sessionKey, sessionBytes, SessionExpiration).Err()
}

func DelToken(ctx context.Context, token string) {
	if err := cache.GetClient().Del(ctx, SessionKey(token)).Err(); err != nil {
		logrus.WithError(err).Errorf("failed deleting session")
	}
}

// ClearUserSession 清除用户会话信息
// 锁定用户、删除用户、重置用户密码后调用
// 可替代原有的 UpdatedAt 失效方式
func ClearUserSession(ctx context.Context, userID string) {
	if userID == "" {
		return
	}
	defer DelUserAuthCache(ctx, userID)

	client := cache.GetClient()
	userSessionKey := UserSessionKey(userID)
	expiredTokenMap, err := client.HGetAll(ctx, userSessionKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.WithError(err).Error("failed getting existed sessions")
		return
	}
	if len(expiredTokenMap) > 0 {
		expiredTokens := make([]string, len(expiredTokenMap))
		expiredSessionKeys := make([]string, len(expiredTokenMap))
		for endport, token := range expiredTokenMap {
			logrus.Debugf("kickout session %q of user %q (endport: %s)", token, userID, endport)
			expiredTokens = append(expiredTokens, token)
			expiredSessionKeys = append(expiredSessionKeys, SessionKey(token))
		}
		err := client.Del(ctx, expiredSessionKeys...).Err()
		if err != nil {
			logrus.WithError(err).Error("failed deleting existed sessions")
		}
	}
}

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
		user, err = ValidateUserToken(ctx, csrfToken)
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
				var session UserSession
				session, err = GetSession(ctx, token)
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

					var userInfo *UserForAuthenticate
					userInfo, err = GetForAuthenticate(ctx, session.UserID)
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
							DelToken(ctx, token)
							err = errcodes.SessionExpired
							c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
							return
						}
					}

					// 保活
					KeepAlive(token)
					user = &LoginUser{
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
				KeepSession(token)
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

func init() {
	sessionTimeout := config.Conf.Web.SessionTimeout
	if sessionTimeout > 0 {
		SessionExpiration = time.Duration(sessionTimeout) * time.Minute
	}
}
