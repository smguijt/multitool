package components

type TransferResponseEntity struct {
	VoucherNumber1            string                              `json:"VoucherNumber1,omitempty" validate:"max=25"`
	VoucherNumber2            string                              `json:"VoucherNumber2,omitempty" validate:"max=25"`
	TransferDate              string                              `json:"TransferDate,omitempty"`
	PaymentBlock              bool                                `json:"PaymentBlock"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	Success                   bool                                `json:"Success"`
	ResponseMessage           string                              `json:"ResponseMessage" validate:"min=1,max=250"`
	SourceSystem              string                              `json:"SourceSystem,omitempty" validate:"max=32"`
}
