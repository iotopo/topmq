package utils

import (
	"fmt"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		expected bool
	}{
		{"长度8", 8, true},
		{"长度12", 12, true},
		{"长度16", 16, true},
		{"长度20", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password := GeneratePassword(tt.length)
			
			// 检查密码长度
			if len(password) != tt.length {
				t.Errorf("GeneratePassword(%d) 生成的密码长度 = %d, 期望 %d", 
					tt.length, len(password), tt.length)
			}
			
			// 检查密码是否满足验证要求
			if !IsValidPassword(password) {
				t.Errorf("GeneratePassword(%d) 生成的密码 '%s' 不满足 IsValidPassword 要求", 
					tt.length, password)
			}
		})
	}
}

func TestGeneratePasswordConsistency(t *testing.T) {
	// 测试多次生成相同长度的密码，确保都能通过验证
	length := 12
	iterations := 100
	
	for i := 0; i < iterations; i++ {
		password := GeneratePassword(length)
		
		if len(password) != length {
			t.Errorf("第 %d 次生成: 密码长度 = %d, 期望 %d", i+1, len(password), length)
		}
		
		if !IsValidPassword(password) {
			t.Errorf("第 %d 次生成: 密码 '%s' 不满足验证要求", i+1, password)
		}
	}
}

func TestGeneratePasswordCharacterDistribution(t *testing.T) {
	password := GeneratePassword(1000) // 生成较长的密码以测试字符分布
	
	// 统计各种字符类型的数量
	letterCount := 0
	digitCount := 0
	specialCount := 0
	
	for _, char := range password {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			letterCount++
		} else if char >= '0' && char <= '9' {
			digitCount++
		} else {
			specialCount++
		}
	}
	
	// 确保每种类型都有字符
	if letterCount == 0 {
		t.Error("生成的密码中没有字母")
	}
	if digitCount == 0 {
		t.Error("生成的密码中没有数字")
	}
	if specialCount == 0 {
		t.Error("生成的密码中没有特殊字符")
	}
	
	t.Logf("字符分布: 字母=%d, 数字=%d, 特殊字符=%d", letterCount, digitCount, specialCount)
}

func TestGeneratePasswordEdgeCases(t *testing.T) {
	// 测试边界情况
	testCases := []struct {
		name     string
		length   int
		expected bool
	}{
		{"最小长度3", 3, true},  // 最小长度应该能包含所有三种字符类型
		{"长度4", 4, true},
		{"长度5", 5, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			password := GeneratePassword(tc.length)
			
			if len(password) != tc.length {
				t.Errorf("密码长度 = %d, 期望 %d", len(password), tc.length)
			}
			
			if !IsValidPassword(password) {
				t.Errorf("密码 '%s' 不满足验证要求", password)
			}
		})
	}
}

// 基准测试
func BenchmarkGeneratePassword(b *testing.B) {
	lengths := []int{8, 12, 16, 20}
	
	for _, length := range lengths {
		b.Run(fmt.Sprintf("长度%d", length), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GeneratePassword(length)
			}
		})
	}
}
