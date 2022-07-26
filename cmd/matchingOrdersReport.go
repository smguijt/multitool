/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"multiTool/components"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

var matchingOrdersToProcess []string

// matchingOrdersReportCmd represents the matchingOrdersReport command
var matchingOrdersReportCmd = &cobra.Command{
	Use:   "matchingOrders",
	Short: "Retrieves all data from the matchingOrder API and MatchingOrderLines API. Creates Excel document of the content",
	Long: `Retrieves all data from the matchingOrder API and MatchingOrderLines API. Creates Excel document of the content. 
			Data is written into two Excel documents under the .multiTool folder`,
	Run: func(cmd *cobra.Command, args []string) {

		// create Excel document
		xlsx, _ := matchingOrdersCmd_CreateExcelDocument()
		xlsx, _ = matchingOrderLinesCmd_CreateExcelDocument(xlsx)

		// get headers
		processMatchingOrders(xlsx)

		// get lines
		for _, orderExternalCode := range matchingOrdersToProcess {
			if len(orderExternalCode) > 0 {
				processMatchingOrderLines(orderExternalCode, xlsx)
			}
		}

		// save excel document
		matchingOrdersCmd_SaveExcelDocument(xlsx, 1, "")
		//matchingOrderLinesCmd_SaveExcelDocument(xlsx, sheetno, "")

	},
}

func processMatchingOrders(xlsx *excelize.File) {
	config.EndpointMethod = "/v1/matchingOrders"
	config.Scope = "matchingOrders.read"

	//matchingOrders.read matchingOrders.write matchingOrders.delete

	switch config.AuthMethod {
	case BASIC:
		var continueationToken2BasicAuth string
		iRowBasicAuth := 0
		//var xlsxBasicAuth *excelize.File
		//var indexBasicAuth int

		iRowBasicAuth = 2
		for {
			jsonBasicAuth, continueationTokenBasicAuth, errMatchingOrdersCmdBasicAuthGetMatchingOrders := matchingOrdersCmd_Basic_Authentication_GetMatchingOrders(continueationToken2BasicAuth)
			if errMatchingOrdersCmdBasicAuthGetMatchingOrders == nil {
				// create Excel Document (report)
				//if iRowBasicAuth == 2 {
				//	xlsxBasicAuth, indexBasicAuth = matchingOrdersCmd_CreateExcelDocument()
				//}

				//create rows
				xlsx, iRowBasicAuth = matchingOrderCmd_AddRecord(jsonBasicAuth, xlsx, iRowBasicAuth)

				if continueationTokenBasicAuth == "" {
					// set active sheet and save document
					//matchingOrdersCmd_SaveExcelDocument(xlsxBasicAuth, indexBasicAuth, "")

					// end loop
					break

				} else {
					continueationToken2BasicAuth = continueationTokenBasicAuth
					log.Printf("continueationToken => %s ", continueationToken2BasicAuth)
					log.Println("fetch next ...")
				}
			}
		}

	case OAUTH2:
		matchingOrdersCmd_OAuth02_Authentication()
		if len(config.Token) > 0 {
			// init
			var continueationToken2 string
			iRow := 0
			//var xlsx *excelize.File
			//var index int

			// fetch
			iRow = 2
			for {
				json, continueationToken, errMatchingOrdersCmdOAuth02GetMatchingOrders := matchingOrdersCmd_OAuth02_GetMatchingOrders(continueationToken2)
				if errMatchingOrdersCmdOAuth02GetMatchingOrders == nil {
					// create Excel Document (report)
					//if iRow == 2 {
					//	xlsx, index = matchingOrdersCmd_CreateExcelDocument()
					//}

					//create rows
					xlsx, iRow = matchingOrderCmd_AddRecord(json, xlsx, iRow)

					if continueationToken == "" {
						// set active sheet and save document
						//matchingOrdersCmd_SaveExcelDocument(xlsx, index, "")

						// end loop
						break

					} else {
						continueationToken2 = continueationToken
						iRow++
						log.Printf("continueationToken => %s ", continueationToken2)
						log.Println("fetch next ...")
					}
				}
			}
		}
	}
}

