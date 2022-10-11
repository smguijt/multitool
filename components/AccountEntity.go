package components

type AccountEntity struct {
	AccountNames              []LanguageTranslationEntity         `json:"AccountNames"`
	Companies                 []AccountCompanyEntity              `json:"Companies"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	LastUpdated               string                              `json:"LastUpdated"`
	AccountCode               string                              `json:"AccountCode" validate:"min=1,max=25"`
}
