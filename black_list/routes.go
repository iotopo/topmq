package black_list

import (
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web"
	"gorm.io/gorm/clause"
	"time"
)

func getBlackList(c *gin.Context) {
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
	q := db.DB.Model(&BlackLists)
	if req.Search != "" {
		q = q.Where(clause.Or(
			db.ClauseExprForContain(q, "username", req.Search),
			db.ClauseExprForContain(q, "client_id", req.Search),
			db.ClauseExprForContain(q, "remote", req.Search),
		))
	}
	var result struct {
		Total int64 `json:"total"`
		Items []struct {
			ID          string     `json:"id"`
			Username    string     `json:"username"`
			ClientID    string     `json:"clientID"`
			Remote      string     `json:"remote"`
			Allow       bool       `json:"allow"`
			Description string     `json:"description"`
			ExpiredAt   *time.Time `json:"expiredAt"`
			CreatedAt   string     `json:"createdAt"`
			UpdatedAt   string     `json:"updatedAt"`
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

func createBlackList(c *gin.Context) {
	var req struct {
		Username    string    `json:"username"`
		ClientID    string    `json:"clientID"`
		Remote      string    `json:"remote"`
		Description string    `json:"description"`
		ExpiredAt   time.Time `json:"expiredAt"`
	}
	if err := c.ShouldBind(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	auth := BlackList{
		ID:          utils.UUIDWithoutDash(),
		Username:    req.Username,
		ClientID:    req.ClientID,
		Remote:      req.Remote,
		Description: req.Description,
		ExpiredAt:   &req.ExpiredAt,
	}
	err := db.DB.Create(&auth).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	clearCacheLater()
	web.SuccessResponse(c, auth.ID)
}

func updateBlackList(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	var req struct {
		Username    string    `json:"username"`
		ClientID    string    `json:"clientID"`
		Remote      string    `json:"remote"`
		Description string    `json:"description"`
		ExpiredAt   time.Time `json:"expiredAt"`
	}
	if err := c.ShouldBind(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	auth := BlackList{
		ID:          id,
		Username:    req.Username,
		ClientID:    req.ClientID,
		Remote:      req.Remote,
		Description: req.Description,
		ExpiredAt:   &req.ExpiredAt,
	}
	err := db.DB.Select("client_id", "username", "remote", "allow").Save(&auth).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	clearCacheLater()
	web.SuccessResponse(c, nil)
}

func deleteBlackList(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	err := db.DB.Delete(&BlackList{}, "id=?", id).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	clearCacheLater()
	web.SuccessResponse(c, nil)
}

func Init() {
	router := web.Router
	api := router.Group("api/v1/black_list", auth.Authorize)
	api.GET("", getBlackList)
	api.POST("", createBlackList)
	api.PUT(":id", updateBlackList)
	api.DELETE(":id", deleteBlackList)
}
