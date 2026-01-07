package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web/errcodes"
	"time"

	"github.com/redis/go-redis/v9"
)

type CsrfUser struct {
	UserID string `json:"user_id"`
}

// 生成一个短期有效的 token，并将其与 value 关联
func GenerateToken(ctx context.Context, value string) (token string, err error) {
	token = utils.UUIDWithoutDash()
	err = cache.GetClient().Set(ctx, CsrfToken(token), value, time.Minute).Err()
	return
}

// 根据 token 查询关联的 value，查询成功时，将 token 作废
func ValidateToken(ctx context.Context, token string) (value string, err error) {
	key := CsrfToken(token)
	client := cache.GetClient()
	value, err = client.Get(ctx, key).Result()
	if err == nil {
		_ = client.Del(ctx, key).Err()
	}
	if errors.Is(err, redis.Nil) {
		err = errcodes.InvalidCsrfToken
	}
	return
}

// 和 GenerateToken 类似，但生成的 token 在指定时间内始终有效，验证成功后不会作废
func GenerateTempToken(ctx context.Context, value string, ttl time.Duration) (token string, err error) {
	token = utils.UUIDWithoutDash()
	err = cache.GetClient().Set(ctx, TempToken(token), value, ttl).Err()
	return
}

// 验证 GenerateTempToken 生成的 token
func ValidateTempToken(ctx context.Context, token string) (value string, err error) {
	key := TempToken(token)
	value, err = cache.GetClient().Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = errcodes.InvalidCsrfToken
	}
	return
}

// GenerateUserToken 生成一个 token，可以用于一次性确定用户身份及所操作的项目
// 注意：token 必须在生成时指定项目，而不能在使用 token 时由前端指定项目
func GenerateUserToken(ctx context.Context, userID string) (token string, err error) {
	var value []byte
	value, err = json.Marshal(CsrfUser{UserID: userID})
	if err != nil {
		return
	}
	token, err = GenerateToken(ctx, string(value))
	return
}

func ValidateUserToken(ctx context.Context, token string) (loginUser *LoginUser, err error) {
	var value string
	value, err = ValidateToken(ctx, token)
	if err != nil {
		return
	}

	var csrfUser CsrfUser
	if err = json.Unmarshal([]byte(value), &csrfUser); err != nil {
		err = errcodes.InvalidCsrfToken
		return
	}

	var user *UserForAuthenticate
	user, err = GetForAuthenticate(context.Background(), csrfUser.UserID)
	if err != nil {
		return
	}
	if user == nil {
		err = errcodes.InvalidCsrfToken
		return
	}
	if user.Locked {
		err = errcodes.InvalidCsrfToken
		return
	}

	loginUser = &LoginUser{
		UID:      csrfUser.UserID,
		Name:     user.Name,
		Username: user.Username,
	}
	return
}
