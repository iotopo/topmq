package mqttserver

import (
	"regexp"
	"strings"
)

// 正则匹配 ${变量名} 格式（变量名支持字母、数字、下划线）
var re = regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)}`)

// ReplaceVariables 替换字符串中的 ${变量名} 为映射中的值
// 参数：
//   - str: 待替换的原始字符串
//   - vars: 变量映射表（key=变量名，value=替换值）
//
// 返回：替换后的字符串
func ReplaceVariables(str string, vars map[string]string) string {

	// 替换逻辑：遍历所有匹配项，从vars中取值替换
	result := re.ReplaceAllStringFunc(str, func(match string) string {
		// 提取 ${} 中的变量名（去掉 ${ 和 }）
		varName := strings.TrimPrefix(match, "${")
		varName = strings.TrimSuffix(varName, "}")

		// 从映射中取值，不存在则保留原占位符
		if val, ok := vars[varName]; ok {
			return val
		}
		return match // 未找到变量时，保留原始 ${变量名}
	})

	return result
}
