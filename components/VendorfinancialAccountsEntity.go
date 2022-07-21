package components

type VendorfinancialAccountsEntity struct {
	FinancialInstitution      FinancialInstitutionEntity          `json:"FinancialInstitution"`
	FinancialAccountIdentifiers []FinancialAccountIdentifierEntity  `json:"FinancialAccountIdentifiers,omitempty"`
}
