package acls

import (
	"context"
	"errors"
	redisCache "github.com/go-redis/cache/v9"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/db"
	"time"
)

// AccessControl 主题授权规则
type AccessControl struct {
	ID          string `gorm:"size:36;primaryKey"`
	ClientID    string `gorm:"size:100"`          // 空表示不限制
	Username    string `gorm:"size:100"`          // 空表示不限制
	Remote      string `gorm:"size:100"`          // 空表示不限制
	Topic       string `gorm:"size:200;not null"` // 主题
	Access      string `gorm:"size:20"`           // d(禁用) r(只读) w(只写) rw(读写)
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var AccessControls = new(AccessControl)

var aclCache = cache.NewCache(&cache.CacheOption{
	LocalCache:    true,
	LocalCacheTTL: time.Minute * 10,
})

const aclCacheKey = "mqtt:acl"

func FindAll(ctx context.Context) ([]AccessControl, error) {
	var items []AccessControl

	// 先从缓存中获取
	if err := aclCache.Get(ctx, config.Conf.Redis.KeyPrefix+aclCacheKey, &items); err != nil {
		if !errors.Is(err, redisCache.ErrCacheMiss) {
			return nil, err
		} // else 未找到缓存，从数据库中获取
	} else {
		return items, nil
	}

	if err := db.DB.Where("topic != ?", "").Find(&items).Error; err != nil {
		return nil, err
	}
	_ = aclCache.Set(&redisCache.Item{Ctx: ctx, Key: config.Conf.Redis.KeyPrefix + aclCacheKey, Value: &items, TTL: 3 * 24 * time.Hour})
	return items, nil
}

func (*AccessControl) FindByPage(search string, pageNum, pageSize int) []*AccessControl {
	var items []*AccessControl
	q := db.DB
	if search != "" {
		q = q.Where(db.ClauseExprForContain(q, "topic", search))
	}
	q.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&items)
	return items
}

func clearCacheLater() {
	_ = aclCache.Delete(context.Background(), config.Conf.Redis.KeyPrefix+aclCacheKey)
}
