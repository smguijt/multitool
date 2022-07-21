package components

type VendorEnrichedFieldsEntity struct {
	IndustryCodes             []VendorIndustryCodeEntity          `json:"IndustryCodes,omitempty"`
	BusinessDescription       []VendorBusinessDescriptionEntity   `json:"BusinessDescription,omitempty"`
	CreditRiskData            []VendorCreditRiskDataEntity        `json:"CreditRiskData,omitempty"`
	ComplianceData            []VendorComplianceDataEntity        `json:"ComplianceData,omitempty"`
}
