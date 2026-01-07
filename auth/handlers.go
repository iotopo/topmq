package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web"
	"github.com/iotopo/topmq/web/errcodes"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
	"time"
)

func authenticate(c *gin.Context) {
	user := GetLoginUser(c)
	data := gin.H{
		"id":          user.UID,
		"name":        user.Name,
		"username":    user.Username,
		"isSuperuser": user.IsSuperuser,
	}
	web.SuccessResponse(c, data)
}

func runWithFailedTimesCheck(c *gin.Context, captchaID, captchaValue string, fn func() (bool, error)) {
	key := LoginFailedTimes(c.ClientIP())
	times, err := cache.GetClient().Get(context.Background(), key).Int()

	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Warnf("failed getting login failed times: %v", err)
	}
	if times >= 3 {
		if captchaID == "" {
			web.HandleError(c, errcodes.BadCaptcha)
			return
		}
		if !verifyCaptcha(c.Request.Context(), captchaID, captchaValue) {
			web.HandleError(c, errcodes.BadCaptcha)
			return
		}
	}
	ok, err := fn()
	if err != nil {
		//if cache.IsAuthError(err) {
		//	web.HandleError(c, errcodes.RedisAuth)
		//	return
		//}
		web.HandleError(c, err)
		return
	}
	if ok {
		// 登录成功，删除失败次数
		_ = cache.GetClient().Del(context.Background(), key).Err()
	} else {
		// 登录失败，增加失败次数
		times++
		_ = cache.SetEx(key, times, time.Minute*5)
	}
}

func login(c *gin.Context) {
	var req struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		CaptchaID    string `json:"captcha_id"`
		CaptchaValue string `json:"captcha_value"`
	}
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if len(config.UsernameBlacklist) > 0 {
		if _, ok := config.UsernameBlacklist[req.Username]; ok {
			web.HandleError(c, errcodes.UserLocked)
			return
		}
	}
	req.Password = strings.TrimSpace(req.Password)
	if req.Password == "" {
		web.BadRequestResponse(c, "password is required")
		return
	}

	passwordBytes, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		web.HandleError(c, errcodes.BadUserOrPassword)
		return
	}
	req.Password = string(passwordBytes)

	runWithFailedTimesCheck(c, req.CaptchaID, req.CaptchaValue, func() (bool, error) {
		var user struct {
			ID           string
			Name         string
			PasswordHash []byte
			PasswordSalt []byte
			Locked       bool
			IsSuperuser  bool
			UpdatedAt    time.Time
		}

		selectFrom := clause.Table{Name: "users", Alias: "u"}
		conditions := []clause.Expression{
			clause.Eq{
				Column: clause.Column{Table: "u", Name: "username"},
				Value:  req.Username,
			},
		}

		columns := []clause.Column{
			{Table: "u", Name: "id"},
			{Table: "u", Name: "name"},
			{Table: "u", Name: "password_hash"},
			{Table: "u", Name: "password_salt"},
			{Table: "u", Name: "locked"},
			{Table: "u", Name: "is_superuser"},
			{Table: "u", Name: "updated_at"},
		}

		result := db.DB.Table("?", selectFrom).
			Clauses(clause.From{Tables: []clause.Table{selectFrom}}).
			Clauses(clause.Select{Columns: columns}).
			Clauses(clause.Where{Exprs: conditions}).
			Limit(1).
			Find(&user)

		if result.Error != nil {
			return false, result.Error
		}
		if result.RowsAffected == 0 {
			web.HandleError(c, errcodes.BadUserOrPassword)
			return false, nil
		}
		if user.Locked {
			return false, errcodes.UserLocked
		}

		if utils.CheckPasswordHash(req.Password, user.PasswordHash, user.PasswordSalt) {
			token, err := NewUserSession(context.Background(), UserSession{
				UserID:    user.ID,
				UpdatedAt: user.UpdatedAt,
				ClientIP:  c.ClientIP(),
			})
			if err != nil {
				return false, err
			}
			web.SuccessResponse(c, gin.H{
				"id":          user.ID,
				"name":        user.Name,
				"username":    req.Username,
				"isSuperuser": user.IsSuperuser,
				"token":       token,
			})
			// 向 context 里设置 LoginUser 原本是不需要的，但是后面的 OperationLog 会用到
			SetLoginUser(c, &LoginUser{
				UID:         user.ID,
				Name:        user.Name,
				Username:    req.Username,
				IsSuperuser: user.IsSuperuser,
			})
			return true, nil
		} else {
			web.HandleError(c, errcodes.BadUserOrPassword)
			return false, nil
		}
	})

}

func logout(c *gin.Context) {
	web.SuccessResponse(c, nil)

	if token := GetAuthToken(c); token != "" {
		DelToken(context.Background(), token)
	}
}

const captchaLength = 4
const captchaExpiration = time.Minute * 5