func processMatchingOrderLines(orderExternalCode string, xlsx *excelize.File) {
	config.EndpointMethod = "/v1/matchingOrderLines"
	config.Scope = "matchingOrderLines.read"

	//matchingOrderLines.read matchingOrderLines.write matchingOrderLines.delete

	switch config.AuthMethod {
	case BASIC:
		var continueationToken2BasicAuth string
		iRowBasicAuth := 0
		//var xlsxBasicAuth *excelize.File
		//var indexBasicAuth int

		iRowBasicAuth = 2
		for {
			jsonBasicAuth, continueationTokenBasicAuth, errmatchingOrderLinesCmdBasicAuthGetMatchingOrders := matchingOrderLinesCmd_Basic_Authentication_GetMatchingOrderLines(orderExternalCode, continueationToken2BasicAuth)
			if errmatchingOrderLinesCmdBasicAuthGetMatchingOrders == nil {
				// create Excel Document (report)
				//if iRowBasicAuth == 2 {
				//	xlsxBasicAuth, indexBasicAuth = matchingOrderLinesCmd_CreateExcelDocument()
				//}

				//create rows
				xlsx, iRowBasicAuth = matchingOrderLineCmd_AddRecord(jsonBasicAuth, xlsx, iRowBasicAuth)

				if continueationTokenBasicAuth == "" {
					// set active sheet and save document
					//matchingOrderLinesCmd_SaveExcelDocument(xlsxBasicAuth, indexBasicAuth, "")

					// end loop
					break

				} else {
					continueationToken2BasicAuth = continueationTokenBasicAuth
					log.Printf("continueationToken => %s ", continueationToken2BasicAuth)
					log.Println("fetch next ...")
				}
			}
		}

	case OAUTH2:
		matchingOrdersCmd_OAuth02_Authentication()
		if len(config.Token) > 0 {
			// init
			var continueationToken2 string
			iRow := 0
			//var xlsx *excelize.File
			//var index int

			// fetch
			iRow = 2
			for {
				json, continueationToken, errmatchingOrderLinesCmdOAuth02GetMatchingOrderLines := matchingOrdersCmd_OAuth02_GetMatchingOrderLines(orderExternalCode, continueationToken2)
				if errmatchingOrderLinesCmdOAuth02GetMatchingOrderLines == nil {
					// create Excel Document (report)
					//if iRow == 2 {
					//	xlsx, index = matchingOrderLinesCmd_CreateExcelDocument()
					//}

					//create rows
					xlsx, iRow = matchingOrderLineCmd_AddRecord(json, xlsx, iRow)

					if continueationToken == "" {
						// set active sheet and save document
						//matchingOrderLinesCmd_SaveExcelDocument(xlsx, index, "")

						// end loop
						break

					} else {
						continueationToken2 = continueationToken
						iRow++
						log.Printf("continueationToken => %s ", continueationToken2)
						log.Println("fetch next ...")
					}
				}
			}
		}
	}
}

func init() {
	reportCmd.AddCommand(matchingOrdersReportCmd)
}

/* Helpers */
func matchingOrdersCmd_writeJson(response string, continueationToken string, externalCode string) {
	// create file
	var fileName string
	home, _ := homedir.Dir()

	fileName = config.EndpointMethod

	if continueationToken != "" {
		fileName += "_" + continueationToken
	}

	if externalCode != "" {
		fileName += "_" + externalCode
	}

	fileName += ".json"

	fileName = strings.Replace(fileName, "/v1/", "", 1)
	file, errHomeDir := os.Create(home + "/.multiTool/" + fileName)
	if errHomeDir != nil {
		log.Fatal(errHomeDir)
	}

	// init buffer
	writer := bufio.NewWriterSize(file, 10)

	_, errWriter := writer.WriteString(response)
	if errWriter != nil {
		log.Fatalf("Got error while writing to a file. Err: %s", errWriter.Error())
	}

	// flush and close file
	writer.Flush()

	log.Printf("%v", "retrieved result written to Home directory folder.")
}

