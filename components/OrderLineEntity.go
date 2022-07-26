package components

type OrderLineEntity struct {
	LineNumber                string                              `json:"LineNumber" validate:"min=1,max=100"`
	Text10                    string                              `json:"Text10,omitempty" validate:"max=250"`
	Numeric4                  float64                             `json:"Numeric4,omitempty"`
	IsSelfApproved            bool                                `json:"IsSelfApproved,omitempty"`
	Text8                     string                              `json:"Text8,omitempty" validate:"max=250"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	OrderLineCoding           []OrderLineCodingEntity             `json:"OrderLineCoding,omitempty"`
	Date5                     string                              `json:"Date5,omitempty"`
	Text3                     string                              `json:"Text3,omitempty" validate:"max=250"`
	Unspsc                    string                              `json:"Unspsc,omitempty" validate:"max=25"`
	RequestedDeliveryDate     string                              `json:"RequestedDeliveryDate,omitempty"`
	Description               string                              `json:"Description,omitempty" validate:"max=1000"`
	ProductCode               string                              `json:"ProductCode,omitempty" validate:"max=250"`
	Text6                     string                              `json:"Text6,omitempty" validate:"max=250"`
	Text1                     string                              `json:"Text1,omitempty" validate:"max=250"`
	IsReceiptRequired         bool                                `json:"IsReceiptRequired,omitempty"`
	TaxSum                    float64                             `json:"TaxSum,omitempty"`
	Numeric3                  float64                             `json:"Numeric3,omitempty"`
	SortNumber                int                                 `json:"SortNumber,omitempty"`
	Date4                     string                              `json:"Date4,omitempty"`
	InvoicedNetSum            float64                             `json:"InvoicedNetSum,omitempty"`
	Quantity                  float64                             `json:"Quantity"`
	Numeric2                  float64                             `json:"Numeric2,omitempty"`
	GoodsReceipts             []GoodsReceiptEntity                `json:"GoodsReceipts,omitempty"`
	GrossPrice                float64                             `json:"GrossPrice,omitempty"`
	Text2                     string                              `json:"Text2,omitempty" validate:"max=250"`
	ContractNumber            string                              `json:"ContractNumber,omitempty" validate:"max=255"`
	ValidFrom                 string                              `json:"ValidFrom,omitempty"`
	ProductName               string                              `json:"ProductName,omitempty" validate:"max=250"`
	Text5                     string                              `json:"Text5,omitempty" validate:"max=250"`
	MatchingMode              string                              `json:"MatchingMode"`
	MaterialGroup             string                              `json:"MaterialGroup,omitempty" validate:"max=250"`
	ActualDeliveryDate        string                              `json:"ActualDeliveryDate,omitempty"`
	OrderExternalCode         string                              `json:"OrderExternalCode" validate:"min=1,max=36"`
	IsReceiptBasedMatching    bool                                `json:"IsReceiptBasedMatching,omitempty"`
	ValidTo                   string                              `json:"ValidTo,omitempty"`
	NetPrice                  float64                             `json:"NetPrice,omitempty"`
	ReferenceUsers            []OrderLineUserEntity               `json:"ReferenceUsers,omitempty"`
	InvoicedQuantity          float64                             `json:"InvoicedQuantity,omitempty"`
	SubUOM                    string                              `json:"SubUOM,omitempty" validate:"max=250"`
	Date3                     string                              `json:"Date3,omitempty"`
	Comment                   string                              `json:"Comment,omitempty" validate:"max=1000"`
	CurrencyCode              string                              `json:"CurrencyCode,omitempty" validate:"min=2,max=3"`
	TaxSum2                   float64                             `json:"TaxSum2,omitempty"`
	Text9                     string                              `json:"Text9,omitempty" validate:"max=250"`
	IsOverreceivalAllowed     bool                                `json:"IsOverreceivalAllowed,omitempty"`
	BuyerProductCode          string                              `json:"BuyerProductCode,omitempty" validate:"max=250"`
	GlobalTradeItemNumber     string                              `json:"GlobalTradeItemNumber,omitempty" validate:"max=250"`
	Numeric1                  float64                             `json:"Numeric1,omitempty"`
	Date2                     string                              `json:"Date2,omitempty"`
	InvoicedGrossSum          float64                             `json:"InvoicedGrossSum,omitempty"`
	PriceUnitDescription      string                              `json:"PriceUnitDescription,omitempty" validate:"max=250"`
	Text7                     string                              `json:"Text7,omitempty" validate:"max=250"`
	TaxCode                   string                              `json:"TaxCode,omitempty" validate:"max=25"`
	IsDeleted                 bool                                `json:"IsDeleted,omitempty"`
	GrossSum                  float64                             `json:"GrossSum,omitempty"`
	LastUpdated               string                              `json:"LastUpdated"`
	TaxPercent2               float64                             `json:"TaxPercent2,omitempty"`
	PriceUnit                 string                              `json:"PriceUnit,omitempty" validate:"max=25"`
	TaxPercent                float64                             `json:"TaxPercent,omitempty"`
	Text4                     string                              `json:"Text4,omitempty" validate:"max=250"`
	NetSum                    float64                             `json:"NetSum,omitempty"`
	IsClosed                  bool                                `json:"IsClosed,omitempty"`
	Numeric5                  float64                             `json:"Numeric5,omitempty"`
	Uom                       string                              `json:"Uom,omitempty" validate:"max=25"`
	Date1                     string                              `json:"Date1,omitempty"`
}
