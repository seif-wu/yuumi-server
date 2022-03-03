package validator

import (
	"fmt"
	"strings"

	v10 "github.com/go-playground/validator/v10"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"

	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Tran 默认语言
var Tran = "zh"

// 定义一个全局翻译器T
var ZhTran ut.Translator
var EnTran ut.Translator

func Trans(tran string) ut.Translator {
	switch tran {
	case "en":
		return EnTran
	case "zh":
		return ZhTran
	default:
		return EnTran
	}
}

func RegisterTranslator(tran string) (t ut.Translator, err error) {
	zhT := zh.New() // 中文翻译器
	enT := en.New() // 英文翻译器
	// 第一个参数是备用（fallback）的语言环境
	// 后面的参数是应该支持的语言环境（支持多个）
	uni := ut.New(enT, zhT, enT)
	t, ok := uni.GetTranslator(tran)
	if !ok {
		return t, fmt.Errorf("uni.GetTranslator(%s) failed", tran)
	}

	// 注册翻译器
	switch tran {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(instance, t)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(instance, t)
	default:
		err = enTranslations.RegisterDefaultTranslations(instance, t)
	}

	if err != nil {
		return t, err
	}

	return t, nil
}

// 注册多种语言
func initTrans() (err error) {
	t, err := RegisterTranslator("en")
	if err != nil {
		return err
	}
	EnTran = t

	t, err = RegisterTranslator("zh")
	if err != nil {
		return err
	}
	ZhTran = t

	return
}

// TranslateArr 翻译错误信息 - 数组模式
func TranslateArr(errs v10.ValidationErrors) (messages []string) {
	for _, v := range errs.Translate(Trans(Tran)) {
		messages = append(messages, v)
	}

	return messages
}

// TranslateStr 翻译错误信息 - 字符串模式
func TranslateStr(errs v10.ValidationErrors, sep string) string {
	return strings.Join(TranslateArr(errs), sep)
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe v10.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) v10.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}
