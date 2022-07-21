package components

type VendorIdentifierEntity struct {
	Id                        string                              `json:"Id" validate:"min=1,max=36"`
	SchemeId                  string                              `json:"SchemeId" validate:"min=1,max=32"`
	DefaultPartyId            bool                                `json:"DefaultPartyId,omitempty"`
}
