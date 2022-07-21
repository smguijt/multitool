package components

type TransferParameterEntity struct {
	Key                       string                              `json:"Key" validate:"min=1,max=250"`
	Value                     string                              `json:"Value" validate:"min=1,max=250"`
}