func matchingOrdersCmd_AddCol(xlsx *excelize.File, sheetName string, colIndex int, rowIndex int, caption interface{}) (*excelize.File, int) {

	// add column
	colName, colErr := excelize.CoordinatesToCellName(colIndex, rowIndex)
	if colErr == nil {
		xlsx.SetCellValue(sheetName, colName, caption)
		newIndex := colIndex + 1
		return xlsx, newIndex
	} else {
		log.Fatalf("%v", colErr.Error())
		return xlsx, colIndex
	}
}

func matchingOrdersCmd_CreateExcelDocument() (*excelize.File, int) {

	var index int
	var sheetName string

	// create Excel Document (report)
	xlsx := excelize.NewFile()

	// create Sheet
	sheetName = "matchingOrder"
	xlsx, _ = matchingOrdersCmd_CreateMatchingOrderDocumentSheet(xlsx, sheetName)

	return xlsx, index
}

func matchingOrderLinesCmd_CreateExcelDocument(xlsx *excelize.File) (*excelize.File, int) {

	var index int
	var sheetName string

	// create Excel Document (report)
	//xlsx := excelize.NewFile()

	// create Sheet
	sheetName = "matchingOrderLines"
	xlsx, _ = matchingOrdersCmd_CreateMatchingOrderLinesDocumentSheet(xlsx, sheetName)

	// create Sheet
	sheetName = "orderLineCodings"
	xlsx, _ = matchingOrdersCmd_CreateMatchingOrderLineCodingSheet(xlsx, sheetName)

	// create Sheet
	sheetName = "goodsReceipts"
	xlsx, _ = matchingOrdersCmd_CreateMatchingOrderLineGoodsreceiptsSheet(xlsx, sheetName)

	// create Sheet
	sheetName = "referenceUsers"
	xlsx, _ = matchingOrdersCmd_CreateMatchingOrderLineReferenceUserSheet(xlsx, sheetName)

	return xlsx, index
}

func matchingOrdersCmd_SaveExcelDocument(xlsx *excelize.File, index int, continueationToken string) {

	home, _ := homedir.Dir()
	homeDir := home + "/.multiTool/"
	var filename string

	if continueationToken == "" {
		filename = "matchingOrders.xlsx"
	} else {
		filename = "matchingOrders_" + continueationToken + ".xlsx"
	}
	prefix := config.Prefix
	if len(prefix) > 0 {
		filename = prefix + "_" + filename
	}

	// set active sheet and save document
	xlsx.SetActiveSheet(1)
	xlsx.DeleteSheet("Sheet1")

	// resize matchingOrders
	var errAutoResize error
	xlsx, errAutoResize = matchingOrdersCmd_AutoResizeExcel("matchingOrder", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "matchingOrder", errAutoResize.Error())
	}

	// resize matchingOrderLines
	xlsx, errAutoResize = matchingOrdersCmd_AutoResizeExcel("matchingOrderLines", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "matchingOrderLines", errAutoResize.Error())
	}

	// resize matchingOrders
	xlsx, errAutoResize = matchingOrdersCmd_AutoResizeExcel("orderLineCodings", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "orderLineCodings", errAutoResize.Error())
	}

	// resize matchingOrders
	xlsx, errAutoResize = matchingOrdersCmd_AutoResizeExcel("goodsReceipts", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "goodsReceipts", errAutoResize.Error())
	}

	// resize matchingOrders
	xlsx, errAutoResize = matchingOrdersCmd_AutoResizeExcel("referenceUsers", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "referenceUsers", errAutoResize.Error())
	}

	// save document
	errDocumentsExcel := xlsx.SaveAs(homeDir + filename)
	if errDocumentsExcel != nil {
		log.Println(errDocumentsExcel)
	} else {
		log.Println("Excel document generated!!")
	}
}

func matchingOrdersCmd_AutoResizeExcel(sheetName string, xlsx *excelize.File) (*excelize.File, error) {
	// Autofit all columns according to their text content
	cols, err := xlsx.GetCols(sheetName)
	if err != nil {
		return xlsx, err
	}

	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return xlsx, err
		}
		xlsx.SetColWidth(sheetName, name, name, float64(largestWidth))
	}

	return xlsx, nil
}

