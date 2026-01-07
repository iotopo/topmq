package cache

import "github.com/redis/go-redis/v9"

func IsAuthError(err error) bool {
	return redis.HasErrorPrefix(err, "ERR invalid password") ||
		redis.HasErrorPrefix(err, "WRONGPASS") ||
		redis.HasErrorPrefix(err, "NOAUTH")
}
