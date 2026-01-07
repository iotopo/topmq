package auth

import (
	"github.com/iotopo/topmq/web"
)

func Init() {
	//router gin.IRouter
	router := web.Router

	router.POST("api/v1/login", login)
	router.POST("api/v1/logout", logout)
	router.GET("api/v1/logout", logout)

	captchaRouter := router.Group("api/v1/captcha")
	captchaRouter.GET("show", isCaptchaShown)
	captchaRouter.GET("refresh", refreshCaptcha)
	captchaRouter.GET("server/:filename", showCaptcha)

	loginGroup := router.Group("api/v1", Authorize)
	loginGroup.GET("authenticate", authenticate)
	loginGroup.PUT("reset_password", changePassword)

	adminApi := router.Group("api/v1/user", Authorize)
	adminApi.GET("paged", getPaged)
	adminApi.GET("username_test", usernameTest)
	adminApi.POST("", create)
	adminApi.PUT(":id", modify)
	adminApi.PUT(":id/lock", lock)
	adminApi.PUT(":id/unlock", unlock)
	adminApi.POST(":id/reset_password", resetPassword)
	adminApi.DELETE(":id", delete_)

	selfApi := router.Group("api/v1/self", Authorize)
	selfApi.GET("", getSelf)
	selfApi.PUT("", modifySelf)
	selfApi.POST("temp_token", genCsrf)
}