/* ADD RECORD ENTRIES */
func matchingOrderCmd_AddRecord(json components.OrderResponse, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	var iRowIndex int
	for _, rec := range json.MatchingOrders {
		// Process Header Records
		iRowIndex = iRow
		xlsx, _ = matchingOrdersCmd_AddMatchingOrdersRecord(rec, xlsx, iRowIndex)

		//increment row number
		iRow++
	}

	return xlsx, iRow
}

func matchingOrderLineCmd_AddRecord(json []components.OrderLineEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	var iRowIndex int
	for _, rec := range json {
		// Process Line Records
		iRowIndex = iRow
		xlsx, _ = matchingOrdersCmd_AddMatchingOrderLinesRecord(rec, xlsx, iRowIndex)

		// add coding line
		for _, recLineCoding := range rec.OrderLineCoding {
			xlsx, _ = matchingOrdersCmd_AddMatchingOrderLineCodingsRecord(recLineCoding, xlsx, iRowIndex)
		}

		for _, recGoodsReceipt := range rec.GoodsReceipts {
			xlsx, _ = matchingOrdersCmd_AddMatchingGoodsReceiptsRecord(recGoodsReceipt, xlsx, iRowIndex)
		}

		for _, recRefUsers := range rec.ReferenceUsers {
			xlsx, _ = matchingOrdersCmd_AddMatchingOrderLineReferenceUsersRecord(recRefUsers, xlsx, iRowIndex)
		}
		//increment row number
		iRow++
	}

	return xlsx, iRow
}

func matchingOrdersCmd_AddMatchingOrdersRecord(json components.OrderEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "matchingOrder"
	colindex := 1

	//add order to slice (required for line processing)
	matchingOrdersToProcess = append(matchingOrdersToProcess, json.ExternalCode)

	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.OrderNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SourceSystem)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CompanyCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CompanyName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.OrganizationElementCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.OrganizationElementName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.OrderTypeCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CurrencyCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Created)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsClosed)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Description)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PurchaseOrganizationCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PurchaseOrganizationName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SupplierCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SupplierName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicingSupplierCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicingSupplierName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.RequestedDeliveryDate)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ActualDeliveryDate)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ValidFrom)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ValidTo)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PaymentTermCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PaymentTermName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text6)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text7)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text8)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text9)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text10)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsInvoiced)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsDelivered)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.LastUpdated)
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.OrderReference)

	return xlsx, iRow
}

func matchingOrdersCmd_AddMatchingOrderLinesRecord(json components.OrderLineEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "matchingOrderLines"
	colindex := 1

	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.OrderExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.LineNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SortNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Quantity)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.NetSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GrossSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CurrencyCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.MatchingMode) //standard, blanket, return
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsReceiptRequired)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsReceiptBasedMatching)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsOverreceivalAllowed)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsClosed)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsDeleted)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsSelfApproved)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Uom)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SubUOM)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.NetPrice)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GrossPrice)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PriceUnit)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PriceUnitDescription)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxPercent)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxPercent2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxSum2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicedQuantity)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicedNetSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicedGrossSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ValidFrom)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ValidTo)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProductCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProductName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.MaterialGroup)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GlobalTradeItemNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Unspsc)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BuyerProductCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ContractNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Description)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Comment)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.RequestedDeliveryDate)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ActualDeliveryDate)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text6)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text7)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text8)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text9)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text10)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date5)
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.LastUpdated)

	return xlsx, iRow
}

