package components

type VendorCandidateModelEntity struct {
	Tags                      []VendorTagEntity                   `json:"Tags,omitempty"`
	Buvid                     string                              `json:"Buvid,omitempty" validate:"max=36"`
	Name                      string                              `json:"Name" validate:"min=1,max=250"`
	ProcessingStatus          VendorProcessingStatusEntity        `json:"ProcessingStatus"`
	CustomFields              []VendorAdditionalFieldEntity       `json:"CustomFields,omitempty"`
	OrderingDetails           VendorOrderingDetailsEntity         `json:"OrderingDetails"`
	VendorCode                string                              `json:"VendorCode,omitempty" validate:"max=25"`
	SupplierAssignedAccountId string                              `json:"SupplierAssignedAccountId,omitempty" validate:"max=128"`
	PaymentTerms              VendorPaymentTermEntity             `json:"PaymentTerms"`
	VendorParent              string                              `json:"VendorParent,omitempty"`
	LastUpdated               string                              `json:"LastUpdated"`
	DeliveryTerm              VendorDeliveryTermEntity            `json:"DeliveryTerm"`
	PaymentMeans              []VendorPaymentMeanEntity           `json:"PaymentMeans,omitempty"`
	Companies                 []VendorCompanyEntity               `json:"Companies,omitempty"`
	Addresses                 []VendorAddressEntity               `json:"Addresses,omitempty"`
	ExternalCode              string                              `json:"ExternalCode,omitempty" validate:"max=36"`
	VendorClass               string                              `json:"VendorClass,omitempty" validate:"max=100"`
	EligibleForSourcing       bool                                `json:"EligibleForSourcing,omitempty"`
	SourceSystem              string                              `json:"SourceSystem,omitempty" validate:"max=36"`
	Identifiers               []VendorIdentifierEntity            `json:"Identifiers,omitempty"`
	Contacts                  []VendorContactEntity               `json:"Contacts,omitempty"`
	Description               string                              `json:"Description,omitempty" validate:"max=250"`
}
