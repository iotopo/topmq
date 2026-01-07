package validators

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

type RegexValidator struct {
	Name  string
	Regex *regexp.Regexp
}

func (v RegexValidator) Provide() RegexValidator {
	return v
}

// 按指定名称提供一个 gin/binding.Validator 使用的正则验证
// 需要主程序调用 RegisterRegexValidators 才会生效
func ProvideRegex(name string, regex *regexp.Regexp) fx.Option {
	provided := RegexValidator{
		Name:  name,
		Regex: regex,
	}
	return fx.Provide(
		fx.Annotate(
			provided.Provide,
			fx.ResultTags(`group:"regexValidators"`),
		),
	)
}

type RegexValidation map[string]*regexp.Regexp

// 按指定名称查找 Regexp 验证
// 需要先调用 RegisterValidators 注册
// 用例：
//
//	type Request struct {
//		PointID `binding:"required,regex=PointID"`
//	}
//
// 1. 此验证只对字符串类型生效，对其他类型无效
// 2. 空字符串不进行验证，如不能为空，可搭配 required 验证
// 3. 未指定名称时，默认与字段名相同；指定的名称无效时，认证失败
func (c RegexValidation) Validate(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		if val == "" {
			return true
		}

		regexName := fl.Param()
		if regexName == "" {
			regexName = fl.FieldName()
		}

		regex := c[regexName]
		if regex == nil {
			return false
		}

		if !regex.MatchString(val) {
			return false
		}
	}
	return true
}

type RegisterRegexValidationDependencies struct {
	fx.In

	Validators []RegexValidator `group:"regexValidators"`
}

func (deps *RegisterRegexValidationDependencies) BuildValidation() RegexValidation {
	validators := make(map[string]*regexp.Regexp, len(deps.Validators))
	for _, validator := range deps.Validators {
		validators[validator.Name] = validator.Regex
	}
	return validators
}

// 向 gin/binding.Validator 注册一个名为 regex 的验证
// 在此之前，应使用 ProvideRegex 注册所有表达式
func RegisterRegexValidation(deps RegisterRegexValidationDependencies) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		vc := deps.BuildValidation()
		v.RegisterValidation("regex", vc.Validate)
	}
}
