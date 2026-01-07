package utils

import (
	"fmt"
	"regexp"
)

func SplitProps(prefix, expression string) (props []string) {
	// 存储提取出的点位和属性，避免重复
	m := make(map[string]bool)

	// 使用正则表达式同时匹配四种格式：
	// 1. data["xxx"] - 双引号格式
	// 2. data['xxx'] - 单引号格式
	// 3. data.xxx - 点号格式（xxx可以包含字母、数字和下划线）
	dataRegex := regexp.MustCompile(fmt.Sprintf(`%s\["([^"]+)"]|%s\['([^']+)']|%s\.([a-zA-Z0-9_]+)`, prefix, prefix, prefix))
	dataMatches := dataRegex.FindAllStringSubmatch(expression, -1)

	// 处理data相关格式
	for _, match := range dataMatches {
		// 检查三个捕获组，分别对应双引号、单引号和点号格式
		for i := 1; i <= 3; i++ {
			if match[i] != "" {
				m[match[i]] = true // true表示是point
				break
			}
		}
	}

	// 生成返回结果
	for k, _ := range m {
		props = append(props, k)
	}

	return
}
