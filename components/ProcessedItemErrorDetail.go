package components

type ProcessedItemErrorDetail struct {
	Type                      string                              `json:"Type"`
	Message                   string                              `json:"Message,omitempty"`
	Category                  string                              `json:"Category"`
}
