package components

type ConsumerLogEventUserView struct {
	LastUpdated               string                              `json:"LastUpdated"`
	Systems                   []BaswareSystem                     `json:"Systems,omitempty"`
	RequestId                 string                              `json:"RequestId,omitempty"`
	EntityType                string                              `json:"EntityType"`
	EntityVersion             string                              `json:"EntityVersion,omitempty"`
	RequestStatus             string                              `json:"RequestStatus"`
}
