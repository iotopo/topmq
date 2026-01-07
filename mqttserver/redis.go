package mqttserver

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/pbkdf2"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/redis/go-redis/v9"
	"hash"
	"strconv"
	"strings"
)

type RedisAuth struct {
	//Hash         string `json:"hash,omitempty" redis:"hash"`
	Salt         string `json:"salt,omitempty" redis:"salt"`
	PasswordHash string `json:"passwordHash,omitempty" redis:"passwordHash"`
}

func (auth *RedisAuth) Verify(password string) bool {
	authHash := config.Conf.Redis.PasswordHash
	if authHash == "plain" {
		return auth.PasswordHash == password
	}

	passwordHash, err := hex.DecodeString(strings.ToLower(auth.PasswordHash))
	if err != nil {
		return false
	}

	if strings.HasPrefix(authHash, "pbkdf2") {
		// pbkdf2 with macfun iterations dklen
		// macfun: md5, sha1, sha256, sha512
		// 示例：pbkdf2,sha256,1000,20
		fields := strings.Split(authHash, ",")
		if len(fields) != 4 {
			return false
		}
		hashFun := fields[1]

		var h func() hash.Hash
		switch hashFun {
		case "md5":
			h = md5.New
		case "sha1":
			h = sha1.New
		case "sha256":
			h = sha256.New
		case "sha512":
			h = sha512.New
		default:
			return false
		}

		iter, err := strconv.Atoi(fields[2])
		keyLen, err := strconv.Atoi(fields[3])
		if iter == 0 {
			return false
		}
		if keyLen == 0 {
			return false
		}
		expected, err := pbkdf2.Key(h, password, []byte(auth.Salt), iter, keyLen)
		// NOTE: 此处已知的 iter, keyLength 不应返回 error
		// 若 len(salt) 过短，此处可能返回 error，但不应短路 ConstantTimeCompare
		return err == nil && subtle.ConstantTimeCompare(expected, passwordHash) > 0
	}

	switch authHash {
	case "sha256":
		h := sha256.New()
		s := h.Sum([]byte(password))
		return bytes.Equal(s, passwordHash)
	case "salt,sha256":
		h := sha256.New()
		if auth.Salt != "" {
			h.Write([]byte(auth.Salt))
		}
		s := h.Sum([]byte(password))
		return bytes.Equal(s, passwordHash)
	case "sha256,salt":
		h := sha256.New()
		if auth.Salt != "" {
			h.Write([]byte(password))
			s := h.Sum([]byte(auth.Salt))
			return bytes.Equal(s, passwordHash)
		} else {
			s := h.Sum([]byte(password))
			return bytes.Equal(s, passwordHash)
		}
	default:
		return false
	}
}

func GetRedisAuth(ctx context.Context, username string) (*RedisAuth, error) {
	key := fmt.Sprintf("%smqtt_user:%s", config.Conf.Redis.KeyPrefix, username)
	var auth RedisAuth

	if config.Conf.Redis.EMQ {
		err := cache.GetClient().HGetAll(ctx, key).Scan(&auth)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, nil
			} else {
				return nil, err
			}
		}
	} else {
		result, err := cache.GetClient().Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, nil
			} else {
				return nil, err
			}
		}
		if result == "" {
			return nil, nil
		}

		err = json.Unmarshal([]byte(result), &auth)
		if err != nil {
			return nil, err
		}
	}
	return &auth, nil
}
