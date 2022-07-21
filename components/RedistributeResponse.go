package components

type RedistributeResponse struct {
	TaskName                  string                              `json:"TaskName,omitempty"`
	TaskStatus                string                              `json:"TaskStatus,omitempty"`
	StatusApiLink             string                              `json:"StatusApiLink,omitempty"`
}
