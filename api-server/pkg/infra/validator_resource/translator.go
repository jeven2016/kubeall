package validator_resource

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
)

type ValidatorTranslator interface {
	Zh() ut.Translator
	En() ut.Translator
}

type validatorTranslatorImpl struct {
	validator    *validator.Validate
	zhTranslator ut.Translator
	enTranslator ut.Translator
}

func NewValidatorTranslator() ValidatorTranslator {
	impl := &validatorTranslatorImpl{}
	impl.initTranslators()
	return impl
}

func (v *validatorTranslatorImpl) initTranslators() {
	enLang := en.New()   //英文翻译器
	zh := zhongwen.New() //中文翻译器

	// default language, optional languages,
	// uni := ut.New(en, en)
	// uni := ut.New(en, zh, tw)
	uni := ut.New(zh, zh, enLang)

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if zhTranslator, found := uni.GetTranslator("zh"); found {
			v.zhTranslator = zhTranslator
			if err := zhtrans.RegisterDefaultTranslations(val, v.zhTranslator); err != nil {
				zap.L().Warn("un expected error during registering chinese translator", zap.Error(err))
			}
		}

		if enTranslator, found := uni.GetTranslator("en"); found {
			v.enTranslator = enTranslator
			if err := entrans.RegisterDefaultTranslations(val, v.enTranslator); err != nil {
				zap.L().Warn("un expected error during registering english translator", zap.Error(err))
			}
		}
	} else {
		zap.L().Warn("failed to get validator engine")
	}
}

func (v *validatorTranslatorImpl) Zh() ut.Translator {
	return v.zhTranslator
}

func (v *validatorTranslatorImpl) En() ut.Translator {
	return v.enTranslator
}
