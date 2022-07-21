package components

type VendorCreditRiskDataEntity struct {
	CreditRatingDataFlag      bool                                `json:"CreditRatingDataFlag,omitempty"`
	CreditRating              string                              `json:"CreditRating,omitempty" validate:"max=1000"`
	ModificationDate          string                              `json:"ModificationDate,omitempty"`
	Source                    string                              `json:"Source,omitempty" validate:"max=50"`
}
