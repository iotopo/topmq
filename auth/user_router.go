package auth

import (
	"context"
	"database/sql"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web"
	"github.com/iotopo/topmq/web/errcodes"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var defaultPage = 1
var defaultPageSize = 10

func getPaged(c *gin.Context) {
	type QueryParams struct {
		Page     *int   `form:"pageNum"`
		PageSize *int   `form:"pageSize"`
		Username string `form:"username"`
		Name     string `form:"name"`
		Email    string `form:"email"`
		Tel      string `form:"tel"`
		Mobile   string `form:"mobile"`
		Locked   string `form:"locked"`
	}

	type ResponseUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Tel      string `json:"tel"`
		Mobile   string `json:"mobile"`
		Locked   bool   `json:"locked"`
	}

	columns := []string{"id", "username", "name", "email", "tel", "mobile", "locked", "is_superuser"}

	type Response struct {
		Total int64          `json:"total"`
		Users []ResponseUser `json:"users"`
	}

	var req QueryParams
	if err := c.BindQuery(&req); err != nil {
		return
	}

	if req.Page == nil {
		req.Page = &defaultPage
	} else if *req.Page <= 0 {
		c.String(http.StatusBadRequest, "requires pageNum > 0")
		return
	}
	if req.PageSize == nil {
		req.PageSize = &defaultPageSize
	} else if *req.PageSize <= 0 {
		c.String(http.StatusBadRequest, "requires pageSize > 0")
		return
	} else if *req.PageSize > 1000 {
		c.String(http.StatusBadRequest, "requires pageSize <= 1000")
		return
	}

	query := db.DB.Model(&Users).
		Where(clause.Eq{Column: "is_superuser", Value: false})

	if req.Name != "" {
		query = query.Where(db.ClauseExprForContain(db.DB, "name", req.Name))
	}
	if req.Username != "" {
		query = query.Where(db.ClauseExprForContain(db.DB, "username", req.Username))
	}
	if req.Email != "" {
		query = query.Where(db.ClauseExprForContain(db.DB, "email", req.Email))
	}
	if req.Tel != "" {
		query = query.Where(db.ClauseExprForContain(db.DB, "tel", req.Tel))
	}
	if req.Mobile != "" {
		query = query.Where(db.ClauseExprForContain(db.DB, "mobile", req.Mobile))
	}
	switch req.Locked {
	case "true", "1":
		query = query.Where(clause.Eq{Column: "locked", Value: true})
	case "false", "0":
		query = query.Where(clause.Eq{Column: "locked", Value: false})
	}
	var res Response
	if err := query.Count(&res.Total).Error; err != nil {
		web.HandleError(c, err)
		return
	}
	if res.Total > 0 {
		query = query.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{{
				Column: clause.Column{Name: "id"},
				Desc:   true,
			}},
		})
		query = query.Select(columns).Limit(*req.PageSize).Offset(*req.PageSize * (*req.Page - 1))
		if err := query.Find(&res.Users).Error; err != nil {
			web.HandleError(c, err)
			return
		}
	}
	web.SuccessResponse(c, res)
}

func usernameTest(c *gin.Context) {
	type Response struct {
		OK bool `json:"ok"`
	}
	username := c.Query("username")
	if username == "" {
		web.SuccessResponse(c, Response{OK: false})
		return
	}
	subquery := db.DB.Model(&Users).Where(clause.Eq{Column: "username", Value: username})
	var err error
	var exists bool
	if db.DB.Dialector.Name() == "dm" {
		var innerExist sql.NullBool
		err = db.DB.Raw("SELECT ? WHERE EXISTS (?)", true, subquery).Scan(&innerExist).Error
		exists = err == nil && innerExist.Valid && innerExist.Bool
	} else {
		err = db.DB.Raw("SELECT EXISTS (?)", subquery).Scan(&exists).Error
	}
	if err != nil {
		web.HandleError(c, err)
		return
	}
	web.SuccessResponse(c, Response{OK: !exists})
}

