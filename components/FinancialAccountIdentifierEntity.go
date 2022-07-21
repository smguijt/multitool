package components

type FinancialAccountIdentifierEntity struct {
	Description               string                              `json:"Description,omitempty" validate:"max=200"`
	Id                        string                              `json:"Id,omitempty" validate:"max=50"`
	AccountHolderName         string                              `json:"AccountHolderName,omitempty" validate:"max=200"`
	CurrencyCode              string                              `json:"CurrencyCode,omitempty" validate:"min=3,max=3"`
	AccountAdditionalData1    string                              `json:"AccountAdditionalData1,omitempty" validate:"max=200"`
	AccountAdditionalData2    string                              `json:"AccountAdditionalData2,omitempty" validate:"max=200"`
	Default                   bool                                `json:"Default,omitempty"`
	SchemeId                  string                              `json:"SchemeId" validate:"min=1,max=36"`
}
