package acls

import (
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web"
)

func getAclRules(c *gin.Context) {
	var req struct {
		PageNum  int    `form:"pageNum"`
		PageSize int    `form:"pageSize"`
		Search   string `form:"search"`
		Type     string `form:"type"`
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
	q := db.DB.Model(&AccessControls)
	if req.Type == "username" {
		q = q.Where("username != ?", "")
		if req.Search != "" {
			q = q.Where(db.ClauseExprForContain(q, "username", req.Search))
		}
	} else if req.Type == "clientID" {
		q = q.Where("client_id != ?", "")
		if req.Search != "" {
			q = q.Where(db.ClauseExprForContain(q, "client_id", req.Search))
		}
	} else if req.Type == "clientIP" {
		q = q.Where("remote != ?", "")
		if req.Search != "" {
			q = q.Where(db.ClauseExprForContain(q, "remote", req.Search))
		}
	} else { // 全部用户
		q = q.Where("remote = ? and username = ? and client_id = ?", "", "", "")
	}
	if req.Search != "" {
		q = q.Where(db.ClauseExprForContain(q, "username", req.Search))
	}
	var result struct {
		Total int64 `json:"total"`
		Items []struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			ClientID    string `json:"clientID"`
			Remote      string `json:"remote"`
			Topic       string `json:"topic"`
			Access      string `json:"access"`
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

func createAcl(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		ClientID    string `json:"clientID"`
		Remote      string `json:"remote"`
		Topic       string `json:"topic"`
		Access      string `json:"access"`
		Description string `json:"description"`
	}
	if err := c.ShouldBind(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	acl := AccessControl{
		ID:          utils.UUIDWithoutDash(),
		Username:    req.Username,
		ClientID:    req.ClientID,
		Remote:      req.Remote,
		Topic:       req.Topic,
		Access:      req.Access,
		Description: req.Description,
	}
	err := db.DB.Create(&acl).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	clearCacheLater()
	web.SuccessResponse(c, acl.ID)
}

func updateAcl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	var req struct {
		Username    string `json:"username"`
		ClientID    string `json:"clientID"`
		Remote      string `json:"remote"`
		Topic       string `json:"topic"`
		Access      string `json:"access"`
		Description string `json:"description"`
	}
	if err := c.ShouldBind(&req); err != nil {
		web.BadRequestResponse(c, err.Error())
		return
	}
	acl := AccessControl{
		ID:          id,
		Username:    req.Username,
		ClientID:    req.ClientID,
		Remote:      req.Remote,
		Topic:       req.Topic,
		Access:      req.Access,
		Description: req.Description,
	}

	err := db.DB.Select("client_id", "username", "remote", "topic", "access").Save(&acl).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	clearCacheLater()
	web.SuccessResponse(c, nil)
}

func deleteAcl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.BadRequestResponse(c, "id is empty")
		return
	}
	defer clearCacheLater()
	err := db.DB.Delete(&AccessControl{}, "id=?", id).Error
	if err != nil {
		web.InternalErrorResponse(c, err)
		return
	}
	web.SuccessResponse(c, nil)
}

func Init() {
	router := web.Router
	api := router.Group("api/v1/acl", auth.Authorize)
	api.GET("", getAclRules)
	api.POST("", createAcl)
	api.PUT(":id", updateAcl)
	api.DELETE(":id", deleteAcl)
	//api.GET(":id/filters", getFilters)
	//api.PUT(":id/filters", updateFilters)
}
