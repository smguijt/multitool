/*
Copyright © 2022 Sietse Guijt
https://github.com/qax-os/excelize

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

// accountingDocumentsCmd represents the accountingDocuments command
var accountingDocumentsCmd = &cobra.Command{
	Use:   "accountingDocuments",
	Short: "Retrieves all data from the accountingDocuments API and creates Excel document of the content",
	Long: `Retrieves all data from the accountingDocuments API. 
	Data is written into an Excel document under the .multiTool folder`,
	Run: func(cmd *cobra.Command, args []string) {

		config.EndpointMethod = "accountingDocuments"
		config.Scope = "accountingDocuments.read"

		switch config.AuthMethod {
		case BASIC:
			var continueationToken2BasicAuth string
			iRowBasicAuth := 0
			var xlsxBasicAuth *excelize.File
			var indexBasicAuth int

			// fetch
			iRowBasicAuth = 2
			for {
				jsonBasicAuth, continueationTokenBasicAuth, errAccountingDocumentsCmdBasicAuthGetAccountingDocuments := accountingDocumentsCmd_Basic_Authentication_GetAccountingDocuments(continueationToken2BasicAuth)
				if errAccountingDocumentsCmdBasicAuthGetAccountingDocuments == nil {
					// create Excel Document (report)
					if iRowBasicAuth == 2 {
						xlsxBasicAuth, indexBasicAuth = accountingDocumentsCmd_CreateExcelDocument()
					}

					//create rows
					xlsxBasicAuth, iRowBasicAuth = accountingDocumentsCmd_AddRecord(jsonBasicAuth, xlsxBasicAuth, iRowBasicAuth)

					if continueationTokenBasicAuth == "" {
						// set active sheet and save document
						accountingDocumentsCmd_SaveExcelDocument(xlsxBasicAuth, indexBasicAuth, "")

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
			accountingDocumentsCmd_OAuth02_Authentication()
			if len(config.Token) > 0 {
				// init
				var continueationToken2 string
				iRow := 0
				var xlsx *excelize.File
				var index int

				// fetch
				iRow = 2
				for {
					json, continueationToken, errAccountingDocumentsCmdOAuth02GetAccountingDocuments := accountingDocumentsCmd_OAuth02_GetAccountingDocuments(continueationToken2)
					if errAccountingDocumentsCmdOAuth02GetAccountingDocuments == nil {
						// create Excel Document (report)
						if iRow == 2 {
							xlsx, index = accountingDocumentsCmd_CreateExcelDocument()
						}

						//create rows
						xlsx, iRow = accountingDocumentsCmd_AddRecord(json, xlsx, iRow)

						if continueationToken == "" {
							// set active sheet and save document
							accountingDocumentsCmd_SaveExcelDocument(xlsx, index, "")

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
	},
}

func init() {
	reportCmd.AddCommand(accountingDocumentsCmd)
}

/* Core Functions */
/* */
func accountingDocumentsCmd_Basic_Authentication_GetAccountingDocuments(continueationToken string) (components.AccountingDocumentResponse, string, error) {
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

	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/accountingDocuments?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
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
		accountingDocumentsCmd_writeJson(string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	tmpW := components.AccountingDocumentResponse{}
	errAccountingDocumentResponse := json.Unmarshal(body, &tmpW)
	if errAccountingDocumentResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errAccountingDocumentResponse.Error())
		}
		return components.AccountingDocumentResponse{}, "", errAccountingDocumentResponse
	}

	return tmpW, continueationToken, nil
}

func accountingDocumentsCmd_OAuth02_Authentication() {
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

func accountingDocumentsCmd_OAuth02_GetAccountingDocuments(continueationToken string) (components.AccountingDocumentResponse, string, error) {

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// get account
	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/accountingDocuments?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
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
		log.Println("getAccountingDocuments - response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		accountingDocumentsCmd_writeJson(string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")
	// log.Printf("CT: %v", continueationToken)

	// unmarshal response to object
	tmpW := components.AccountingDocumentResponse{}
	errAccountingDocumentsResponse := json.Unmarshal(body, &tmpW)
	if errAccountingDocumentsResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errAccountingDocumentsResponse.Error())
		}
		return components.AccountingDocumentResponse{}, "", errAccountingDocumentsResponse
	}

	return tmpW, continueationToken, nil
}

