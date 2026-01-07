package accounts

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web"
	"strings"
)

func createAccount(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		ClientID    string `json:"clientID"`
		Remote      string `json:"remote"`
		Description string `json:"description"`
	}
	if err := c.ShouldBind(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	user := Account{
		ID:          utils.UUIDWithoutDash(),
		Username:    req.Username,
		Password:    req.Password,
		ClientID:    req.ClientID,
		Remote:      req.Remote,
		Description: req.Description,
	}
	err := db.DB.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: accounts.username") {
			web.FailureResponseWithArgs(c, "err_username_duplicated", err.Error(), nil)
			return
		}
		web.InternalErrorResponse(c, err)
		return
	}
	_ = accountCache.Delete(context.Background(), getUserCacheKey(user.Username))
	web.SuccessResponse(c, user.ID)
}

func updateAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		ClientID    string `json:"clientID"`
		Remote      string `json:"remote"`
		Description string `json:"description"`
	}
	if err := c.ShouldBind(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	user := Account{
		ID:          id,
		Username:    req.Username,
		Password:    req.Password,
		ClientID:    req.ClientID,
		Remote:      req.Remote,
		Description: req.Description,
	}
	err := db.DB.Select("username", "password", "client_id", "remote", "description").Save(&user).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	web.SuccessResponse(c, nil)
}

func deleteAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	err := runAndClearCache(id, func() error {
		return db.DB.Delete(&Account{}, "id=?", id).Error
	})
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	web.SuccessResponse(c, nil)
}

func enabledAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	err := runAndClearCache(id, func() error {
		return db.DB.Model(&Accounts).Where("id=?", id).Update("disabled", false).Error
	})
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	web.SuccessResponse(c, nil)
}

func disabledAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	err := runAndClearCache(id, func() error {
		return db.DB.Model(&Accounts).Where("id=?", id).Update("disabled", true).Error
	})
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	web.SuccessResponse(c, nil)
}

func getAccounts(c *gin.Context) {
	var req struct {
		PageNum  int    `form:"pageNum"`
		PageSize int    `form:"pageSize"`
		Search   string `form:"search"`
	}
	if err := c.BindQuery(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	q := db.DB.Model(&Accounts)
	if req.Search != "" {
		q = q.Where(db.ClauseExprForContain(q, "username", req.Search))
	}
	var result struct {
		Total int64 `json:"total"`
		Items []struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			Password    string `json:"password"`
			ClientID    string `json:"clientID"`
			Remote      string `json:"remote"`
			Disabled    bool   `json:"disabled"`
			Description string `json:"description"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"items"`
	}
	err := q.Count(&result.Total).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	if result.Total > 0 {
		err = q.Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Order("created_at desc").Find(&result.Items).Error
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
	}
	web.SuccessResponse(c, result)
}

//func getFilters(c *gin.Context) {
//	id := c.Param("id")
//	if id == "" {
//		web.BadRequestResponse(c, "id is empty")
//		return
//	}
//	var result struct {
//		Filters map[string]string `json:"filters" gorm:"serializer:json"`
//	}
//	err := db.DB.Model(&Accounts).Where("id = ?", id).First(&result).Error
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			web.HandleError(c, errcodes.NotFound)
//		} else {
//			web.InternalErrorResponse(c, err)
//		}
//		return
//	}
//	web.SuccessResponse(c, result.Filters)
//}

func Init() {
	router := web.Router
	api := router.Group("api/v1/account", auth.Authorize)
	api.GET("", getAccounts)
	api.POST("", createAccount)
	api.PUT(":id", updateAccount)
	api.DELETE(":id", deleteAccount)
	api.POST(":id/disable", disabledAccount)
	api.POST(":id/enable", enabledAccount)
	//api.GET(":id/filters", getFilters)
	//api.PUT(":id/filters", updateFilters)
}
