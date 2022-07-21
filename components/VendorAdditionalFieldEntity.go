package components

type VendorAdditionalFieldEntity struct {
	GroupName                 string                              `json:"GroupName,omitempty" validate:"max=500"`
	Name                      string                              `json:"Name" validate:"min=1,max=25"`
	Value                     string                              `json:"Value" validate:"min=1,max=250"`
}
