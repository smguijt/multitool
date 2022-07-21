package components

type VendorEntity struct {
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	Contacts                  []VendorContactEntity               `json:"Contacts,omitempty"`
	VendorParent              string                              `json:"VendorParent,omitempty"`
	Identifiers               []VendorIdentifierEntity            `json:"Identifiers,omitempty"`
	Name                      string                              `json:"Name" validate:"min=1,max=250"`
	SourceSystem              string                              `json:"SourceSystem,omitempty" validate:"max=36"`
	VendorClass               string                              `json:"VendorClass,omitempty" validate:"max=100"`
	ProcessingStatus          VendorProcessingStatusEntity        `json:"ProcessingStatus"`
	LastUpdated               string                              `json:"LastUpdated"`
	CustomFields              []VendorAdditionalFieldEntity       `json:"CustomFields,omitempty"`
	VendorCode                string                              `json:"VendorCode" validate:"min=1,max=25"`
	Description               string                              `json:"Description,omitempty" validate:"max=250"`
	PaymentMeans              []VendorPaymentMeanEntity           `json:"PaymentMeans,omitempty"`
	EligibleForSourcing       bool                                `json:"EligibleForSourcing,omitempty"`
	Tags                      []VendorTagEntity                   `json:"Tags,omitempty"`
	Companies                 []VendorCompanyEntity               `json:"Companies"`
	PaymentTerms              VendorPaymentTermEntity             `json:"PaymentTerms"`
	SupplierAssignedAccountId string                              `json:"SupplierAssignedAccountId,omitempty" validate:"max=128"`
	Buvid                     string                              `json:"Buvid,omitempty" validate:"max=36"`
	DeliveryTerm              VendorDeliveryTermEntity            `json:"DeliveryTerm"`
	OrderingDetails           VendorOrderingDetailsEntity         `json:"OrderingDetails"`
	Addresses                 []VendorAddressEntity               `json:"Addresses,omitempty"`
}
