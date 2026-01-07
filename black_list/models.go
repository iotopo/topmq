package black_list

import (
	"context"
	"errors"
	redisCache "github.com/go-redis/cache/v9"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/db"
	"time"
)

// BlackList 认证规则， 用于定义黑名单
type BlackList struct {
	ID          string `gorm:"size:36;primaryKey"`
	ClientID    string `gorm:"size:100"` // 空表示不限制
	Username    string `gorm:"size:100"` // 空表示不限制
	Remote      string `gorm:"size:100"` // 空表示不限制
	Allow       bool   // true 表示白名单，false 表示黑名单
	Description string `gorm:"type:text"`
	ExpiredAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var BlackLists = new(BlackList)

var blackListCache = cache.NewCache(&cache.CacheOption{
	LocalCache:    true,
	LocalCacheTTL: time.Minute * 10,
})

const blackListCacheKey = "mqtt:black_list"

func FindAll(ctx context.Context) ([]BlackList, error) {
	var auths []BlackList

	// 先从缓存中获取
	if err := blackListCache.Get(ctx, config.Conf.Redis.KeyPrefix+blackListCacheKey, &auths); err != nil {
		if !errors.Is(err, redisCache.ErrCacheMiss) {
			return nil, err
		} // else 未找到缓存，从数据库中获取
	} else {
		return auths, nil
	}

	if err := db.DB.Find(&auths).Error; err != nil {
		return nil, err
	}
	_ = blackListCache.Set(&redisCache.Item{Ctx: ctx, Key: config.Conf.Redis.KeyPrefix + blackListCacheKey, Value: &auths, TTL: 3 * 24 * time.Hour})
	return auths, nil
}

func clearCacheLater() {
	_ = blackListCache.Delete(context.Background(), config.Conf.Redis.KeyPrefix+blackListCacheKey)
}
