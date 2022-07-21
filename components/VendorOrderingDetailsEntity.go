package components

type VendorOrderingDetailsEntity struct {
	OrderingFormat            string                              `json:"OrderingFormat,omitempty"`
	HasActiveCatalog          bool                                `json:"HasActiveCatalog,omitempty"`
	DeliverOrderAutomatically bool                                `json:"DeliverOrderAutomatically,omitempty"`
	OrderingLanguage          string                              `json:"OrderingLanguage,omitempty" validate:"max=8"`
	OrderEmail                string                              `json:"OrderEmail,omitempty" validate:"min=3,max=320"`
	OrderProcessType          string                              `json:"OrderProcessType,omitempty" validate:"max=50"`
	AutomaticallyReceiveOnOrder bool                                `json:"AutomaticallyReceiveOnOrder,omitempty"`
	MinimumOrderAllowed       float64                             `json:"MinimumOrderAllowed,omitempty"`
	NoFreeformItems           bool                                `json:"NoFreeformItems,omitempty"`
	SendToNetwork             bool                                `json:"SendToNetwork,omitempty"`
	OrderingMessageLanguage   string                              `json:"OrderingMessageLanguage,omitempty" validate:"max=8"`
	CreateOrderAutomatically  bool                                `json:"CreateOrderAutomatically,omitempty"`
	IsTaxable                 bool                                `json:"IsTaxable,omitempty"`
	AutomaticallyReceiveOnInvoice bool                                `json:"AutomaticallyReceiveOnInvoice,omitempty"`
}
