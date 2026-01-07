package auth

import "fmt"

func CsrfToken(token string) string {
	return fmt.Sprintf("topmq:csrf-token:%s", token)
}
func TempToken(token string) string {
	return fmt.Sprintf("topmq:temp-token:%s", token)
}

// UserForAuthenticateV2 新版用户认证缓存
func UserForAuthenticateV2(userID string) string {
	return fmt.Sprintf("topmq:authenticate_v2:%s", userID)
}

func SessionKey(sessionKey string) string {
	return fmt.Sprintf("topmq:session:%s", sessionKey)
}

func UserSessionKey(userID string) string {
	return fmt.Sprintf("topmq:user-session:%s", userID)
}
