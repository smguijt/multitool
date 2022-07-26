package components

type GoodsReceiptLine struct {
	DeliveryNoteNumber        string                              `json:"DeliveryNoteNumber,omitempty" validate:"max=250"`
	ProductSerialNumber       string                              `json:"ProductSerialNumber,omitempty" validate:"max=250"`
	ReceivedQuantity          float64                             `json:"ReceivedQuantity"`
	ReferenceGRLineExternalCode string                              `json:"ReferenceGRLineExternalCode,omitempty" validate:"max=100"`
	Comment                   string                              `json:"Comment,omitempty" validate:"max=255"`
	NotifyFault               bool                                `json:"NotifyFault"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=100"`
	OrderLineExternalCode     string                              `json:"OrderLineExternalCode" validate:"min=1,max=100"`
}
