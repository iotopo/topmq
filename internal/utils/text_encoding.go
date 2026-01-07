package utils

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

type TextEncoding string

const (
	TextEncodingUndefined TextEncoding = ""
	TextEncodingUTF8      TextEncoding = "utf8"
	TextEncodingUTF16     TextEncoding = "utf16" // deprecated
	TextEncodingUTF16Le   TextEncoding = "utf16le"
	TextEncodingUTF16Be   TextEncoding = "utf16be"
	TextEncodingGB2312    TextEncoding = "gb2312" // deprecated
	TextEncodingGBK       TextEncoding = "gbk"
	TextEncodingGB18030   TextEncoding = "gb18030"
)

// 根据 src 指定的编码解析 src 为字符串
// 返回 error 为空时，保证 string 为 UTF-8 编码的
func DecodeText(src []byte, enc TextEncoding) (string, error) {
	var codec encoding.Encoding

	switch enc {
	case TextEncodingGB2312:
		// deprecated: GBK 完全兼容 GB2312，建议直接使用 GBK
		codec = simplifiedchinese.HZGB2312
	case TextEncodingGBK:
		codec = simplifiedchinese.GBK
	case TextEncodingGB18030:
		codec = simplifiedchinese.GB18030
	case TextEncodingUTF16:
		// deprecated: 根据 BOM 决定，无 BOM 时报错，建议使用 UTF16Le 或 UTF16Be
		codec = unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	case TextEncodingUTF16Le:
		// 有 BOM 时根据 BOM 决定，无 BOM 时默认
		codec = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case TextEncodingUTF16Be:
		// 有 BOM 时根据 BOM 决定，无 BOM 时默认 BE
		codec = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case TextEncodingUndefined, TextEncodingUTF8:
		// 虽然在 Go 中，所有 string 字面量默认是 UTF-8 编码
		// 但将非 UTF-8 编码的 []byte 转换为 string 并不会进行检查
		// 换言之，string 可以是非 UTF-8 编码的
		// 因此这里显式进行检查，并允许字符串开头带有 UTF-8 BOM
		codec = unicode.UTF8BOM
	default:
		return "", fmt.Errorf("unknown encoding %q", enc)
	}

	utf8Bytes, err := codec.NewDecoder().Bytes(src)
	if err != nil {
		return "", fmt.Errorf("failed decoding string in %s: %v", enc, err)
	}

	return string(utf8Bytes), nil
}
