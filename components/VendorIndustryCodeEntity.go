package components

type VendorIndustryCodeEntity struct {
	SchemeId                  string                              `json:"SchemeId,omitempty" validate:"max=50"`
	SchemeDescription         string                              `json:"SchemeDescription,omitempty" validate:"max=50"`
	ModificationDate          string                              `json:"ModificationDate,omitempty"`
	Source                    string                              `json:"Source,omitempty" validate:"max=50"`
	Code                      string                              `json:"Code,omitempty" validate:"max=50"`
	Description               string                              `json:"Description,omitempty" validate:"max=256"`
}
