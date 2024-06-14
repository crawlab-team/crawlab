package interfaces

type I18nService interface {
	AddTranslations(t []Translation)
	GetTranslations() (t []Translation)
}
