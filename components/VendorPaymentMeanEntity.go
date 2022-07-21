package components

type VendorPaymentMeanEntity struct {
	FinancialAccounts         []VendorfinancialAccountsEntity     `json:"FinancialAccounts,omitempty"`
	Default                   bool                                `json:"Default,omitempty"`
	PaymentMeansCode          string                              `json:"PaymentMeansCode" validate:"min=1,max=25"`
	Description               string                              `json:"Description,omitempty" validate:"max=200"`
	CurrencyCode              string                              `json:"CurrencyCode,omitempty" validate:"min=3,max=3"`
}
