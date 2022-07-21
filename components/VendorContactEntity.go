package components

type VendorContactEntity struct {
	Telephone                 string                              `json:"Telephone,omitempty" validate:"max=50"`
	Telefax                   string                              `json:"Telefax,omitempty" validate:"max=50"`
	Email                     string                              `json:"Email,omitempty" validate:"max=200"`
	Website                   string                              `json:"Website,omitempty" validate:"max=200"`
	Role                      string                              `json:"Role"`
	Name                      string                              `json:"Name" validate:"min=1,max=200"`
	Description               string                              `json:"Description,omitempty" validate:"max=200"`
}
