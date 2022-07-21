package components

type GenericListCompanyEntity struct {
	CompanyCode               string                              `json:"CompanyCode" validate:"min=2,max=32"`
	Active                    bool                                `json:"Active"`
}
