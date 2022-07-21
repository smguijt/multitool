package components

type RemoveResponseEntity struct {
	ResponseMessage           string                              `json:"ResponseMessage" validate:"max=250"`
	SourceSystem              string                              `json:"SourceSystem,omitempty" validate:"max=32"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=2,max=36"`
	Success                   bool                                `json:"Success,omitempty"`
}