//x-lint:ignore U1000 Ignore unused function temporarily for debugging
func matchingOrdersCmd_AddMatchingOrderLineCodingsRecord(json components.OrderLineCodingEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "orderLineCodings"
	colindex := 1

	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.RowIndex)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.NetTotal)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GrossTotal)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.AccountCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.AccountName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CostCenterCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CostCenterName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProjectCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProjectName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ConversionNumerator)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ConversionDenominator)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ConversionDeNumerator)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxPercent)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxPercent2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxSum2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.PartnerProfitCenter)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.FixedAssetCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.FixedAssetName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.MaterialGroup)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.FixedAssetSubCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.FixedAssetSubName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InternalOrderCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InternalOrderName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProfitCenterCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProfitCenterName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BusinessUnitCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BusinessUnitName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProjectSubCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProjectSubName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.EmployeeCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.EmployeeName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.VehicleNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.VehicleName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SalesOrderCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SalesOrderName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SalesOrderSubCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SalesOrderSubName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CustomerCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CustomerName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.AccAssignmentCategoryCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.AccAssignmentCategoryName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BudgetCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BudgetName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ServiceCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ServiceName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BusinessAreaCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BusinessAreaName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.TaxJurisdictionCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SubUOM)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.WorkOrderCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.WorkOrderName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.WorkOrderSubCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.WorkOrderSubName)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.CommitmentItem)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ControllingArea)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.FunctionalArea)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode6)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode7)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode8)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode9)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimCode10)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName6)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName7)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName8)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName9)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DimName10)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Num1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Num2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Num3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Num4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Num5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.AllocatedQuantity)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Network)
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.NetworkActivity)

	return xlsx, iRow
}

//x-lint:ignore U1000 Ignore unused function temporarily for debugging
func matchingOrdersCmd_AddMatchingGoodsReceiptsRecord(json components.GoodsReceiptEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "goodsReceipts"
	colindex := 1

	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GoodsReceiptNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GoodsReceiptLineNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ReferenceGRExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DeliveryNoteNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.BestFitGrouping)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Quantity)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.NetSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GrossSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.NetPrice)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GrossPrice)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.IsDeleted)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicedQuantity)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicedNetSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.InvoicedGrossSum)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.UnitOfMeasure)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.SubUnitOfMeasure)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GoodsReceiptType)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ReceiveMethod)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.VoucherNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.GoodsReceiptNote)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.FiscalYear)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.DeliveryDate)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ProductSerialNumber)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text6)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text7)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text8)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text9)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Text10)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric4)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Numeric5)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date1)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date2)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date3)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date4)
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.Date5)

	return xlsx, iRow
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func matchingOrdersCmd_AddMatchingOrderLineReferenceUsersRecord(json components.OrderLineUserEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "referenceUsers"
	colindex := 1

	//xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.ExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.UserExternalCode)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.UserEmail)
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.UserRole) // Buyer, Owner, ReferencePerson, Other
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, iRow, json.LastUpdated)

	return xlsx, iRow
}

/* CREATE SHEETS AND ADD HEADERS */
func matchingOrdersCmd_CreateMatchingOrderDocumentSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	colindex := 1
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "orderNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "sourceSystem")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "companyCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "companyName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationElementCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationElementName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "orderTypeCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "created")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isClosed")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "purchaseOrganizationCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "purchaseOrganizationName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicingSupplierCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicingSupplierName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "requestedDeliveryDate")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "actualDeliveryDate")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "validFrom")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "validTo")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentTermCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentTermName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text6")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text7")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text8")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text9")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text10")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isInvoiced")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isDelivered")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "lastUpdated")
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "orderReference")

	return xlsx, sheetIndex
}

func matchingOrdersCmd_CreateMatchingOrderLinesDocumentSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	colindex := 1
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "orderExternalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "lineNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "sortNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "quantity")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "netSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "grossSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "matchingMode") //standard, blanket, return
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isReceiptRequired")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isReceiptBasedMatching")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isOverreceivalAllowed")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isClosed")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isDeleted")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isSelfApproved")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "uom")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "subUOM")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "netPrice")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "grossPrice")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "priceUnit")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "priceUnitDescription")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicedQuantity")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicedNetSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicedGrossSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "validFrom")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "validTo")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "productCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "productName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "materialGroup")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "globalTradeItemNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "unspsc")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "buyerProductCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "contractNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "comment")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "requestedDeliveryDate")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "actualDeliveryDate")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text6")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text7")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text8")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text9")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text10")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date5")
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "lastUpdated")

	return xlsx, sheetIndex
}