func accountingDocumentsCmd_writeJson(response string, continueationToken string) {
	// create file
	var fileName string
	home, _ := homedir.Dir()
	if continueationToken == "" {
		fileName = config.EndpointMethod + ".json"
	} else {
		fileName = config.EndpointMethod + "_" + continueationToken + ".json"
	}

	file, errHomeDir := os.Create(home + "/.multiTool/" + fileName)
	if errHomeDir != nil {
		log.Fatal(errHomeDir)
	}

	// init buffer
	writer := bufio.NewWriterSize(file, 10)

	_, errWriter := writer.WriteString(response)
	// _, errWriter := writer.WriteString(PrettyPrint(response))
	if errWriter != nil {
		log.Fatalf("Got error while writing to a file. Err: %s", errWriter.Error())
	}

	// flush and close file
	writer.Flush()

	log.Printf("%v", "retrieved result written to Home directory folder.")
}

func accountingDocumentsCmd_AddCol(xlsx *excelize.File, sheetName string, colIndex int, rowIndex int, caption interface{}) (*excelize.File, int) {

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

func accountingDocumentsCmd_CreateExcelDocument() (*excelize.File, int) {

	var index int
	var sheetName string

	// create Excel Document (report)
	xlsx := excelize.NewFile()

	// create Sheet
	sheetName = "accountingDocument"
	xlsx, _ = accountingDocumentsCmd_CreateAccountingDocumentSheet(xlsx, sheetName)

	// create sheet
	sheetName = "orderNumbers"
	xlsx, _ = accountingDocumentsCmd_CreateOrderNumbersSheet(xlsx, sheetName)

	// create sheet
	sheetName = "codingRows"
	xlsx, _ = accountingDocumentsCmd_CreateCodingRowsSheet(xlsx, sheetName)

	// create sheet
	sheetName = "transferParameters"
	xlsx, _ = accountingDocumentsCmd_CreateTransferParametersSheet(xlsx, sheetName)

	// create sheet
	sheetName = "transferResponses"
	xlsx, _ = accountingDocumentsCmd_CreateTransferResponsesSheet(xlsx, sheetName)

	// create sheet
	sheetName = "removeResponses"
	xlsx, _ = accountingDocumentsCmd_CreateRemoveResponsesSheet(xlsx, sheetName)

	// create sheet
	sheetName = "prebookResponses"
	xlsx, _ = accountingDocumentsCmd_CreatePrebookResponsesSheet(xlsx, sheetName)

	// create sheet
	sheetName = "paymentResponses"
	xlsx, index = accountingDocumentsCmd_CreatePaymentResponsesSheet(xlsx, sheetName)

	return xlsx, index

}

func accountingDocumentsCmd_SaveExcelDocument(xlsx *excelize.File, index int, continueationToken string) {

	home, _ := homedir.Dir()
	homeDir := home + "/.multiTool/"
	var filename string

	if continueationToken == "" {
		filename = "accountingDocuments.xlsx"
	} else {
		filename = "accountingDocuments_" + continueationToken + ".xlsx"
	}
	prefix := config.Prefix
	if len(prefix) > 0 {
		filename = prefix + "_" + filename
	}

	// set active sheet and save document
	xlsx.SetActiveSheet(1)
	xlsx.DeleteSheet("Sheet1")

	// resize accountingDocuments
	var errAutoResize error
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("accountingDocument", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "accountingDocument", errAutoResize.Error())
	}

	// resize orderNumbers
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("orderNumbers", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "orderNumbers", errAutoResize.Error())
	}

	// resize codingRows
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("codingRows", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "codingRows", errAutoResize.Error())
	}

	// resize transferParameters
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("transferParameters", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "transferParameters", errAutoResize.Error())
	}

	// resize transferResponses
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("transferResponses", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "transferResponses", errAutoResize.Error())
	}

	// resize removeResponses
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("removeResponses", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "removeResponses", errAutoResize.Error())
	}

	// resize prebookResponses
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("prebookResponses", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "prebookResponses", errAutoResize.Error())
	}

	// resize paymentResponses
	xlsx, errAutoResize = accountingDocumentsCmd_AutoResizeExcel("paymentResponses", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "paymentResponses", errAutoResize.Error())
	}

	// save document
	errAccountingDocumentsExcel := xlsx.SaveAs(homeDir + filename)
	if errAccountingDocumentsExcel != nil {
		log.Println(errAccountingDocumentsExcel)
	} else {
		log.Println("Excel document generated!!")
	}
}

