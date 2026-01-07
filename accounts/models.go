package accounts

import (
	"context"
	"errors"
	"fmt"
	redisCache "github.com/go-redis/cache/v9"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/web/errcodes"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID          string `gorm:"size:36;primaryKey"`
	Username    string `gorm:"size:20;not null;uniqueIndex"`
	Password    string `gorm:"size:20;not null"`
	ClientID    string `gorm:"size:100"` // 空表示不限制
	Remote      string `gorm:"size:100"` // 空表示不限制
	Description string `gorm:"type:text"`
	Disabled    bool   `gorm:"not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	//Filters     map[string]string `gorm:"serializer:json"`
}

var Accounts = new(Account)

var accountCache = cache.NewCache(&cache.CacheOption{
	LocalCache:    true,
	LocalCacheTTL: time.Minute * 10,
})

func getUserCacheKey(username string) string {
	return fmt.Sprintf("%smqtt:account:%s", config.Conf.Redis.KeyPrefix, username)
}

//func (*Account) UpdateFilter(id string, filter map[string]string) error {
//	return Accounts.runAndClearCache(id, func() error {
//		return db.DB.Model(&Account{}).Where("id=?", id).Update("filters", filter).Error
//	})
//}

func runAndClearCache(id string, fn func() error) error {
	var user Account
	db.DB.Select("username").Where("id = ?", id).Limit(1).Find(&user)
	err := fn()
	if err == nil && user.Username != "" {
		_ = accountCache.Delete(context.Background(), getUserCacheKey(user.Username))
	}
	return err
}

func FindByUsername(username string) (*Account, error) {
	cacheKey := getUserCacheKey(username)

	var user Account

	// 先从缓存中获取
	if err := accountCache.Get(context.Background(), cacheKey, &user); err != nil {
		if !errors.Is(err, redisCache.ErrCacheMiss) {
			return nil, err
		} // else 未找到缓存，从数据库中获取
	} else {
		if user.ID == "" {
			return nil, errcodes.NotFound
		}
		return &user, nil
	}

	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = accountCache.Set(&redisCache.Item{Ctx: context.Background(), Key: cacheKey, Value: &user, TTL: 3 * 24 * time.Hour})
			return nil, nil
		}
		return nil, err
	}
	_ = accountCache.Set(&redisCache.Item{Ctx: context.Background(), Key: cacheKey, Value: &user, TTL: 3 * 24 * time.Hour})
	return &user, nil
}

func FindByPage(search string, pageNum, pageSize int) []*Account {
	var users []*Account
	q := db.DB
	if search != "" {
		q = q.Where(db.ClauseExprForContain(q, "username", search))
	}
	q.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users)
	return users
}