func matchingOrdersCmd_CreateMatchingOrderLineCodingSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	colindex := 1
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "rowIndex")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "netTotal")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "grossTotal")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "accountCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "accountName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "costCenterCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "costCenterName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "projectCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "projectName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "conversionNumerator")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "conversionDenominator")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "conversionDeNumerator")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "partnerProfitCenter")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "materialGroup")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetSubCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetSubName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "internalOrderCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "internalOrderName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "profitCenterCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "profitCenterName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "businessUnitCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "businessUnitName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "projectSubCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "projectSubName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "employeeCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "employeeName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "vehicleNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "vehicleName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderSubCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderSubName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "customerCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "customerName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "accAssignmentCategoryCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "accAssignmentCategoryName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "budgetCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "budgetName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "serviceCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "serviceName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "businessAreaCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "businessAreaName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "taxJurisdictionCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "subUOM")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderSubCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderSubName")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "commitmentItem")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "controllingArea")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "functionalArea")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode6")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode7")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode8")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode9")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode10")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName6")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName7")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName8")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName9")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName10")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "num1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "num2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "num3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "num4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "num5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "allocatedQuantity")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "network")
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "networkActivity")

	return xlsx, sheetIndex
}

func matchingOrdersCmd_CreateMatchingOrderLineGoodsreceiptsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	colindex := 1
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "goodsReceiptNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "goodsReceiptLineNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "referenceGRExternalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "deliveryNoteNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "bestFitGrouping")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "quantity")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "netSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "grossSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "netPrice")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "grossPrice")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "isDeleted")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicedQuantity")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicedNetSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "invoicedGrossSum")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "unitOfMeasure")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "subUnitOfMeasure")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "goodsReceiptType")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "receiveMethod")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "goodsReceiptNote")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "fiscalYear")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "deliveryDate")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "productSerialNumber")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text6")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text7")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text8")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text9")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "text10")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric4")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric5")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date1")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date2")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date3")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date4")
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "date5")

	return xlsx, sheetIndex
}

func matchingOrdersCmd_CreateMatchingOrderLineReferenceUserSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	colindex := 1
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "userExternalCode")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "userEmail")
	xlsx, colindex = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "userRole") // Buyer, Owner, ReferencePerson, Other
	xlsx, _ = matchingOrdersCmd_AddCol(xlsx, sheetName, colindex, 1, "lastUpdated")

	return xlsx, sheetIndex
}

/* DATA FETCH */
func matchingOrdersCmd_Basic_Authentication_GetMatchingOrders(continueationToken string) (components.OrderResponse, string, error) {
	log.Println("Authentication method => BASIC")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Scope => %v", config.Scope)
	log.Printf("Endpoint method => %v", config.EndpointMethod)
	log.Printf("pageSize => %v", config.PageSize)

	if config.PageSize == 0 {
		viper.Set("api_pagesize", 100)
		_ = viper.WriteConfig()
	}

	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/matchingOrders?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}
	req.SetBasicAuth(config.Username, config.Password)

	// add continuation token to fetch next page result
	if continueationToken != "" {
		req.Header.Add("x-amz-meta-continuationtoken", continueationToken)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	} else {
		log.Println("basicauth response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		matchingOrdersCmd_writeJson(string([]byte(body)), continueationToken, "")
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	tmpW := components.OrderResponse{}
	errDocumentResponse := json.Unmarshal(body, &tmpW)
	if errDocumentResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errDocumentResponse.Error())
		}
		return components.OrderResponse{}, "", errDocumentResponse
	}

	return tmpW, continueationToken, nil
}