func accountingDocumentsCmd_AutoResizeExcel(sheetName string, xlsx *excelize.File) (*excelize.File, error) {
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

func accountingDocumentsCmd_AddRecord(json components.AccountingDocumentResponse, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	var iRowIndex int
	for _, rec := range json.AccountingDocuments {

		// Process Header Records
		iRowIndex = iRow
		xlsx, _ = accountingDocumentsCmd_AddAccountingDocumentRecord(rec, xlsx, iRowIndex)

		//* Processing Order numbers
		iRowIndex = iRow
		for _, ordernumber := range rec.OrderNumbers {
			xlsx, _ = accountingDocumentsCmd_AddOrderNumbersRecord(ordernumber, xlsx, iRowIndex)
			iRowIndex++
		}

		// Processing Coding rows
		iRowIndex = iRow
		for _, codingRow := range rec.CodingRows {
			xlsx, _ = accountingDocumentsCmd_AddCodingRowsRecord(codingRow, xlsx, iRowIndex)
			iRowIndex++
		}

		// Processing payment responses rows
		iRowIndex = iRow
		for _, paymentResponse := range rec.PaymentResponses {
			xlsx, _ = accountingDocumentsCmd_AddPaymentResponsesRecord(paymentResponse, xlsx, iRowIndex)
			iRowIndex++
		}

		// Processing transfer response rows
		iRowIndex = iRow
		for _, transferResponse := range rec.TransferResponses {
			xlsx, _ = accountingDocumentsCmd_AddTransferResponsesRecord(transferResponse, xlsx, iRowIndex)
			iRowIndex++
		}

		// Processing transfer parameters rows
		iRowIndex = iRow
		for _, transferParam := range rec.TransferParameters {
			xlsx, _ = accountingDocumentsCmd_AddTransferParametersRecord(transferParam, xlsx, iRowIndex)
			iRowIndex++
		}

		// Processing remove response rows
		iRowIndex = iRow
		for _, removeResponse := range rec.RemoveResponses {
			xlsx, _ = accountingDocumentsCmd_AddRemoveResponsesRecord(removeResponse, xlsx, iRowIndex)
			iRowIndex++
		}

		// Processing prebook response rows
		iRowIndex = iRow
		for _, prebookResponse := range rec.PrebookResponses {
			xlsx, _ = accountingDocumentsCmd_AddPrebookResponsesRecord(prebookResponse, xlsx, iRowIndex)
			iRowIndex++
		}

		iRow++
	}

	return xlsx, iRow
}

/* Create Excel Worksheets */
/* Helper Functions */
/* */
func accountingDocumentsCmd_CreateAccountingDocumentSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// add columns
	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceId")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "bumid")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "processingStatus")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "originService")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "companyCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "companyName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "netSum")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "grossSum")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "baseLineDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dueDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "cashDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "cashPercent")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "cashSum")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "referencePerson")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "contractNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceTypeCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceTypeName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentMethod")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentBlock")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentTermCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentTermName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentTermExternalCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentPlanReference")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentRevelsalDocument")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "codingDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "prebooked")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "referenceNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierBankIBAN")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierBankBBAN")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierBankBIC")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierBankName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountingPeriod")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountingGroup")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text6")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text7")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text8")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text9")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text10")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text11")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text12")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text13")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text14")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text15")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text16")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text17")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text18")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text19")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text20")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text21")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text22")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text23")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text24")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text25")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text26")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text27")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text28")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text29")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text30")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric6")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric7")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric8")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric9")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric10")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric11")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric12")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric13")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric14")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric15")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric16")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric17")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric18")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric19")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "numeric20")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date6")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date7")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date8")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date9")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date10")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceImageURL")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "invoiceImageToken")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "lastUpdated")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "cashSumCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "cashSumOrganization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "exchangeRateBaseDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "exchangeRateCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "exchangeRateOrganization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum1Company")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum1Organization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2Company")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2Organization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCodeCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCodeOrganization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationElementCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationElementName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "grossSumCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "grossSumOrganization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "netSumCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "netSumOrganization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "parentInvoiceBumId")

	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierSourceSystemId")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreateOrderNumbersSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderNumber")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreateCodingRowsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "rowIndex")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "costCenterCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "costCenterName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "deliveryNoteNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "lastComment")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "matchingType")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "grossTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "grossTotalCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "grossTotalOrganization")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "netTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "netTotalCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "netTotalOrganization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLineNetTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLineGrossTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "receivedQuantity")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "receivedNetPrice")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "receivedGrossPrice")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "receivedNetTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "receivedGrossTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationElementName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "organizationElementCode")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "freightSlip")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "conversionNumerator")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "conversionDeNumerator")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "conversionDenominator")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLineGrossTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderedQuantity")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderNetTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderGrossTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderedNetPrice")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderedGrossPrice")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "internalOrderCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "internalOrderName")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderItemNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLineNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "contractNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "allocatedQuantity")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "matchedQuantity")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "matchedNetSum")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "matchedGrossSum")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "closeOrder")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "plant")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "profitCenterCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "profitCenterName")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "businessUnitCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "businessUnitName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "projectCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "projectName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "projectSubCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "projectSubName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "employeeCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "employeeName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "vehicleNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "vehicleName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderName")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderSubCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "salesOrderSubName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "customerCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "customerName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "conditionType")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "accAssignmentCategoryCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "accAssignmentCategoryName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "budgetCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "budgetName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "serviceCode	")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "serviceName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "businessAreaCode")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "businessAreaName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "productCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "productName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "buyerName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLineDescription")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLinePriceUnit")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLinePriceUnitDescription")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "partnerProfitCenter")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetName	")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "goodsReceiptItemNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "goodsReceiptNumber")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "materialGroup")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "fiscalYear")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetSubCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "fixedAssetSubName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderLineUOM")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "subUOM")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderSubCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "workOrderSubName	")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "commitmentItem")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "controllingArea")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "functionalArea")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "network")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "networkActivity")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderCodingRowNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "ownerName")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxJurisdictionCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxPercent2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSumCompany")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSumOrganization")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2Company")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "taxSum2Organization")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "rowOrigin")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "plannedAdditionalCostType")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode6")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode7")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode8")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode9")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode10")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName6")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName7")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName8")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimName9")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "dimCode10")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "num1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "num2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "num3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "num4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "num5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text4")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "text5")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date1")

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date3")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date4")
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "date5")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreateTransferParametersSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "key")
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "value")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreateTransferResponsesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "success")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "responseMessage")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "sourceSystem")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "transferDate")
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentBlock")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreateRemoveResponsesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "success")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "responseMessage")
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "sourceSystem")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreatePrebookResponsesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "success")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "responseMessage")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "sourceSystem")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber1")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "voucherNumber2")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "prebookDate")
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentBlock")

	return xlsx, sheetIndex
}

