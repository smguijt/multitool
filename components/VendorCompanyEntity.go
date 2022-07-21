package components

type VendorCompanyEntity struct {
	CompanyCode               string                              `json:"CompanyCode" validate:"min=2,max=32"`
	InheritToChildUnits       bool                                `json:"InheritToChildUnits,omitempty"`
}
