package components

type GoodsReceiptEntity struct {
	Text1                     string                              `json:"Text1,omitempty" validate:"max=250"`
	ReceiveMethod             int                                 `json:"ReceiveMethod,omitempty"`
	GrossPrice                float64                             `json:"GrossPrice,omitempty"`
	Date2                     string                              `json:"Date2,omitempty"`
	Date5                     string                              `json:"Date5,omitempty"`
	NetSum                    float64                             `json:"NetSum,omitempty"`
	GoodsReceiptLineNumber    string                              `json:"GoodsReceiptLineNumber,omitempty" validate:"max=25"`
	GoodsReceiptNumber        string                              `json:"GoodsReceiptNumber" validate:"min=1,max=100"`
	IsDeleted                 bool                                `json:"IsDeleted,omitempty"`
	ProductSerialNumber       string                              `json:"ProductSerialNumber,omitempty" validate:"max=250"`
	Text4                     string                              `json:"Text4,omitempty" validate:"max=250"`
	Text9                     string                              `json:"Text9,omitempty" validate:"max=250"`
	InvoicedQuantity          float64                             `json:"InvoicedQuantity,omitempty"`
	Text5                     string                              `json:"Text5,omitempty" validate:"max=250"`
	GrossSum                  float64                             `json:"GrossSum,omitempty"`
	InvoicedNetSum            float64                             `json:"InvoicedNetSum,omitempty"`
	Numeric5                  float64                             `json:"Numeric5,omitempty"`
	Text7                     string                              `json:"Text7,omitempty" validate:"max=250"`
	Text6                     string                              `json:"Text6,omitempty" validate:"max=250"`
	InvoicedGrossSum          float64                             `json:"InvoicedGrossSum,omitempty"`
	BestFitGrouping           string                              `json:"BestFitGrouping,omitempty" validate:"max=250"`
	Quantity                  float64                             `json:"Quantity"`
	Text8                     string                              `json:"Text8,omitempty" validate:"max=250"`
	DeliveryNoteNumber        string                              `json:"DeliveryNoteNumber,omitempty" validate:"max=250"`
	DeliveryDate              string                              `json:"DeliveryDate,omitempty"`
	Numeric2                  float64                             `json:"Numeric2,omitempty"`
	GoodsReceiptType          int                                 `json:"GoodsReceiptType,omitempty"`
	VoucherNumber             string                              `json:"VoucherNumber,omitempty" validate:"max=100"`
	ReferenceGRExternalCode   string                              `json:"ReferenceGRExternalCode,omitempty" validate:"max=36"`
	Text10                    string                              `json:"Text10,omitempty" validate:"max=250"`
	Text3                     string                              `json:"Text3,omitempty" validate:"max=250"`
	UnitOfMeasure             string                              `json:"UnitOfMeasure,omitempty" validate:"max=25"`
	Numeric1                  float64                             `json:"Numeric1,omitempty"`
	FiscalYear                string                              `json:"FiscalYear,omitempty"`
	Date1                     string                              `json:"Date1,omitempty"`
	NetPrice                  float64                             `json:"NetPrice,omitempty"`
	SubUnitOfMeasure          string                              `json:"SubUnitOfMeasure,omitempty" validate:"max=250"`
	Text2                     string                              `json:"Text2,omitempty" validate:"max=250"`
	Date3                     string                              `json:"Date3,omitempty"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	GoodsReceiptNote          string                              `json:"GoodsReceiptNote,omitempty" validate:"max=250"`
	Numeric3                  float64                             `json:"Numeric3,omitempty"`
	Numeric4                  float64                             `json:"Numeric4,omitempty"`
	Date4                     string                              `json:"Date4,omitempty"`
}
