package components

type VendorBusinessDescriptionEntity struct {
	Source                    string                              `json:"Source,omitempty" validate:"max=50"`
	Type                      string                              `json:"Type,omitempty" validate:"max=50"`
	ModificationDate          string                              `json:"ModificationDate,omitempty"`
}
