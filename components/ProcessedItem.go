package components

type ProcessedItem struct {
	Status                    string                              `json:"Status"`
	Errors                    []ProcessedItemErrorDetail          `json:"Errors,omitempty"`
	ExternalCode              string                              `json:"ExternalCode,omitempty"`
}
