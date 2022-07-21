package components

type PaymentResponseEntity struct {
	PaidTotal                 float64                             `json:"PaidTotal,omitempty"`
	PaymentReversalDocument   string                              `json:"PaymentReversalDocument,omitempty" validate:"max=250"`
	PaymentTermCode           string                              `json:"PaymentTermCode,omitempty" validate:"max=25"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	SourceSystem              string                              `json:"SourceSystem,omitempty" validate:"max=32"`
	PaymentDate               string                              `json:"PaymentDate,omitempty"`
	CashDiscount              float64                             `json:"CashDiscount,omitempty"`
	ResponseMessage           string                              `json:"ResponseMessage,omitempty" validate:"max=250"`
	PaymentBlock              bool                                `json:"PaymentBlock"`
	PaymentDocument           string                              `json:"PaymentDocument,omitempty" validate:"max=250"`
	Success                   bool                                `json:"Success"`
	CheckNumber               string                              `json:"CheckNumber,omitempty" validate:"max=25"`
	PaymentMessage            string                              `json:"PaymentMessage,omitempty" validate:"max=100"`
	PaymentNumber             string                              `json:"PaymentNumber,omitempty" validate:"max=250"`
	PaymentMethod             string                              `json:"PaymentMethod,omitempty" validate:"max=25"`
}