func matchingOrderLinesCmd_Basic_Authentication_GetMatchingOrderLines(orderExternalCode string, continueationToken string) ([]components.OrderLineEntity, string, error) {
	log.Println("Authentication method => BASIC")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Scope => %v", config.Scope)
	log.Printf("Endpoint method => %v", config.EndpointMethod)
	log.Printf("pageSize => %v", config.PageSize)

	if config.PageSize == 0 {
		viper.Set("api_pagesize", 100)
		_ = viper.WriteConfig()
	}

	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/matchingOrderLines?OrderExternalCode="+orderExternalCode+"&pageSize="+strconv.Itoa(config.PageSize)+"", nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}
	req.SetBasicAuth(config.Username, config.Password)

	// add continuation token to fetch next page result
	if continueationToken != "" {
		req.Header.Add("x-amz-meta-continuationtoken", continueationToken)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	} else {
		log.Println("basicauth response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		matchingOrdersCmd_writeJson(string([]byte(body)), continueationToken, orderExternalCode)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	//tmpW := components.OrderLineResponse{}
	tmpW := []components.OrderLineEntity{}
	errDocumentResponse := json.Unmarshal(body, &tmpW)
	if errDocumentResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errDocumentResponse.Error())
		}
		//return components.OrderLineResponse{}, "", errDocumentResponse
		return []components.OrderLineEntity{}, "", errDocumentResponse
	}

	return tmpW, continueationToken, nil
}

/* OAUTH02 */
func matchingOrdersCmd_OAuth02_Authentication() {
	log.Println("Authentication method => OAUTH2")

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Scope => %v", config.Scope)
	log.Printf("Endpoint method => %v", config.EndpointMethod)
	log.Printf("pageSize => %v", config.PageSize)

	if config.PageSize == 0 {
		viper.Set("api_pagesize", 100)
		_ = viper.WriteConfig()
	}

	// create form data to post
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	data.Set("scope", config.Scope)

	// create POST request header and link form data
	req, err := http.NewRequest("POST", config.EndpointUrl+"v1/tokens", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}

	// add authentication header
	req.SetBasicAuth(config.ClientId, config.ClientSecret)

	// add header params (form urlencoded and length)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// send request
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	} else {
		log.Println("oauth2 response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// process response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(string(body))

	var accessToken AccessTokenResponse
	err2 := json.Unmarshal([]byte(string(body)), &accessToken)
	if err2 != nil {
		log.Fatal("AccessTokenResponse Decode: ", err2)
	} else {
		config.Token = accessToken.Access_Token
		viper.Set("api_token", config.Token)
		_ = viper.WriteConfig()
	}
}

func matchingOrdersCmd_OAuth02_GetMatchingOrders(continueationToken string) (components.OrderResponse, string, error) {

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// get account
	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/matchingOrders?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}

	// add authentication header
	bearer := "Bearer " + config.Token
	req.Header.Add("Authorization", bearer)

	// add continuation token to fetch next page result
	if continueationToken != "" {
		req.Header.Add("x-amz-meta-continuationtoken", continueationToken)
	}

	// send request
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	} else {
		log.Println("getMatchingOrders - response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		matchingOrdersCmd_writeJson(string([]byte(body)), continueationToken, "")
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	tmpW := components.OrderResponse{}
	errDocumentsResponse := json.Unmarshal(body, &tmpW)
	if errDocumentsResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errDocumentsResponse.Error())
		}
		return components.OrderResponse{}, "", errDocumentsResponse
	}

	return tmpW, continueationToken, nil
}

func matchingOrdersCmd_OAuth02_GetMatchingOrderLines(orderExternalCode string, continueationToken string) ([]components.OrderLineEntity, string, error) {

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// get account
	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/matchingOrderLines?OrderExternalCode="+orderExternalCode+"&pageSize="+strconv.Itoa(config.PageSize)+"", nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}

	// add authentication header
	bearer := "Bearer " + config.Token
	req.Header.Add("Authorization", bearer)

	// add continuation token to fetch next page result
	if continueationToken != "" {
		req.Header.Add("x-amz-meta-continuationtoken", continueationToken)
	}

	// send request
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	} else {
		log.Println("getMatchingOrderLines - response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		matchingOrdersCmd_writeJson(string([]byte(body)), continueationToken, orderExternalCode)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	//tmpW := components.OrderLineResponse{}
	tmpW := []components.OrderLineEntity{}
	errDocumentsResponse := json.Unmarshal(body, &tmpW)
	if errDocumentsResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errDocumentsResponse.Error())
		}
		//return components.OrderLineResponse{}, "", errDocumentsResponse
		return []components.OrderLineEntity{}, "", errDocumentsResponse
	}

	return tmpW, continueationToken, nil
}
