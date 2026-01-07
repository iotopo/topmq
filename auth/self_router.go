package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/web"
	"github.com/iotopo/topmq/web/errcodes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

func getSelf(c *gin.Context) {
	type Response struct {
		Username string `json:"username" gorm:"-"`
		Name     string `json:"name" gorm:"-"`
		Email    string `json:"email"`
		Tel      string `json:"tel"`
		Mobile   string `json:"mobile"`

		// WechatWebBound   bool    `json:"wechat_web_bound" gorm:"-"`
		WechatMiniBound bool `json:"wechat_mini_bound" gorm:"-"`
		// WechatUnionID    *string `json:"-"`
		// WechatWebOpenID  *string `json:"-"`
		// WechatMiniOpenID *string `json:"-"`
	}
	token := GetLoginUser(c)
	var res Response
	user := User{ID: token.UID}
	columns := []string{"email", "tel", "mobile"}
	// if config.Config.Wechat.OpenPlatformEnabled {
	// 	columns = append(columns, "wechat_union_id")
	// } else {
	// 	columns = append(columns, "wechat_web_open_id", "wechat_mini_open_id")
	// }
	result := db.DB.Model(&user).
		Where(&user, "id").
		Select(columns).
		Limit(1).
		Find(&res)
	if result.Error != nil {
		web.HandleError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	res.Name = token.Name
	res.Username = token.Username
	// if config.Config.Wechat.OpenPlatformEnabled {
	// 	res.WechatWebBound = res.WechatUnionID != nil
	// 	res.WechatMiniBound = res.WechatUnionID != nil
	// } else {
	// 	res.WechatWebBound = res.WechatWebOpenID != nil
	// 	res.WechatMiniBound = res.WechatMiniOpenID != nil
	// }
	//wechatAccount, err := service.WechatMini.GetOpenID(user.ID)
	//if err != nil {
	//	web.HandleError(c, result.Error)
	//	return
	//}
	//res.WechatMiniBound = wechatAccount != ""

	web.SuccessResponse(c, res)
}

func modifySelf(c *gin.Context) {
	type RequestBody struct {
		Name string `json:"name"`
		Tel  string `json:"tel"`
		// Email  string `json:"email"`
		// Email string `json:"mobile"`
	}

	var req RequestBody
	if err := c.BindJSON(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}

	token := GetLoginUser(c)

	user := User{
		ID:   token.UID,
		Name: req.Name,
	}
	if req.Tel != "" {
		user.Tel = &req.Tel
	}
	// if req.Email != "" {
	//	user.Email = &req.Email
	// }
	// if req.Email != "" {
	//	user.Email = &req.Email
	// }
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.
			Model(&user).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select([]string{"username"}).
			Limit(1).
			Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errcodes.NotFound
		}
		return tx.Select([]string{"name", "tel"}).Save(&user).Error
	})
	if err != nil {
		web.HandleError(c, err)
		return
	}
	web.SuccessResponse(c, nil)
	DelUserAuthCache(context.Background(), user.ID)
}

func genCsrf(c *gin.Context) {
	user := GetLoginUser(c)
	token, err := GenerateUserToken(context.Background(), user.UID)
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	web.SuccessResponse(c, gin.H{
		"token": token,
	})
}
