package components

type StandardCodingEntity struct {
	DimCode10                 string                              `json:"DimCode10,omitempty" validate:"max=250"`
	TaxCode                   string                              `json:"TaxCode,omitempty" validate:"max=250"`
	SalesOrderName            string                              `json:"SalesOrderName,omitempty" validate:"max=250"`
	DimName5                  string                              `json:"DimName5,omitempty" validate:"max=250"`
	Date2                     string                              `json:"Date2,omitempty"`
	InternalOrderCode         string                              `json:"InternalOrderCode,omitempty" validate:"max=250"`
	Num5                      float64                             `json:"Num5,omitempty"`
	OrganizationElementName   string                              `json:"OrganizationElementName,omitempty" validate:"max=250"`
	FixedAssetCode            string                              `json:"FixedAssetCode,omitempty" validate:"max=250"`
	InternalOrderName         string                              `json:"InternalOrderName,omitempty" validate:"max=250"`
	CommitmentItem            string                              `json:"CommitmentItem,omitempty" validate:"max=250"`
	TaxSumCompany             float64                             `json:"TaxSumCompany,omitempty"`
	OrderLineGrossTotal       float64                             `json:"OrderLineGrossTotal,omitempty"`
	EmployeeName              string                              `json:"EmployeeName,omitempty" validate:"max=250"`
	VehicleNumber             string                              `json:"VehicleNumber,omitempty" validate:"max=250"`
	DimCode1                  string                              `json:"DimCode1,omitempty" validate:"max=250"`
	DimName9                  string                              `json:"DimName9,omitempty" validate:"max=250"`
	ContractNumber            string                              `json:"ContractNumber,omitempty" validate:"max=255"`
	BudgetName                string                              `json:"BudgetName,omitempty" validate:"max=250"`
	Date5                     string                              `json:"Date5,omitempty"`
	WorkOrderName             string                              `json:"WorkOrderName,omitempty" validate:"max=250"`
	DimCode8                  string                              `json:"DimCode8,omitempty" validate:"max=250"`
	BusinessUnitName          string                              `json:"BusinessUnitName,omitempty" validate:"max=250"`
	GoodsReceiptItemNumber    string                              `json:"GoodsReceiptItemNumber,omitempty" validate:"max=250"`
	MatchingType              int                                 `json:"MatchingType,omitempty"`
	NetTotalCompany           float64                             `json:"NetTotalCompany,omitempty"`
	AccountCode               string                              `json:"AccountCode,omitempty" validate:"max=50"`
	TaxSum2Organization       float64                             `json:"TaxSum2Organization,omitempty"`
	CostCenterName            string                              `json:"CostCenterName,omitempty" validate:"max=250"`
	DimName1                  string                              `json:"DimName1,omitempty" validate:"max=250"`
	FixedAssetSubName         string                              `json:"FixedAssetSubName,omitempty" validate:"max=250"`
	ProductCode               string                              `json:"ProductCode,omitempty" validate:"max=250"`
	NetTotal                  float64                             `json:"NetTotal,omitempty"`
	Date4                     string                              `json:"Date4,omitempty"`
	DimName6                  string                              `json:"DimName6,omitempty" validate:"max=250"`
	OrderItemNumber           string                              `json:"OrderItemNumber,omitempty"`
	Text1                     string                              `json:"Text1,omitempty" validate:"max=250"`
	ConversionNumerator       float64                             `json:"ConversionNumerator,omitempty"`
	OrganizationElementCode   string                              `json:"OrganizationElementCode,omitempty" validate:"max=250"`
	OwnerName                 string                              `json:"OwnerName,omitempty" validate:"max=250"`
	TaxJurisdictionCode       string                              `json:"TaxJurisdictionCode,omitempty" validate:"max=250"`
	GrossTotalOrganization    float64                             `json:"GrossTotalOrganization,omitempty"`
	Num2                      float64                             `json:"Num2,omitempty"`
	TaxSum2Company            float64                             `json:"TaxSum2Company,omitempty"`
	TaxSumOrganization        float64                             `json:"TaxSumOrganization,omitempty"`
	MaterialGroup             string                              `json:"MaterialGroup,omitempty" validate:"max=250"`
	TaxPercent2               float64                             `json:"TaxPercent2,omitempty"`
	ProfitCenterCode          string                              `json:"ProfitCenterCode,omitempty" validate:"max=250"`
	ControllingArea           string                              `json:"ControllingArea,omitempty" validate:"max=250"`
	MatchedQuantity           float64                             `json:"MatchedQuantity,omitempty"`
	EmployeeCode              string                              `json:"EmployeeCode,omitempty" validate:"max=250"`
	BuyerName                 string                              `json:"BuyerName,omitempty" validate:"max=250"`
	SalesOrderCode            string                              `json:"SalesOrderCode,omitempty" validate:"max=250"`
	OrderedQuantity           float64                             `json:"OrderedQuantity,omitempty"`
	ServiceCode               string                              `json:"ServiceCode,omitempty" validate:"max=250"`
	BusinessUnitCode          string                              `json:"BusinessUnitCode,omitempty" validate:"max=250"`
	OrderLineNetTotal         float64                             `json:"OrderLineNetTotal,omitempty"`
	FiscalYear                string                              `json:"FiscalYear,omitempty" validate:"max=250"`
	WorkOrderSubName          string                              `json:"WorkOrderSubName,omitempty" validate:"max=250"`
	ProjectName               string                              `json:"ProjectName,omitempty" validate:"max=250"`
	CloseOrder                string                              `json:"CloseOrder,omitempty" validate:"max=100"`
	ProfitCenterName          string                              `json:"ProfitCenterName,omitempty" validate:"max=250"`
	Date3                     string                              `json:"Date3,omitempty"`
	Text3                     string                              `json:"Text3,omitempty" validate:"max=250"`
	Num3                      float64                             `json:"Num3,omitempty"`
	Date1                     string                              `json:"Date1,omitempty"`
	DimCode5                  string                              `json:"DimCode5,omitempty" validate:"max=250"`
	OrderLinePriceUnit        string                              `json:"OrderLinePriceUnit,omitempty" validate:"max=250"`
	TaxSum2                   float64                             `json:"TaxSum2,omitempty"`
	GrossTotal                float64                             `json:"GrossTotal,omitempty"`
	SalesOrderSubName         string                              `json:"SalesOrderSubName,omitempty" validate:"max=250"`
	TaxSum                    float64                             `json:"TaxSum,omitempty"`
	BudgetCode                string                              `json:"BudgetCode,omitempty" validate:"max=250"`
	ProductName               string                              `json:"ProductName,omitempty" validate:"max=250"`
	Plant                     string                              `json:"Plant,omitempty" validate:"max=250"`
	ReceivedQuantity          float64                             `json:"ReceivedQuantity,omitempty"`
	OrderGrossTotal           float64                             `json:"OrderGrossTotal,omitempty"`
	Num1                      float64                             `json:"Num1,omitempty"`
	ReceivedNetPrice          float64                             `json:"ReceivedNetPrice,omitempty"`
	TaxPercent                float64                             `json:"TaxPercent,omitempty"`
	Text2                     string                              `json:"Text2,omitempty" validate:"max=250"`
	DimCode3                  string                              `json:"DimCode3,omitempty" validate:"max=250"`
	ConversionDenominator     float64                             `json:"ConversionDenominator,omitempty"`
	DimCode2                  string                              `json:"DimCode2,omitempty" validate:"max=250"`
	MatchedGrossSum           float64                             `json:"MatchedGrossSum,omitempty"`
	ProjectCode               string                              `json:"ProjectCode,omitempty" validate:"max=250"`
	FixedAssetSubCode         string                              `json:"FixedAssetSubCode,omitempty" validate:"max=250"`
	Text4                     string                              `json:"Text4,omitempty" validate:"max=250"`
	Network                   string                              `json:"Network,omitempty" validate:"max=250"`
	MatchedNetSum             float64                             `json:"MatchedNetSum,omitempty"`
	AllocatedQuantity         float64                             `json:"AllocatedQuantity,omitempty"`
	LastComment               string                              `json:"LastComment,omitempty" validate:"max=1000"`
	WorkOrderCode             string                              `json:"WorkOrderCode,omitempty" validate:"max=250"`
	ConversionDeNumerator     float64                             `json:"ConversionDeNumerator,omitempty"`
	OrderLineUOM              string                              `json:"OrderLineUOM,omitempty" validate:"max=250"`
	ProjectSubName            string                              `json:"ProjectSubName,omitempty" validate:"max=250"`
	NetworkActivity           string                              `json:"NetworkActivity,omitempty" validate:"max=250"`
	BusinessAreaName          string                              `json:"BusinessAreaName,omitempty" validate:"max=250"`
	Num4                      float64                             `json:"Num4,omitempty"`
	DimCode6                  string                              `json:"DimCode6,omitempty" validate:"max=250"`
	ExternalCode              string                              `json:"ExternalCode" validate:"min=1,max=36"`
	ReceivedNetTotal          float64                             `json:"ReceivedNetTotal,omitempty"`
	WorkOrderSubCode          string                              `json:"WorkOrderSubCode,omitempty" validate:"max=250"`
	AccAssignmentCategoryName string                              `json:"AccAssignmentCategoryName,omitempty" validate:"max=250"`
	DimCode4                  string                              `json:"DimCode4,omitempty" validate:"max=250"`
	DeliveryNoteNumber        string                              `json:"DeliveryNoteNumber,omitempty" validate:"max=250"`
	OrderNumber               string                              `json:"OrderNumber,omitempty" validate:"max=1000"`
	DimName10                 string                              `json:"DimName10,omitempty" validate:"max=250"`
	PlannedAdditionalCostType int                                 `json:"PlannedAdditionalCostType"`
	FixedAssetName            string                              `json:"FixedAssetName,omitempty" validate:"max=250"`
	ServiceName               string                              `json:"ServiceName,omitempty" validate:"max=250"`
	AccAssignmentCategoryCode string                              `json:"AccAssignmentCategoryCode,omitempty" validate:"max=250"`
	OrderLineDescription      string                              `json:"OrderLineDescription,omitempty" validate:"max=250"`
	SalesOrderSubCode         string                              `json:"SalesOrderSubCode,omitempty" validate:"max=250"`
	OrderCodingRowNumber      string                              `json:"OrderCodingRowNumber,omitempty" validate:"max=250"`
	PartnerProfitCenter       string                              `json:"PartnerProfitCenter,omitempty" validate:"max=250"`
	FunctionalArea            string                              `json:"FunctionalArea,omitempty" validate:"max=250"`
	GrossTotalCompany         float64                             `json:"GrossTotalCompany,omitempty"`
	DimName4                  string                              `json:"DimName4,omitempty" validate:"max=250"`
	DimCode9                  string                              `json:"DimCode9,omitempty" validate:"max=250"`
	CustomerCode              string                              `json:"CustomerCode,omitempty" validate:"max=250"`
	DimCode7                  string                              `json:"DimCode7,omitempty" validate:"max=250"`
	DimName2                  string                              `json:"DimName2,omitempty" validate:"max=250"`
	DimName3                  string                              `json:"DimName3,omitempty" validate:"max=250"`
	ReceivedGrossPrice        float64                             `json:"ReceivedGrossPrice,omitempty"`
	VehicleName               string                              `json:"VehicleName,omitempty" validate:"max=250"`
	DimName7                  string                              `json:"DimName7,omitempty" validate:"max=250"`
	GoodsReceiptNumber        string                              `json:"GoodsReceiptNumber,omitempty" validate:"max=250"`
	BusinessAreaCode          string                              `json:"BusinessAreaCode,omitempty" validate:"max=250"`
	NetTotalOrganization      float64                             `json:"NetTotalOrganization,omitempty"`
	ProjectSubCode            string                              `json:"ProjectSubCode,omitempty" validate:"max=250"`
	SubUOM                    string                              `json:"SubUOM,omitempty" validate:"max=250"`
	OrderedGrossPrice         float64                             `json:"OrderedGrossPrice,omitempty"`
	AccountName               string                              `json:"AccountName,omitempty" validate:"max=250"`
	ReceivedGrossTotal        float64                             `json:"ReceivedGrossTotal,omitempty"`
	DimName8                  string                              `json:"DimName8,omitempty" validate:"max=250"`
	RowIndex                  int                                 `json:"RowIndex,omitempty"`
	OrderedNetPrice           float64                             `json:"OrderedNetPrice,omitempty"`
	CustomerName              string                              `json:"CustomerName,omitempty" validate:"max=250"`
	OrderNetTotal             float64                             `json:"OrderNetTotal,omitempty"`
	Text5                     string                              `json:"Text5,omitempty" validate:"max=250"`
	FreightSlip               string                              `json:"FreightSlip,omitempty"`
	CostCenterCode            string                              `json:"CostCenterCode,omitempty" validate:"max=200"`
	OrderLinePriceUnitDescription string                              `json:"OrderLinePriceUnitDescription,omitempty" validate:"max=250"`
	OrderLineNumber           string                              `json:"OrderLineNumber,omitempty" validate:"max=2000"`
	ConditionType             string                              `json:"ConditionType,omitempty" validate:"max=50"`
	RowOrigin                 string                              `json:"RowOrigin,omitempty" validate:"max=50"`
}
