package components

type AccountingDocumentEntity struct {
	Text18                    string                              `json:"Text18,omitempty" validate:"max=100"`
	GrossSum                  float64                             `json:"GrossSum,omitempty"`
	CashDate                  string                              `json:"CashDate,omitempty"`
	CodingDate                string                              `json:"CodingDate,omitempty"`
	Numeric17                 float64                             `json:"Numeric17,omitempty"`
	CashSumOrganization       float64                             `json:"CashSumOrganization,omitempty"`
	InvoiceImageURL           string                              `json:"InvoiceImageURL,omitempty" validate:"max=250"`
	PaymentMethod             string                              `json:"PaymentMethod,omitempty" validate:"max=25"`
	InvoiceDate               string                              `json:"InvoiceDate,omitempty"`
	VoucherNumber2            string                              `json:"VoucherNumber2,omitempty" validate:"max=25"`
	CashSum                   float64                             `json:"CashSum,omitempty"`
	VoucherNumber1            string                              `json:"VoucherNumber1,omitempty" validate:"max=25"`
	Text17                    string                              `json:"Text17,omitempty" validate:"max=100"`
	OrganizationCode          string                              `json:"OrganizationCode,omitempty" validate:"max=25"`
	Text30                    string                              `json:"Text30,omitempty" validate:"max=100"`
	LastUpdated               string                              `json:"LastUpdated"`
	RemoveResponses           []RemoveResponseEntity              `json:"RemoveResponses,omitempty"`
	ReferencePerson           string                              `json:"ReferencePerson,omitempty" validate:"max=255"`
	Text23                    string                              `json:"Text23,omitempty" validate:"max=100"`
	DueDate                   string                              `json:"DueDate,omitempty"`
	TransferParameters        []TransferParameterEntity           `json:"TransferParameters,omitempty"`
	Text11                    string                              `json:"Text11,omitempty" validate:"max=100"`
	SupplierCode              string                              `json:"SupplierCode" validate:"max=25"`
	Numeric3                  float64                             `json:"Numeric3,omitempty"`
	Numeric16                 float64                             `json:"Numeric16,omitempty"`
	BaseLineDate              string                              `json:"BaseLineDate,omitempty"`
	OriginService             int                                 `json:"OriginService"`
	Text8                     string                              `json:"Text8,omitempty" validate:"max=250"`
	Text27                    string                              `json:"Text27,omitempty" validate:"max=100"`
	SupplierName              string                              `json:"SupplierName,omitempty" validate:"max=250"`
	Text29                    string                              `json:"Text29,omitempty" validate:"max=100"`
	Text2                     string                              `json:"Text2,omitempty" validate:"max=250"`
	NetSum                    float64                             `json:"NetSum,omitempty"`
	TaxSum1Organization       float64                             `json:"TaxSum1Organization,omitempty"`
	Text25                    string                              `json:"Text25,omitempty" validate:"max=100"`
	TaxSum2Company            float64                             `json:"TaxSum2Company,omitempty"`
	CurrencyCodeCompany       string                              `json:"CurrencyCodeCompany,omitempty" validate:"max=3"`
	SupplierBankIBAN          string                              `json:"SupplierBankIBAN,omitempty" validate:"max=50"`
	InvoiceTypeCode           string                              `json:"InvoiceTypeCode,omitempty" validate:"max=25"`
	Date4                     string                              `json:"Date4,omitempty"`
	PaymentResponses          []PaymentResponseEntity             `json:"PaymentResponses,omitempty"`
	Date5                     string                              `json:"Date5,omitempty"`
	Text4                     string                              `json:"Text4,omitempty" validate:"max=250"`
	Text22                    string                              `json:"Text22,omitempty" validate:"max=100"`
	Text14                    string                              `json:"Text14,omitempty" validate:"max=100"`
	CodingRows                []StandardCodingEntity              `json:"CodingRows,omitempty"`
	Numeric12                 float64                             `json:"Numeric12,omitempty"`
	Text7                     string                              `json:"Text7,omitempty" validate:"max=250"`
	Text15                    string                              `json:"Text15,omitempty" validate:"max=100"`
	PaymentBlock              bool                                `json:"PaymentBlock"`
	GrossSumCompany           float64                             `json:"GrossSumCompany,omitempty"`
	Numeric9                  float64                             `json:"Numeric9,omitempty"`
	OrganizationName          string                              `json:"OrganizationName,omitempty" validate:"max=250"`
	CompanyName               string                              `json:"CompanyName,omitempty" validate:"min=2,max=250"`
	CurrencyCodeOrganization  string                              `json:"CurrencyCodeOrganization,omitempty" validate:"min=2,max=3"`
	PaymentRevelsalDocument   string                              `json:"PaymentRevelsalDocument,omitempty" validate:"max=250"`
	Numeric14                 float64                             `json:"Numeric14,omitempty"`
	Text6                     string                              `json:"Text6,omitempty" validate:"max=250"`
	ExchangeRateOrganization  float64                             `json:"ExchangeRateOrganization,omitempty"`
	Numeric7                  float64                             `json:"Numeric7,omitempty"`
	TaxPercent2               float64                             `json:"TaxPercent2,omitempty"`
	Text20                    string                              `json:"Text20,omitempty" validate:"max=100"`
	Numeric5                  float64                             `json:"Numeric5,omitempty"`
	Numeric4                  float64                             `json:"Numeric4,omitempty"`
	ExchangeRateBaseDate      string                              `json:"ExchangeRateBaseDate,omitempty"`
	Numeric1                  float64                             `json:"Numeric1,omitempty"`
	NetSumCompany             float64                             `json:"NetSumCompany,omitempty"`
	Text9                     string                              `json:"Text9,omitempty" validate:"max=250"`
	Text5                     string                              `json:"Text5,omitempty" validate:"max=250"`
	NetSumOrganization        float64                             `json:"NetSumOrganization,omitempty"`
	Text3                     string                              `json:"Text3,omitempty" validate:"max=250"`
	Date1                     string                              `json:"Date1,omitempty"`
	Numeric18                 float64                             `json:"Numeric18,omitempty"`
	AccountingPeriod          string                              `json:"AccountingPeriod,omitempty" validate:"max=50"`
	Prebooked                 bool                                `json:"Prebooked,omitempty"`
	TaxPercent1               float64                             `json:"TaxPercent1,omitempty"`
	TaxSum2                   float64                             `json:"TaxSum2,omitempty"`
	Date9                     string                              `json:"Date9,omitempty"`
	OrganizationElementCode   string                              `json:"OrganizationElementCode,omitempty" validate:"min=1,max=25"`
	ProcessingStatus          string                              `json:"ProcessingStatus"`
	ParentInvoiceBumId        string                              `json:"ParentInvoiceBumId,omitempty" validate:"max=36"`
	VoucherDate               string                              `json:"VoucherDate,omitempty"`
	CompanyCode               string                              `json:"CompanyCode" validate:"min=2,max=32"`
	Text1                     string                              `json:"Text1,omitempty" validate:"max=250"`
	OrderNumbers              []string                            `json:"OrderNumbers,omitempty"`
	OrganizationElementName   string                              `json:"OrganizationElementName,omitempty" validate:"min=1,max=250"`
	Date2                     string                              `json:"Date2,omitempty"`
	Numeric6                  float64                             `json:"Numeric6,omitempty"`
	CashPercent               float64                             `json:"CashPercent,omitempty"`
	SupplierBankBIC           string                              `json:"SupplierBankBIC,omitempty" validate:"max=250"`
	Text21                    string                              `json:"Text21,omitempty" validate:"max=100"`
	InvoiceNumber             string                              `json:"InvoiceNumber,omitempty" validate:"max=100"`
	PrebookResponses          []PrebookResponseEntity             `json:"PrebookResponses,omitempty"`
	Text26                    string                              `json:"Text26,omitempty" validate:"max=100"`
	Numeric2                  float64                             `json:"Numeric2,omitempty"`
	PaymentPlanReference      string                              `json:"PaymentPlanReference,omitempty" validate:"max=255"`
	Date3                     string                              `json:"Date3,omitempty"`
	Description               string                              `json:"Description,omitempty" validate:"max=250"`
	Date6                     string                              `json:"Date6,omitempty"`
	Date10                    string                              `json:"Date10,omitempty"`
	Text12                    string                              `json:"Text12,omitempty" validate:"max=100"`
	TransferResponses         []TransferResponseEntity            `json:"TransferResponses,omitempty"`
	SupplierSourceSystemId    string                              `json:"SupplierSourceSystemId,omitempty" validate:"max=36"`
	PaymentTermExternalCode   string                              `json:"PaymentTermExternalCode,omitempty" validate:"max=100"`
	PaymentTermCode           string                              `json:"PaymentTermCode,omitempty" validate:"max=25"`
	Numeric11                 float64                             `json:"Numeric11,omitempty"`
	Text19                    string                              `json:"Text19,omitempty" validate:"max=100"`
	Numeric8                  float64                             `json:"Numeric8,omitempty"`
	Numeric13                 float64                             `json:"Numeric13,omitempty"`
	CurrencyCode              string                              `json:"CurrencyCode" validate:"min=2,max=3"`
	Text13                    string                              `json:"Text13,omitempty" validate:"max=100"`
	ReferenceNumber           string                              `json:"ReferenceNumber,omitempty" validate:"max=50"`
	Bumid                     string                              `json:"Bumid,omitempty" validate:"max=36"`
	Numeric20                 float64                             `json:"Numeric20,omitempty"`
	PaymentTermName           string                              `json:"PaymentTermName,omitempty" validate:"max=250"`
	Date7                     string                              `json:"Date7,omitempty"`
	ExchangeRateCompany       float64                             `json:"ExchangeRateCompany,omitempty"`
	Date8                     string                              `json:"Date8,omitempty"`
	InvoiceTypeName           string                              `json:"InvoiceTypeName,omitempty" validate:"max=250"`
	SupplierBankName          string                              `json:"SupplierBankName,omitempty" validate:"max=250"`
	Numeric15                 float64                             `json:"Numeric15,omitempty"`
	TaxCode                   string                              `json:"TaxCode,omitempty" validate:"max=25"`
	AccountingGroup           string                              `json:"AccountingGroup,omitempty" validate:"max=50"`
	Text24                    string                              `json:"Text24,omitempty" validate:"max=100"`
	TaxSum1                   float64                             `json:"TaxSum1,omitempty"`
	Text28                    string                              `json:"Text28,omitempty" validate:"max=100"`
	Text10                    string                              `json:"Text10,omitempty" validate:"max=250"`
	TaxSum2Organization       float64                             `json:"TaxSum2Organization,omitempty"`
	CashSumCompany            float64                             `json:"CashSumCompany,omitempty"`
	Numeric19                 float64                             `json:"Numeric19,omitempty"`
	ContractNumber            string                              `json:"ContractNumber,omitempty" validate:"max=250"`
	TaxSum1Company            float64                             `json:"TaxSum1Company,omitempty"`
	InvoiceId                 string                              `json:"InvoiceId" validate:"min=1,max=36"`
	SupplierBankBBAN          string                              `json:"SupplierBankBBAN,omitempty" validate:"max=50"`
	InvoiceImageToken         string                              `json:"InvoiceImageToken,omitempty" validate:"max=250"`
	Numeric10                 float64                             `json:"Numeric10,omitempty"`
	Text16                    string                              `json:"Text16,omitempty" validate:"max=100"`
	GrossSumOrganization      float64                             `json:"GrossSumOrganization,omitempty"`
}
