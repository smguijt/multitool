package components

type BaswareSystem struct {
	Module                    string                              `json:"Module,omitempty"`
	SystemStatus              string                              `json:"SystemStatus"`
	LastUpdated               string                              `json:"LastUpdated"`
	Items                     []ProcessedItem                     `json:"Items,omitempty"`
	System                    string                              `json:"System,omitempty"`
}
