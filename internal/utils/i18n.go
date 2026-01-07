package utils

import (
	"reflect"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// 根据 cfg 提供的配置，用 localizer 加载一批翻译文本
// target 必须是一个结构体指针，每一个出现在 cfg 中的 public 的 string，会进行翻译
func LocalizeStruct(localizer *i18n.Localizer, cfg map[string]*i18n.LocalizeConfig, target any) (err error) {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		panic("LocalizeStruct: target must be a non-nil pointer")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		panic("LocalizeStruct: target must be a pointer to struct")
	}
	for _, field := range reflect.VisibleFields(val.Type()) {
		if field.Type.Kind() != reflect.String {
			continue
		}
		if fieldCfg, ok := cfg[field.Name]; ok {
			fieldVal, err := localizer.Localize(fieldCfg)
			if err != nil {
				return err
			}
			val.FieldByName(field.Name).SetString(fieldVal)
		}
	}
	return
}
