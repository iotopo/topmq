package utils

import (
	cryptoRand "crypto/rand"
	"io"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rs/xid"
)

//goland:noinspection GoUnusedExportedFunction
func UUID() string {
	id, err := uuid.NewV7()
	if err == nil {
		return id.String()
	}
	return uuid.New().String()
}

//goland:noinspection GoUnusedExportedFunction
func UUIDWithoutDash() string {
	return strings.ReplaceAll(UUID(), "-", "")
}

//goland:noinspection GoUnusedExportedFunction
func XID() string {
	return xid.New().String()
}

//goland:noinspection GoUnusedExportedFunction
func ULID() string {
	return ulid.Make().String()
}

func RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, _ = io.ReadFull(cryptoRand.Reader, bytes)
	return bytes
}

// 字母和数字，移除了 oODlL01 等易混淆字母
const randomStringLetters = "abcdefghijkmnpqrstuvwxyzABCEFGHIJKMNPQRSTUVWXYZ23456789"

// RandomString returns a random string with a fixed length
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randomStringLetters[rand.Intn(len(randomStringLetters))]
	}

	return string(b)
}
