package components

type VendorComplianceTypeEntity struct {
	NoOfDisqualifiedDirector  int                                 `json:"NoOfDisqualifiedDirector,omitempty"`
	NoOfIdv                   int                                 `json:"NoOfIdv,omitempty"`
	NoOfSanctions             int                                 `json:"NoOfSanctions,omitempty"`
	NoOfPep                   int                                 `json:"NoOfPep,omitempty"`
	NoOfLawenforcement        int                                 `json:"NoOfLawenforcement,omitempty"`
	NoOfFinancialRegulator    int                                 `json:"NoOfFinancialRegulator,omitempty"`
	NoOfAdverseMedia          int                                 `json:"NoOfAdverseMedia,omitempty"`
	NoOfInsolvency            int                                 `json:"NoOfInsolvency,omitempty"`
}
