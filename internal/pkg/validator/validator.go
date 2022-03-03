package validator

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"yuumi/internal/pkg/util/funcutil"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type TransMessage struct {
	Tran    string
	Message string
}

var instance *validator.Validate
var mu sync.Mutex

func init() {
	if instance == nil {
		mu.Lock()
		if instance == nil {
			ok := true
			instance, ok = binding.Validator.Engine().(*validator.Validate)

			if !ok {
				panic(fmt.Errorf("bind validator engine error"))
			}
			log.Println("初始化验证器")

			// 注册一个获取json tag的自定义方法
			instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get(Tran), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})

			err := initTrans()
			if err != nil {
				panic(fmt.Errorf("init translator error: %s", err.Error()))
			}
		}
		mu.Unlock()
	}
}

// RegisterValidation 注册验证函数
func RegisterValidation(tagName string, fn validator.Func) (err error) {
	err = instance.RegisterValidation(tagName, fn)
	if err != nil {
		return err
	}
	return nil
}

func RegisterTranslation(tag string, trans ut.Translator, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc) (err error) {
	if err := instance.RegisterTranslation(tag, trans, registerFn, translationFn); err != nil {
		return err
	}
	return
}

// RegisterValidationAndTran 注册验证函数以及添加翻译
func RegisterValidationAndTran(fn validator.Func, transMessage []TransMessage) (err error) {
	tagName := funcutil.GetName(fn, '/', '.')
	if err = RegisterValidation(tagName, fn); err != nil {
		return err
	}

	for _, v := range transMessage {
		if err = RegisterTranslation(
			tagName,
			Trans(v.Tran),
			registerTranslator(tagName, v.Message),
			translate,
		); err != nil {
			return err
		}
	}

	return
}