func accountingDocumentsCmd_CreatePaymentResponsesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "success")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "responseMessage")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "sourceSystem")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentDate")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentTermCode")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentBlock")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentMethod")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paidTotal")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "cashDiscount")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "checkNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentMessage")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentNumber")
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentDocument")
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentReversalDocument")

	return xlsx, sheetIndex
}

/* Add Data to the Excel Worksheets */
/* Helper Functions */
/* */
func accountingDocumentsCmd_AddAccountingDocumentRecord(accountingDocumentReponse components.AccountingDocumentEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "accountingDocument"
	colindex := 1

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceId)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Bumid)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.VoucherNumber1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.VoucherNumber2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.VoucherDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ProcessingStatus)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.OriginService)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CompanyCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CompanyName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.OrganizationCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.OrganizationName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.SupplierCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.SupplierName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CurrencyCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.NetSum)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.GrossSum)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxPercent1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxPercent2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxSum1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxSum2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.BaseLineDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.DueDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CashDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CashPercent)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CashSum)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ReferencePerson)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ContractNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Description)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceTypeCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceTypeName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentMethod)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentBlock)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentTermCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentTermName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentTermExternalCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentPlanReference)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.PaymentRevelsalDocument)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CodingDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Prebooked)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ReferenceNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.SupplierBankIBAN)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.SupplierBankBBAN)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.SupplierBankBIC)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.SupplierBankName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.AccountingPeriod)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.AccountingGroup)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text6)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text7)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text8)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text9)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text10)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text11)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text12)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text13)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text14)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text15)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text16)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text17)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text18)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text19)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text20)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text21)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text22)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text23)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text24)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text25)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text26)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text27)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text28)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text29)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Text30)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric6)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric7)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric8)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric9)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric10)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric11)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric12)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric13)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric14)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric15)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric16)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric17)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric18)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric19)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Numeric20)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date6)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date7)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date8)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date9)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.Date10)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceImageURL)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.InvoiceImageToken)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.LastUpdated)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CashSumCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CashSumOrganization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ExchangeRateBaseDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ExchangeRateCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ExchangeRateOrganization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxSum1Company)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxSum1Organization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxSum2Company)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.TaxSum2Organization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CurrencyCodeCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.CurrencyCodeOrganization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.OrganizationElementCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.OrganizationElementName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.GrossSumCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.GrossSumOrganization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.NetSumCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.NetSumOrganization)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, accountingDocumentReponse.ParentInvoiceBumId)

	return xlsx, iRow
}

func accountingDocumentsCmd_AddOrderNumbersRecord(ordernumber string, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "orderNumber"
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, 1, iRow, ordernumber)
	return xlsx, iRow
}

