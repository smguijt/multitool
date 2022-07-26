package components

type OrderLineCodingEntity struct {
	Num3                      float64                             `json:"Num3,omitempty"`
	AllocatedQuantity         float64                             `json:"AllocatedQuantity"`
	ConversionDeNumerator     float64                             `json:"ConversionDeNumerator,omitempty"`
	Text2                     string                              `json:"Text2,omitempty" validate:"max=250"`
	DimCode3                  string                              `json:"DimCode3,omitempty" validate:"max=100"`
	CustomerCode              string                              `json:"CustomerCode,omitempty" validate:"max=50"`
	ProfitCenterName          string                              `json:"ProfitCenterName,omitempty" validate:"max=50"`
	DimCode7                  string                              `json:"DimCode7,omitempty" validate:"max=100"`
	Text1                     string                              `json:"Text1,omitempty" validate:"max=250"`
	BudgetName                string                              `json:"BudgetName,omitempty" validate:"max=50"`
	DimCode2                  string                              `json:"DimCode2,omitempty" validate:"max=100"`
	Date2                     string                              `json:"Date2,omitempty"`
	BusinessAreaName          string                              `json:"BusinessAreaName,omitempty" validate:"max=50"`
	DimCode1                  string                              `json:"DimCode1,omitempty" validate:"max=100"`
	Date4                     string                              `json:"Date4,omitempty"`
	SalesOrderSubName         string                              `json:"SalesOrderSubName,omitempty" validate:"max=50"`
	ProjectSubCode            string                              `json:"ProjectSubCode,omitempty" validate:"max=50"`
	AccountCode               string                              `json:"AccountCode,omitempty" validate:"max=50"`
	WorkOrderCode             string                              `json:"WorkOrderCode,omitempty" validate:"max=50"`
	CostCenterName            string                              `json:"CostCenterName,omitempty" validate:"max=250"`
	CustomerName              string                              `json:"CustomerName,omitempty" validate:"max=50"`
	Date1                     string                              `json:"Date1,omitempty"`
	DimName6                  string                              `json:"DimName6,omitempty" validate:"max=100"`
	ProjectSubName            string                              `json:"ProjectSubName,omitempty" validate:"max=50"`
	InternalOrderName         string                              `json:"InternalOrderName,omitempty" validate:"max=50"`
	CostCenterCode            string                              `json:"CostCenterCode,omitempty" validate:"max=200"`
	ConversionDenominator     float64                             `json:"ConversionDenominator,omitempty"`
	Num5                      float64                             `json:"Num5,omitempty"`
	GrossTotal                float64                             `json:"GrossTotal,omitempty"`
	Network                   string                              `json:"Network,omitempty" validate:"max=250"`
	ServiceCode               string                              `json:"ServiceCode,omitempty" validate:"max=50"`
	WorkOrderName             string                              `json:"WorkOrderName,omitempty" validate:"max=50"`
	AccAssignmentCategoryCode string                              `json:"AccAssignmentCategoryCode,omitempty" validate:"max=50"`
	ControllingArea           string                              `json:"ControllingArea,omitempty" validate:"max=50"`
	TaxSum                    float64                             `json:"TaxSum,omitempty"`
	SalesOrderCode            string                              `json:"SalesOrderCode,omitempty" validate:"max=50"`
	PartnerProfitCenter       string                              `json:"PartnerProfitCenter,omitempty" validate:"max=250"`
	BusinessUnitCode          string                              `json:"BusinessUnitCode,omitempty" validate:"max=50"`
	DimCode4                  string                              `json:"DimCode4,omitempty" validate:"max=100"`
	Num1                      float64                             `json:"Num1,omitempty"`
	TaxPercent                float64                             `json:"TaxPercent,omitempty"`
	FixedAssetSubCode         string                              `json:"FixedAssetSubCode,omitempty" validate:"max=50"`
	BusinessAreaCode          string                              `json:"BusinessAreaCode,omitempty" validate:"max=50"`
	InternalOrderCode         string                              `json:"InternalOrderCode,omitempty" validate:"max=50"`
	DimName7                  string                              `json:"DimName7,omitempty" validate:"max=100"`
	CommitmentItem            string                              `json:"CommitmentItem,omitempty" validate:"max=50"`
	DimCode8                  string                              `json:"DimCode8,omitempty" validate:"max=100"`
	RowIndex                  int                                 `json:"RowIndex"`
	FixedAssetSubName         string                              `json:"FixedAssetSubName,omitempty" validate:"max=50"`
	VehicleNumber             string                              `json:"VehicleNumber,omitempty" validate:"max=50"`
	DimName3                  string                              `json:"DimName3,omitempty" validate:"max=100"`
	FunctionalArea            string                              `json:"FunctionalArea,omitempty" validate:"max=50"`
	EmployeeName              string                              `json:"EmployeeName,omitempty" validate:"max=50"`
	SalesOrderSubCode         string                              `json:"SalesOrderSubCode,omitempty" validate:"max=50"`
	DimCode5                  string                              `json:"DimCode5,omitempty" validate:"max=100"`
	DimCode9                  string                              `json:"DimCode9,omitempty" validate:"max=100"`
	DimName1                  string                              `json:"DimName1,omitempty" validate:"max=100"`
	AccAssignmentCategoryName string                              `json:"AccAssignmentCategoryName,omitempty" validate:"max=50"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	Text3                     string                              `json:"Text3,omitempty" validate:"max=250"`
	ProfitCenterCode          string                              `json:"ProfitCenterCode,omitempty" validate:"max=50"`
	FixedAssetName            string                              `json:"FixedAssetName,omitempty" validate:"max=250"`
	FixedAssetCode            string                              `json:"FixedAssetCode,omitempty" validate:"max=250"`
	TaxSum2                   float64                             `json:"TaxSum2,omitempty"`
	DimName4                  string                              `json:"DimName4,omitempty" validate:"max=100"`
	Date5                     string                              `json:"Date5,omitempty"`
	Num4                      float64                             `json:"Num4,omitempty"`
	MaterialGroup             string                              `json:"MaterialGroup,omitempty" validate:"max=100"`
	Text4                     string                              `json:"Text4,omitempty" validate:"max=250"`
	BusinessUnitName          string                              `json:"BusinessUnitName,omitempty" validate:"max=50"`
	Date3                     string                              `json:"Date3,omitempty"`
	SubUOM                    string                              `json:"SubUOM,omitempty" validate:"max=50"`
	Text5                     string                              `json:"Text5,omitempty" validate:"max=250"`
	ConversionNumerator       float64                             `json:"ConversionNumerator,omitempty"`
	DimName9                  string                              `json:"DimName9,omitempty" validate:"max=100"`
	ProjectName               string                              `json:"ProjectName,omitempty" validate:"max=250"`
	DimCode6                  string                              `json:"DimCode6,omitempty" validate:"max=100"`
	AccountName               string                              `json:"AccountName,omitempty" validate:"max=250"`
	NetworkActivity           string                              `json:"NetworkActivity,omitempty" validate:"max=250"`
	SalesOrderName            string                              `json:"SalesOrderName,omitempty" validate:"max=50"`
	DimName10                 string                              `json:"DimName10,omitempty" validate:"max=100"`
	TaxJurisdictionCode       string                              `json:"TaxJurisdictionCode,omitempty" validate:"max=50"`
	EmployeeCode              string                              `json:"EmployeeCode,omitempty" validate:"max=50"`
	VehicleName               string                              `json:"VehicleName,omitempty" validate:"max=50"`
	TaxCode                   string                              `json:"TaxCode,omitempty" validate:"max=25"`
	BudgetCode                string                              `json:"BudgetCode,omitempty" validate:"max=50"`
	DimName2                  string                              `json:"DimName2,omitempty" validate:"max=100"`
	DimCode10                 string                              `json:"DimCode10,omitempty" validate:"max=100"`
	ProjectCode               string                              `json:"ProjectCode,omitempty" validate:"max=25"`
	ServiceName               string                              `json:"ServiceName,omitempty" validate:"max=50"`
	TaxPercent2               float64                             `json:"TaxPercent2,omitempty"`
	DimName5                  string                              `json:"DimName5,omitempty" validate:"max=100"`
	WorkOrderSubCode          string                              `json:"WorkOrderSubCode,omitempty" validate:"max=50"`
	Num2                      float64                             `json:"Num2,omitempty"`
	NetTotal                  float64                             `json:"NetTotal,omitempty"`
	DimName8                  string                              `json:"DimName8,omitempty" validate:"max=100"`
	WorkOrderSubName          string                              `json:"WorkOrderSubName,omitempty" validate:"max=50"`
}
