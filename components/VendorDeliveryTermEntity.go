package components

type VendorDeliveryTermEntity struct {
	DeliveryTermCode          string                              `json:"DeliveryTermCode" validate:"min=1,max=25"`
	DeliveryLocation          string                              `json:"DeliveryLocation,omitempty" validate:"max=100"`
	Description               string                              `json:"Description,omitempty" validate:"max=200"`
}
