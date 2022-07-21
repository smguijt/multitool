package components

type VendorAddressPartEntity struct {
	Value                     string                              `json:"Value,omitempty" validate:"max=200"`
	Key                       string                              `json:"Key"`
}
