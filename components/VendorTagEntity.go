package components

type VendorTagEntity struct {
	Value                     string                              `json:"Value" validate:"min=1,max=250"`
	TagGroup                  string                              `json:"TagGroup,omitempty" validate:"max=50"`
	Name                      string                              `json:"Name" validate:"min=1,max=50"`
}
