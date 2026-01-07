package secrets

import (
	"strings"
	"testing"
)

func TestEncodeEmptyString(t *testing.T) {
	// 测试空字符串加密
	result, err := Encode("")
	if err != nil {
		t.Errorf("Encode empty string should not return error: %v", err)
	}
	if result != "" {
		t.Errorf("Encode empty string should return empty string, got: %q", result)
	}
}

func TestDecodeEmptyString(t *testing.T) {
	// 测试空字符串解密
	result, err := Decode("")
	if err != nil {
		t.Errorf("Decode empty string should not return error: %v", err)
	}
	if result != "" {
		t.Errorf("Decode empty string should return empty string, got: %q", result)
	}
}

func TestDecodeNonEncrypted(t *testing.T) {
	// 测试解密非加密文本
	testCases := []string{
		"plain text",
		"test@123!@#$%",
		"测试中文",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			result, err := Decode(tc)
			if err != nil {
				t.Errorf("Decode() returned error for non-encrypted text: %v", err)
			}
			if result != tc {
				t.Errorf("Decode() = %v, want %v", result, tc)
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	// 测试加密解密
	testCases := []string{
		"Tq7@nK2$wX9",
		"hello world",
		"test@123!@#$%",
		"测试中文",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			// 加密
			encoded, err := Encode(tc)
			if err != nil {
				t.Errorf("Encode() error: %v", err)
				return
			}

			// 验证加密结果有正确的头部
			if !strings.HasPrefix(encoded, header) {
				t.Errorf("Encode() result = %v, should start with %v", encoded, header)
			}

			// 解密
			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Decode() error: %v", err)
				return
			}

			if decoded != tc {
				t.Errorf("Decode() = %v, want %v", decoded, tc)
			}
		})
	}
}

func TestDecodeInvalidHeader(t *testing.T) {
	// 测试解密只有头部但没有数据的文本
	_, err := Decode("ENCv1|")
	if err == nil {
		t.Error("Decode() should return error for invalid encrypted text")
	}
}