func LoginFailedTimes(ip string) string {
	return fmt.Sprintf(config.Conf.Redis.KeyPrefix+"login-failed-times:%s", ip)
}

func Captcha(captchaID string) string {
	return fmt.Sprintf(config.Conf.Redis.KeyPrefix+"captcha:%s", captchaID)
}

func isCaptchaShown(c *gin.Context) {
	type Response struct {
		Show bool `json:"show"`
	}

	loginFailedTimesKey := LoginFailedTimes(c.ClientIP())
	loginFailedTimes, err := cache.GetClient().Get(c.Request.Context(), loginFailedTimesKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Warnf("failed getting login failed times: %v", err)
	}
	if loginFailedTimes >= 3 {
		web.SuccessResponse(c, Response{Show: true})
	} else {
		web.SuccessResponse(c, Response{Show: false})
	}
}

func refreshCaptcha(c *gin.Context) {
	type Response struct {
		CaptchaID string `json:"captcha_id"`
	}
	captchaIDBytes := make([]byte, 10)
	_, _ = rand.Read(captchaIDBytes)
	captchaID := hex.EncodeToString(captchaIDBytes)
	captchaValue := captcha.RandomDigits(captchaLength)
	err := cache.SetEx(Captcha(captchaID), captchaValue, captchaExpiration)
	if err != nil {
		logrus.Errorf("failed setting captcha to redis: %v", err)
		return
	}
	web.SuccessResponse(c, Response{CaptchaID: captchaID})
}

func showCaptcha(c *gin.Context) {
	captchaFilename := c.Param("filename")
	if !strings.HasSuffix(captchaFilename, ".png") {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	captchaID := strings.TrimSuffix(captchaFilename, ".png")
	if captchaID == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	captchaValue, err := cache.GetClient().Get(c.Request.Context(), Captcha(captchaID)).Bytes()
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		if !errors.Is(err, redis.Nil) {
			logrus.Errorf("failed getting captcha from redis: %v", err)
		}
		return
	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Content-Type", "image/png")
	_, _ = captcha.NewImage(captchaID, captchaValue, captcha.StdWidth, captcha.StdHeight).WriteTo(c.Writer)
}

func verifyCaptcha(ctx context.Context, captchaID string, captchaString string) bool {
	if captchaString == "" {
		return false
	}
	digits := make([]byte, 0, len(captchaString))
	for _, d := range []byte(captchaString) {
		switch {
		case '0' <= d && d <= '9':
			digits = append(digits, d-'0')
		case d == ' ' || d == ',':
			// ignore
		default:
			return false
		}
	}
	if len(digits) == 0 {
		return false
	}
	captchaValue, err := cache.GetClient().Get(ctx, Captcha(captchaID)).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logrus.Errorf("failed getting captcha from redis: %v", err)
		}
		return false
	}
	// 无论验证成功还是失败，验证码在验证时都被删除
	cache.GetClient().Del(ctx, Captcha(captchaID))
	return bytes.Equal(digits, captchaValue)
}

func changePassword(c *gin.Context) {
	type RequestBody struct {
		OldPassword string `json:"oldPassword"`
		Password    string `json:"password"`
	}

	var req RequestBody
	if err := c.BindJSON(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}

	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.Password = strings.TrimSpace(req.Password)
	if req.OldPassword == "" {
		web.BadRequestResponse(c, "old password is required")
		return
	}
	if req.Password == "" {
		web.BadRequestResponse(c, "password is required")
		return
	}
	if len(req.Password) < config.Conf.MinPwdLen {
		web.BadRequestResponse(c, "password is too short")
		return
	}
	if !utils.IsValidPassword(req.Password) {
		web.BadRequestResponse(c, "password is invalid")
		return
	}

	token := GetLoginUser(c)

	user := User{ID: token.UID}
	result := db.DB.Select([]string{"password_hash", "password_salt"}).Limit(1).Find(&user)
	if result.Error != nil {
		web.HandleError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		web.HandleError(c, errcodes.NotFound)
		return
	}

	if utils.CheckPasswordHash(req.OldPassword, user.PasswordHash, user.PasswordSalt) {
		user.PasswordHash, user.PasswordSalt = utils.GeneratePasswordHash(req.Password)
		user.UpdatedAt = time.Now()
		result := db.DB.Select([]string{"password_hash", "password_salt", "updated_at"}).Save(&user)

		if result.Error != nil {
			web.HandleError(c, result.Error)
			return
		}
		if result.RowsAffected == 0 {
			web.HandleError(c, errcodes.NotFound)
			return
		}

		ClearUserSession(context.Background(), user.ID)

		web.SuccessResponse(c, nil)
	} else {
		// web.HandleError(c, errcodes.BadUserOrPassword)
		web.FailureResponse(c, "old_password_incorrect", "")
	}

}
