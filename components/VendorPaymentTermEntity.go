package components

type VendorPaymentTermEntity struct {
	Description               string                              `json:"Description,omitempty" validate:"max=200"`
	PaymentTermCode           string                              `json:"PaymentTermCode" validate:"min=1,max=25"`
}
