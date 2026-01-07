package auth

import (
	"context"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/db/datatypes"
	"github.com/iotopo/topmq/internal/utils"
	"time"

	redisCache "github.com/go-redis/cache/v9"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

//	type UserService struct {
//		ctx context.Context
//	}
type User struct {
	ID           string         `gorm:"size:36;primaryKey"`
	Username     string         `gorm:"size:20;not null"`
	PasswordHash datatypes.Blob `gorm:"not null"`
	PasswordSalt datatypes.Blob `gorm:"not null"`
	Name         string         `gorm:"size:100;not null;index"`
	Locked       bool           `json:",omitempty" gorm:"not null;default:false;index"`
	IsSuperuser  bool           `json:"isSuperuser"`

	Email  *string `json:",omitempty" gorm:"size:64;index"`
	Tel    *string `json:",omitempty" gorm:"size:64;index"`
	Mobile *string `json:",omitempty" gorm:"size:64;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var Users User

type UserForSelect struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserForAuthenticate struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"user_name"`
	Locked    bool   `json:"locked,omitempty"`
	UpdatedAt time.Time
}

type UserSession struct {
	UserID    string
	UpdatedAt time.Time
	ClientIP  string
	//AppID      string
}

type LoginUser struct {
	UID         string
	Name        string
	Username    string
	IsSuperuser bool
}

func (u *LoginUser) IsAnonymous() bool {
	return u.UID == ""
}

var userForAuthenticateCache = cache.NewCache(&cache.CacheOption{
	LocalCache: true,
})

func DelUserAuthCache(ctx context.Context, userID string) {
	cacheKey := UserForAuthenticateV2(userID)
	err := userForAuthenticateCache.Delete(ctx, cacheKey)
	if err != nil {
		logrus.WithError(err).Errorf("delete user for authenticate (%s) from cache failed", userID)
	}
}

func GetForAuthenticate(ctx context.Context, userID string) (*UserForAuthenticate, error) {
	cacheKey := UserForAuthenticateV2(userID)
	var user UserForAuthenticate
	err := userForAuthenticateCache.Get(ctx, cacheKey, &user)
	if err != nil {
		if !errors.Is(err, redisCache.ErrCacheMiss) {
			logrus.WithError(err).Errorf("get user for authenticate (%s) from cache failed", userID)
		}
	} else {
		if user.ID == "" {
			return nil, nil
		}
		return &user, nil
	}

	selectFrom := clause.Table{Name: "users", Alias: "u"}
	conditions := []clause.Expression{
		clause.Eq{
			Column: clause.Column{Table: "u", Name: "id"},
			Value:  userID,
		},
	}
	columns := []clause.Column{
		{Table: "u", Name: "id"},
		{Table: "u", Name: "name"},
		{Table: "u", Name: "username"},
		{Table: "u", Name: "locked"},
		{Table: "u", Name: "updated_at"},
	}
	result := db.DB.Table("?", selectFrom).
		Clauses(clause.Where{Exprs: conditions}).
		Clauses(clause.Select{Columns: columns}).
		Limit(1).
		Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		err = userForAuthenticateCache.Set(&redisCache.Item{Ctx: context.Background(), Key: cacheKey, Value: &UserForAuthenticate{}, TTL: 15 * time.Minute})
		if err != nil {
			logrus.WithError(err).Errorf("save user for authenticate (%s) to cache failed", userID)
		}
		return nil, nil
	}

	err = userForAuthenticateCache.Set(&redisCache.Item{Ctx: context.Background(), Key: cacheKey, Value: &user, TTL: 15 * time.Minute})
	if err != nil {
		logrus.WithError(err).Errorf("save user for authenticate (%s) to cache failed", userID)
	}
	//if cached, err := json.Marshal(&user); err == nil {
	//	if err := redisClient.Set(context.Background(), redis_keys.UserForAuthenticate(userID), cached, 24*time.Hour).Err(); err != nil {
	//		logrus.Debugf("failed setting cached UserForAuthenticate: %v", err)
	//	}
	//}
	return &user, nil
}

func InitDefaultUser() error {
	logrus.Info("setting up admin user")
	conn := db.DB
	su := User{IsSuperuser: true}
	result := conn.Where(&su, "is_superuser").Limit(1).Find(&su)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		username := "admin"
		password := "admin"

		su.ID = utils.XID()
		su.Username = username
		su.Name = "管理员"
		su.PasswordHash, su.PasswordSalt = utils.GeneratePasswordHash(password)
		su.IsSuperuser = true
		if err := conn.Create(&su).Error; err != nil {
			return err
		}
	}
	return nil
}
