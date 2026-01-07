package utils

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/subtle"
	"math/rand"
	"regexp"
)

/**
 * WARNING
 * 由于计算机硬件算力的发展，pbkdf2 和 bcrypt 等密钥生成算法已经逐渐变得不安全
 * 在新的 PKCS 标准出版之前，目前比较流行的算法是 scrypt 和 argon2
**/
const PasswordSaltLength = 16
const PasswordHashLength = 32
const PasswordIterations = 10_0000

//goland:noinspection GoUnusedExportedFunction
func GeneratePasswordHash(password string) ([]byte, []byte) {
	salt := RandomBytes(PasswordSaltLength)
	hash, err := pbkdf2.Key(sha256.New, password, salt, PasswordIterations, PasswordHashLength)
	if err != nil {
		// SAFETY: 此处已知的 len(salt), iter, keyLength 不应返回 error
		panic(err)
	}
	return hash, salt
}

//goland:noinspection GoUnusedExportedFunction
func CheckPasswordHash(password string, hash, salt []byte) bool {
	expected, err := pbkdf2.Key(sha256.New, password, salt, PasswordIterations, PasswordHashLength)
	// NOTE: 此处已知的 iter, keyLength 不应返回 error
	// 若 len(salt) 过短，此处可能返回 error，但不应短路 ConstantTimeCompare
	return err == nil && subtle.ConstantTimeCompare(expected, hash) > 0
}

// IsValidPassword 校验密码是否符合要求
// 至少包含字母、数字和1个特殊字符
func IsValidPassword(password string) bool {
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]`).MatchString(password)

	return hasLetter && hasDigit && hasSpecial
}

func GeneratePassword(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "!@#$%^&*()_+-=[]{};':\"\\|,.<>/?"

	// 确保每种类型至少有一个字符
	password := make([]byte, length)

	// 强制包含每种类型
	password[0] = letters[rand.Intn(len(letters))]
	password[1] = digits[rand.Intn(len(digits))]
	password[2] = special[rand.Intn(len(special))]

	// 剩余位置随机填充
	all := letters + digits + special
	for i := 3; i < length; i++ {
		password[i] = all[rand.Intn(len(all))]
	}

	// 打乱密码顺序
	for i := len(password) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}
