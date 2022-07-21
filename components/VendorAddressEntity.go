package components

type VendorAddressEntity struct {
	CountryId                 string                              `json:"CountryId,omitempty" validate:"min=2,max=2"`
	Locality                  string                              `json:"Locality,omitempty" validate:"max=50"`
	Default                   bool                                `json:"Default,omitempty"`
	Description               string                              `json:"Description,omitempty" validate:"max=200"`
	CountrySubEntityDescription string                              `json:"CountrySubEntityDescription,omitempty" validate:"max=50"`
	Name                      string                              `json:"Name" validate:"min=1,max=200"`
	AddressType               string                              `json:"AddressType,omitempty"`
	AddressLine3              string                              `json:"AddressLine3,omitempty" validate:"max=200"`
	CityName                  string                              `json:"CityName,omitempty" validate:"max=50"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	AdditionalAddressFields   []VendorAddressPartEntity           `json:"AdditionalAddressFields,omitempty"`
	PoBox                     string                              `json:"PoBox,omitempty" validate:"max=50"`
	CountrySubEntity          string                              `json:"CountrySubEntity,omitempty" validate:"max=50"`
	AddressLine2              string                              `json:"AddressLine2,omitempty" validate:"max=200"`
	PostalZone                string                              `json:"PostalZone,omitempty" validate:"max=50"`
	StreetName                string                              `json:"StreetName,omitempty" validate:"max=50"`
	AddressLine1              string                              `json:"AddressLine1,omitempty" validate:"max=200"`
}
