package components

type VendorComplianceDataEntity struct {
	ComplianceType            VendorComplianceTypeEntity          `json:"ComplianceType"`
	ModificationDate          string                              `json:"ModificationDate,omitempty"`
	Source                    string                              `json:"Source,omitempty" validate:"max=50"`
}
