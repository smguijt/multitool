package components

type FinancialInstitutionEntity struct {
	SchemeId                  string                              `json:"SchemeId" validate:"min=1,max=36"`
	BranchId                  string                              `json:"BranchId,omitempty" validate:"max=50"`
	BranchIdSchemeId          string                              `json:"BranchIdSchemeId,omitempty" validate:"max=36"`
	CountryId                 string                              `json:"CountryId,omitempty" validate:"min=2,max=2"`
	Name                      string                              `json:"Name" validate:"max=250"`
	Id                        string                              `json:"Id" validate:"min=1,max=20"`
}
