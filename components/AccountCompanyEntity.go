package components

type AccountCompanyEntity struct {
	Inherit                   bool                                `json:"Inherit"`
	Active                    bool                                `json:"Active"`
	CompanyCode               string                              `json:"CompanyCode" validate:"min=2,max=32"`
}
