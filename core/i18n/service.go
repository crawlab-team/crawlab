package i18n

import "github.com/crawlab-team/crawlab/core/interfaces"

var translations []interfaces.Translation

var _svc interfaces.I18nService

type Service struct {
}

func (svc *Service) AddTranslations(t []interfaces.Translation) {
	translations = append(translations, t...)
}

func (svc *Service) GetTranslations() (t []interfaces.Translation) {
	return translations
}

func GetI18nService(cfgPath string) (svc2 interfaces.I18nService, err error) {
	if _svc != nil {
		return _svc, nil
	}

	_svc, err = NewI18nService()
	if err != nil {
		return nil, err
	}

	return _svc, nil
}

func ProvideGetI18nService(cfgPath string) func() (svc interfaces.I18nService, err error) {
	return func() (svc interfaces.I18nService, err error) {
		return GetI18nService(cfgPath)
	}
}

func NewI18nService() (svc2 interfaces.I18nService, err error) {
	svc := &Service{}

	return svc, nil
}
