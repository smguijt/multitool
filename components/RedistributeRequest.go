package components

type RedistributeRequest struct {
	SubscribedService         string                              `json:"SubscribedService"`
	LastUpdated               string                              `json:"LastUpdated,omitempty"`
	ExternalCodes             []string                            `json:"ExternalCodes,omitempty"`
	ListKey                   string                              `json:"ListKey,omitempty"`
	Version                   string                              `json:"Version,omitempty"`
	EntityType                string                              `json:"EntityType"`
}
