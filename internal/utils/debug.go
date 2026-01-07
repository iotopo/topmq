package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func DebugBytesWithBuffer(buf *strings.Builder, data []byte) {
	dataLen := len(data)
	if dataLen == 0 {
		_, _ = buf.WriteString("(0 bytes)")
		return
	}
	blocks := dataLen / 16
	rest := dataLen % 16
	if rest > 0 {
		blocks++
	}
	// 由于 dataLen == 0 已 return，必有 blocks > 0
	// 最后一行的起始地址是 (blocks - 1) * 16，用其计算起始地址的对其宽度
	addrWidth := len(strconv.FormatInt(int64(blocks-1)*16, 16))
	// 上边框
	for i := addrWidth + 2; i > 0; i-- {
		_ = buf.WriteByte(' ')
	}
	_, _ = buf.WriteString("|-------------|-------------|-------------|-------------|\n")
	// 数据主体
	for addr, b := range data {
		if addr%16 == 0 {
			if addr > 0 {
				_ = buf.WriteByte('\n')
			}
			// 起始地址
			addrStr := strconv.FormatInt(int64(addr), 16)
			for i := addrWidth - len(addrStr); i > 0; i-- {
				_ = buf.WriteByte('0')
			}
			_, _ = buf.WriteString(addrStr)
			_, _ = buf.WriteString(": | ")
		} else {
			_ = buf.WriteByte(' ')
		}
		_, _ = fmt.Fprintf(buf, "%02x", b)
		// 右侧边框
		switch addr % 16 {
		case 3, 7, 11, 15:
			_, _ = buf.WriteString(" |")
		}
	}
	// 如果最后一行没有填满，补充
	if rest > 0 {
		cells := rest / 4
		cellRest := rest % 4
		// 先填满一个格子
		if cellRest > 0 {
			cells++
			for i := (4 - cellRest) * 3; i > 0; i-- {
				_ = buf.WriteByte(' ')
			}
			_, _ = buf.WriteString(" |")
		}
		// 填满其余空白格子
		for i := cells; i < 4; i++ {
			for j := 0; j < 13; j++ {
				_ = buf.WriteByte(' ')
			}
			_ = buf.WriteByte('|')
		}
	}
	_ = buf.WriteByte('\n')
	// 下边框
	for i := addrWidth + 2; i > 0; i-- {
		_ = buf.WriteByte(' ')
	}
	_, _ = buf.WriteString("|-------------|-------------|-------------|-------------|")
	// 附加信息
	_, _ = fmt.Fprintf(buf, "\n(%d bytes)", dataLen)
}

// 将字节串格式化为每行 16 个字节的美观格式，供调试用
//
//goland:noinspection GoUnusedExportedFunction
func DebugBytes(data []byte) string {
	var buf strings.Builder
	DebugBytesWithBuffer(&buf, data)
	return buf.String()
}
