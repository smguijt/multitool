package components

type VendorProcessingStatusEntity struct {
	PaymentDenied             bool                                `json:"PaymentDenied,omitempty"`
	Active                    bool                                `json:"Active,omitempty"`
}