func accountingDocumentsCmd_AddCodingRowsRecord(codingRow components.StandardCodingEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "codingRows"
	colindex := 1

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ExternalCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.RowIndex)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.AccountCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.AccountName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.CostCenterCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.CostCenterName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DeliveryNoteNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.LastComment)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.MatchingType)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.GrossTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.GrossTotalCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.GrossTotalOrganization)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.NetTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.NetTotalCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.NetTotalOrganization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLineNetTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLineGrossTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ReceivedQuantity)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ReceivedNetPrice)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ReceivedGrossPrice)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ReceivedNetTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ReceivedGrossTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrganizationElementName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrganizationElementCode)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FreightSlip)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ConversionNumerator)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ConversionDeNumerator)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ConversionDenominator)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLineGrossTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderedQuantity)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderNetTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderGrossTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderedNetPrice)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderedGrossPrice)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.InternalOrderCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.InternalOrderName)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderItemNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLineNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ContractNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.AllocatedQuantity)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.MatchedQuantity)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.MatchedNetSum)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.MatchedGrossSum)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.CloseOrder)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Plant)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProfitCenterCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProfitCenterName)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BusinessUnitCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BusinessUnitName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProjectCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProjectName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProjectSubCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProjectSubName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.EmployeeCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.EmployeeName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.VehicleNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.VehicleName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.SalesOrderCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.SalesOrderName)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.SalesOrderSubCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.SalesOrderSubName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.CustomerCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.CustomerName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ConditionType)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.AccAssignmentCategoryCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.AccAssignmentCategoryName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BudgetCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BudgetName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ServiceCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ServiceName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BusinessAreaCode)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BusinessAreaName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProductCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ProductName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.BuyerName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLineDescription)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLinePriceUnit)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLinePriceUnitDescription)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.PartnerProfitCenter)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FixedAssetCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FixedAssetName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.GoodsReceiptItemNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.GoodsReceiptNumber)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.MaterialGroup)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FiscalYear)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FixedAssetSubCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FixedAssetSubName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderLineUOM)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.SubUOM)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.WorkOrderCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.WorkOrderName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.WorkOrderSubCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.WorkOrderSubName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.CommitmentItem)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.ControllingArea)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.FunctionalArea)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Network)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.NetworkActivity)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OrderCodingRowNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.OwnerName)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxJurisdictionCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxPercent)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxPercent2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxSum)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxSumCompany)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxSumOrganization)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxSum2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxSum2Company)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.TaxSum2Organization)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.RowOrigin)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.PlannedAdditionalCostType)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode6)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode7)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode8)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode9)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode10)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName6)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName7)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName8)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimName9)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.DimCode10)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Num1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Num2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Num3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Num4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Num5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Text1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Text2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Text3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Text4)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Text5)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Date1)

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Date2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Date3)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Date4)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, codingRow.Date5)

	return xlsx, iRow
}

func accountingDocumentsCmd_AddTransferParametersRecord(params components.TransferParameterEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "transferParameters"
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, 1, iRow, params.Key)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, 2, iRow, params.Value)

	return xlsx, iRow
}

func accountingDocumentsCmd_AddTransferResponsesRecord(transferResponse components.TransferResponseEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "transferResponses"
	colindex := 1

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.ExternalCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.Success)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.ResponseMessage)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.SourceSystem)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.VoucherNumber1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.VoucherNumber2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.TransferDate)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, transferResponse.PaymentBlock)

	return xlsx, iRow
}

func accountingDocumentsCmd_AddRemoveResponsesRecord(removeResponse components.RemoveResponseEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "removeResponses"
	colindex := 1

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, removeResponse.ExternalCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, removeResponse.Success)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, removeResponse.ResponseMessage)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, removeResponse.SourceSystem)

	return xlsx, iRow
}

func accountingDocumentsCmd_AddPrebookResponsesRecord(prebookResponse components.PrebookResponseEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "prebookResponses"
	colindex := 1

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.ExternalCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.Success)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.ResponseMessage)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.SourceSystem)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.VoucherNumber1)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.VoucherNumber2)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.PrebookDate)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, prebookResponse.PaymentBlock)

	return xlsx, iRow
}

func accountingDocumentsCmd_AddPaymentResponsesRecord(paymentResponse components.PaymentResponseEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "paymentResponses"
	colindex := 1

	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.ExternalCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.Success)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.ResponseMessage)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.SourceSystem)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentDate)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentTermCode)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentBlock)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentMethod)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaidTotal)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.CashDiscount)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.CheckNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentMessage)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentNumber)
	xlsx, colindex = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentDocument)
	xlsx, _ = accountingDocumentsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentResponse.PaymentReversalDocument)

	return xlsx, iRow
}
