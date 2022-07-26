package components

type OrderEntity struct {
	PaymentTermCode           string                              `json:"PaymentTermCode,omitempty" validate:"max=25"`
	OrderNumber               string                              `json:"OrderNumber" validate:"min=1,max=100"`
	OrganizationElementName   string                              `json:"OrganizationElementName,omitempty" validate:"max=250"`
	Text5                     string                              `json:"Text5,omitempty" validate:"max=250"`
	LastUpdated               string                              `json:"LastUpdated"`
	CompanyName               string                              `json:"CompanyName,omitempty" validate:"max=250"`
	SupplierName              string                              `json:"SupplierName,omitempty" validate:"max=250"`
	InvoicingSupplierName     string                              `json:"InvoicingSupplierName,omitempty" validate:"max=250"`
	OrganizationElementCode   string                              `json:"OrganizationElementCode" validate:"min=2,max=25"`
	Text10                    string                              `json:"Text10,omitempty" validate:"max=250"`
	IsClosed                  bool                                `json:"IsClosed,omitempty"`
	Numeric4                  float64                             `json:"Numeric4,omitempty"`
	Text4                     string                              `json:"Text4,omitempty" validate:"max=250"`
	OrderReference            string                              `json:"OrderReference,omitempty" validate:"max=100"`
	CompanyCode               string                              `json:"CompanyCode" validate:"min=1,max=25"`
	IsInvoiced                bool                                `json:"IsInvoiced,omitempty"`
	Description               string                              `json:"Description,omitempty" validate:"max=1000"`
	Created                   string                              `json:"Created,omitempty"`
	PaymentTermName           string                              `json:"PaymentTermName,omitempty" validate:"max=250"`
	Text1                     string                              `json:"Text1,omitempty" validate:"max=250"`
	SourceSystem              string                              `json:"SourceSystem" validate:"min=1,max=36"`
	Text3                     string                              `json:"Text3,omitempty" validate:"max=250"`
	Numeric1                  float64                             `json:"Numeric1,omitempty"`
	Date4                     string                              `json:"Date4,omitempty"`
	ValidTo                   string                              `json:"ValidTo,omitempty"`
	Text9                     string                              `json:"Text9,omitempty" validate:"max=250"`
	InvoicingSupplierCode     string                              `json:"InvoicingSupplierCode,omitempty" validate:"max=25"`
	OrderTypeCode             string                              `json:"OrderTypeCode,omitempty" validate:"max=25"`
	Text7                     string                              `json:"Text7,omitempty" validate:"max=250"`
	Date1                     string                              `json:"Date1,omitempty"`
	PurchaseOrganizationName  string                              `json:"PurchaseOrganizationName,omitempty" validate:"max=250"`
	Date5                     string                              `json:"Date5,omitempty"`
	Date2                     string                              `json:"Date2,omitempty"`
	ValidFrom                 string                              `json:"ValidFrom,omitempty"`
	CurrencyCode              string                              `json:"CurrencyCode" validate:"min=2,max=3"`
	ActualDeliveryDate        string                              `json:"ActualDeliveryDate,omitempty"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	Text2                     string                              `json:"Text2,omitempty" validate:"max=250"`
	PurchaseOrganizationCode  string                              `json:"PurchaseOrganizationCode,omitempty" validate:"max=25"`
	Numeric5                  float64                             `json:"Numeric5,omitempty"`
	SupplierCode              string                              `json:"SupplierCode" validate:"min=1,max=25"`
	Text8                     string                              `json:"Text8,omitempty" validate:"max=250"`
	Numeric2                  float64                             `json:"Numeric2,omitempty"`
	IsDelivered               bool                                `json:"IsDelivered,omitempty"`
	Numeric3                  float64                             `json:"Numeric3,omitempty"`
	Date3                     string                              `json:"Date3,omitempty"`
	RequestedDeliveryDate     string                              `json:"RequestedDeliveryDate,omitempty"`
	Text6                     string                              `json:"Text6,omitempty" validate:"max=250"`
}