func create(c *gin.Context) {
	type RequestBody struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Tel      string `json:"tel"`
		Mobile   string `json:"mobile"`
	}

	type Response struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	var req RequestBody
	if err := c.BindJSON(&req); err != nil {
		return
	}

	// token := authorization.GetLoginUser(c)

	user := User{
		ID:       utils.XID(),
		Username: req.Username,
		Name:     req.Name,
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Tel != "" {
		user.Tel = &req.Tel
	}
	if req.Mobile != "" {
		user.Mobile = &req.Mobile
	}
	password := utils.GeneratePassword(config.Conf.MinPwdLen)
	user.PasswordHash, user.PasswordSalt = utils.GeneratePasswordHash(password)
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// tenant := auth.Tenant{ID: user.TenantID}
		// result := tx.Model(&tenant).
		//	Clauses(clause.Locking{Strength: "UPDATE"}).
		//	Select([]string{"user_count", "user_limit"}).
		//	Limit(1).
		//	Find(&tenant)
		// if result.Error != nil {
		//	return result.Error
		// }
		// if result.RowsAffected == 0 {
		//	return failure.SimpleUnauthorized
		// }
		// if tenant.UserLimit != nil && tenant.UserCount >= *tenant.UserLimit {
		//	return errcodes.UserCountLimit
		// }

		err := tx.Create(&user).Error
		if err != nil {
			return err
		}

		// tenant.UserCount += 1
		// if err := tx.Model(&tenant).Update("user_count", tenant.UserCount).Error; err != nil {
		//	return err
		// }
		return nil
	})
	if err != nil {
		web.HandleError(c, err)
		return
	}

	res := Response{
		ID:       user.ID,
		Password: password,
	}
	web.SuccessResponse(c, res)
}

func modify(c *gin.Context) {
	type RequestBody struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Tel    string `json:"tel"`
		Mobile string `json:"mobile"`
	}

	columns := []string{"name", "email", "tel", "mobile"}

	var req RequestBody
	if err := c.BindJSON(&req); err != nil {
		return
	}

	user := User{
		ID:   c.Param("id"),
		Name: req.Name,
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Tel != "" {
		user.Tel = &req.Tel
	}
	if req.Mobile != "" {
		user.Mobile = &req.Mobile
	}

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
		return tx.Select(columns).Save(&user).Error
	})

	if err != nil {
		web.HandleError(c, err)
		return
	}

	web.SuccessResponse(c, nil)
}

func lockOrUnlock(c *gin.Context, locked bool) {
	user := User{ID: c.Param("id")}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select([]string{"username", "locked"}).
			Limit(1).
			Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errcodes.NotFound
		}
		if user.Locked != locked {
			return tx.Model(&user).Update("locked", locked).Error
		} else {
			return nil
		}
	})
	if err != nil {
		web.HandleError(c, err)
		return
	}

	if locked {
		ClearUserSession(context.Background(), user.ID)
	} else {
		DelUserAuthCache(context.Background(), user.ID)
	}

	web.SuccessResponse(c, nil)
}

func lock(c *gin.Context) {
	lockOrUnlock(c, true)
}

func unlock(c *gin.Context) {
	lockOrUnlock(c, false)
}

func resetPassword(c *gin.Context) {
	type Response struct {
		Password string `json:"password"`
	}

	user := User{ID: c.Param("id")}
	password := utils.GeneratePassword(config.Conf.MinPwdLen)
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select([]string{"username"}).
			Limit(1).
			Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errcodes.NotFound
		}

		user.PasswordHash, user.PasswordSalt = utils.GeneratePasswordHash(password)
		return tx.Select([]string{"password_hash", "password_salt", "updated_at"}).Save(&user).Error
	})
	if err != nil {
		web.HandleError(c, err)
		return
	}

	ClearUserSession(context.Background(), user.ID)

	web.SuccessResponse(c, Response{Password: password})
}

func delete_(c *gin.Context) {
	user := User{ID: c.Param("id")}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&user).
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

		err := tx.Delete(&user).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		web.HandleError(c, err)
		return
	}
	web.SuccessResponse(c, nil)
	ClearUserSession(context.Background(), user.ID)
}
