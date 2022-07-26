package components

type OrderLineUserEntity struct {
	UserExternalCode          string                              `json:"UserExternalCode,omitempty" validate:"max=36"`
	UserEmail                 string                              `json:"UserEmail,omitempty" validate:"max=250"`
	UserRole                  string                              `json:"UserRole,omitempty"`
	LastUpdated               string                              `json:"LastUpdated"`
}
