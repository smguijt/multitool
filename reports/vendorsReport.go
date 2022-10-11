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

// vendorsCmd represents the vendors command
var vendorsCmd = &cobra.Command{
	Use:   "vendors",
	Short: "Retrieves all data from the vendors API and creates Excel document of the content",
	Long: `Retrieves all data from the vendors API. 
	Data is written into an Excel document under the .multiTool folder`,
	Run: func(cmd *cobra.Command, args []string) {
		config.EndpointMethod = "vendors"
		config.Scope = "vendors.read"

		_externalCode, _ := cmd.Flags().GetString("externalCode")
		if _externalCode == "externalCode" {
			_externalCode = ""
		}
		log.Printf("ExternalCode: %s", _externalCode)

		switch config.AuthMethod {
		case BASIC:
			getVendorBasicData(_externalCode)
		case OAUTH2:
			getVendorSecuredData(_externalCode)
		}

	},
}

func init() {
	reportCmd.AddCommand(vendorsCmd)

	vendorsCmd.PersistentFlags().String("externalCode", "externalCode", "")
}

func getVendorBasicData(key string) {
	// var listArray []string
	var continueationToken2BasicAuth string
	var iRowBasicAuth int
	var xlsxBasicAuth *excelize.File
	var indexBasicAuth int

	iRowBasicAuth = 0
	continueationToken2BasicAuth = ""

	// fetch
	iRowBasicAuth = 2
	for {
		jsonBasicAuth, continueationTokenBasicAuth, errvendorsCmdBasicAuthVendorsDocuments := vendorsCmd_Basic_Authentication_GetVendorDocuments(continueationToken2BasicAuth, key)

		if errvendorsCmdBasicAuthVendorsDocuments == nil {
			if iRowBasicAuth == 2 {
				xlsxBasicAuth, indexBasicAuth = vendorsCmd_CreateExcelDocument()
			}

			//create rows
			xlsxBasicAuth, iRowBasicAuth = vendorsCmd_AddRecord(jsonBasicAuth, xlsxBasicAuth, iRowBasicAuth)

			if continueationTokenBasicAuth == "" {
				vendorsCmd_SaveExcelDocument(xlsxBasicAuth, indexBasicAuth, "", key)
				break

			} else {

				continueationToken2BasicAuth = continueationTokenBasicAuth
				log.Printf("continueationToken => %s ", continueationToken2BasicAuth)
				time.Sleep(10 * time.Millisecond)
				log.Println("fetch next ...")
			}
		} else {
			break
		}
	}
}

func getVendorSecuredData(key string) {
	vendorsCmd_OAuth02_AuthenticationToken()
	if len(config.Token) > 0 {
		// var listArray []string
		var continueationToken2OAuth02 string
		var iRowBasicAuth int
		var xlsxBasicAuth *excelize.File
		var indexBasicAuth int

		iRowBasicAuth = 0
		continueationToken2OAuth02 = ""

		// fetch
		iRowBasicAuth = 2
		for {
			jsonBasicAuth, continueationTokenOAuth02, errvendorsCmdOAuth02GetVendorDocuments := vendorsCmd_OAuth02_GetVendorDocuments(continueationToken2OAuth02, key)

			if errvendorsCmdOAuth02GetVendorDocuments == nil {
				// create Excel Document (report)
				if iRowBasicAuth == 2 {
					xlsxBasicAuth, indexBasicAuth = vendorsCmd_CreateExcelDocument()
				}

				//create rows
				xlsxBasicAuth, iRowBasicAuth = vendorsCmd_AddRecord(jsonBasicAuth, xlsxBasicAuth, iRowBasicAuth)

				if continueationTokenOAuth02 == "" {
					vendorsCmd_SaveExcelDocument(xlsxBasicAuth, indexBasicAuth, "", key)
					break

				} else {

					continueationToken2OAuth02 = continueationTokenOAuth02
					log.Printf("continueationToken => %s ", continueationToken2OAuth02)
					time.Sleep(10 * time.Millisecond)
					log.Println("fetch next ...")
				}
			} else {
				break
			}
		}
	}
}

func vendorsCmd_CreateExcelDocument() (*excelize.File, int) {

	var index int
	var sheetName string

	// create Excel Document (report)
	xlsx := excelize.NewFile()

	// create Sheet
	sheetName = "Vendors"
	xlsx, _ = vendorsCmd_CreateVendorsDocumentSheet(xlsx, sheetName)

	sheetName = "Identifiers"
	xlsx, _ = vendorsCmd_CreateIdentifiersSheet(xlsx, sheetName)

	sheetName = "Addresses"
	xlsx, _ = vendorsCmd_CreateAddressesSheet(xlsx, sheetName)

	sheetName = "AdditionalAddressFields"
	xlsx, _ = vendorsCmd_CreateAdditionalAddressFieldsSheet(xlsx, sheetName)

	sheetName = "Contacts"
	xlsx, _ = vendorsCmd_CreateContactsSheet(xlsx, sheetName)

	sheetName = "DeliveryTerms"
	xlsx, _ = vendorsCmd_CreateDeliveryTermsSheet(xlsx, sheetName)

	sheetName = "PaymentMeans"
	xlsx, _ = vendorsCmd_CreatePaymentMeansSheet(xlsx, sheetName)

	sheetName = "FinancialInstitutions"
	xlsx, _ = vendorsCmd_CreateFinancialInstitutionSheet(xlsx, sheetName)

	sheetName = "FinancialAccountIdentifiers"
	xlsx, _ = vendorsCmd_CreateFinancialAccountIdentifiersSheet(xlsx, sheetName)

	sheetName = "OrderingDetails"
	xlsx, _ = vendorsCmd_CreateOrderingDetailsSheet(xlsx, sheetName)

	sheetName = "ProcessingStatus"
	xlsx, _ = vendorsCmd_CreateProcessingStatusSheet(xlsx, sheetName)

	sheetName = "Tags"
	xlsx, _ = vendorsCmd_CreateTagsSheet(xlsx, sheetName)

	sheetName = "CustomFields"
	xlsx, _ = vendorsCmd_CreateCustomFieldsSheet(xlsx, sheetName)

	sheetName = "Companies"
	vendorsCmd_CreateCompaniesSheet(xlsx, sheetName)

	return xlsx, index
}

func vendorsCmd_SaveExcelDocument(xlsx *excelize.File, index int, continueationToken string, externalCode string) {

	home, _ := homedir.Dir()
	homeDir := home + "/.multiTool/"
	var filename string

	if continueationToken == "" {
		if externalCode == "" {
			filename = "vendors" + ".xlsx"
		} else {
			filename = "vendors" + "_" + externalCode + ".xlsx"
		}

	} else {
		if externalCode == "" {
			filename = "vendors" + "_" + continueationToken + ".xlsx"
		} else {
			filename = "vendors" + "_" + externalCode + "_" + continueationToken + ".xlsx"
		}
	}
	prefix := config.Prefix
	if len(prefix) > 0 {
		filename = prefix + "_" + filename
	}

	// set active sheet and save document
	xlsx.SetActiveSheet(1)
	xlsx.DeleteSheet("Sheet1")

	// resize list sheet
	var errAutoResize error
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("Vendors", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "Vendors", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("Identifiers", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "Identifiers", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("Addresses", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "Addresses", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("AdditionalAddressFields", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "AdditionalAddressFields", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("Contacts", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "Contacts", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("DeliveryTerms", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "DeliveryTerms", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("PaymentMeans", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "PaymentMeans", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("FinancialInstitutions", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "FinancialInstitutions", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("FinancialAccountIdentifiers", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "FinancialAccountIdentifiers", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("OrderingDetails", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "OrderingDetails", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("ProcessingStatus", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "ProcessingStatus", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("Tags", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "Tags", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("CustomFields", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "CustomFields", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = vendorsCmd_AutoResizeExcel("Companies", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "Companies", errAutoResize.Error())
	}

	// save document
	errDocumentsExcel := xlsx.SaveAs(homeDir + filename)
	if errDocumentsExcel != nil {
		log.Println(errDocumentsExcel)
	} else {
		log.Println("Excel document generated!!")
	}
}

func vendorsCmd_writeJson(response string, continueationToken string) {
	// create file
	var fileName string
	home, _ := homedir.Dir()
	if continueationToken == "" {
		fileName = config.EndpointMethod + "_" + "vendors" + ".json"
	} else {
		fileName = config.EndpointMethod + "_" + "vendors" + "_" + continueationToken + ".json"
	}

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

func vendorsCmd_CreateVendorsDocumentSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorCode")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorParent")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "buvid")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "name")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "sourceSystem")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "supplierAssignedAccountId")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "eligibleForSourcing")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorClass")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "LastUpdated")

	return xlsx, sheetIndex
}

func vendorsCmd_Basic_Authentication_GetVendorDocuments(continueationToken string, key string) ([]components.VendorEntity, string, error) {
	log.Println("Authentication method => BASIC")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Scope => %v", config.Scope)
	log.Printf("Endpoint method => %v", config.EndpointMethod)
	log.Printf("pageSize => %v", config.PageSize)

	if config.PageSize == 0 {
		viper.Set("api_pagesize", 500)
		_ = viper.WriteConfig()
	}

	var webMethod string
	webMethod = config.EndpointUrl + "v1/vendors/" + "?pageSize=" + strconv.Itoa(config.PageSize) + "&system=P2P"
	if len(key) > 0 {
		webMethod = config.EndpointUrl + "v1/vendors/" + key
	}

	req, err := http.NewRequest("GET", webMethod, nil)
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
		vendorsCmd_writeJson(string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object.
	tmpW := []components.VendorEntity{}
	if len(key) > 0 {
		_tmpW := components.VendorEntity{}
		errVendorResponse := json.Unmarshal(body, &_tmpW)
		if errVendorResponse != nil {
			if response.StatusCode == 404 {
				log.Printf("%s", "No data retrieved!")
			} else {
				log.Printf("Got error %s", errVendorResponse.Error())
			}
			return []components.VendorEntity{}, "", errVendorResponse
		} else {
			tmpW = append(tmpW, _tmpW)
		}

	} else {
		errVendorResponse := json.Unmarshal(body, &tmpW)
		if errVendorResponse != nil {
			if response.StatusCode == 404 {
				log.Printf("%s", "No data retrieved!")
			} else {
				log.Printf("Got error %s", errVendorResponse.Error())
			}
			return []components.VendorEntity{}, "", errVendorResponse
		}
	}

	return tmpW, continueationToken, nil
}

func vendorsCmd_AddRecord(json []components.VendorEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	var iRowIndex int
	var iRowIndex2 int = 2
	//var iRowIndex3 int = 1
	var iRowIndex2a int = 2
	var iRowIndex3a int = 2

	for _, rec := range json {
		iRowIndex = iRow

		// Add vendors
		xlsx, _ = vendorsCmd_AddDocumentRecord(rec, xlsx, iRowIndex)

		// Add Identifiers
		iRowIndex = iRow
		for _, identifier := range rec.Identifiers {
			xlsx, _ = vendorsCmd_AddIdentifierRecords(rec.ExternalCode, identifier, xlsx, iRowIndex)
			iRowIndex++
		}

		// Add Addressess
		iRowIndex = iRow
		for _, address := range rec.Addresses {
			xlsx, _ = vendorsCmd_AddAddressFieldRecords(rec.ExternalCode, address, xlsx, iRowIndex)

			// Add Additional Addresses
			for _, addressPart := range address.AdditionalAddressFields {
				xlsx, _ = vendorsCmd_AddAdditionalAddressFieldRecords(rec.ExternalCode, address.ExternalCode, addressPart, xlsx, iRowIndex2)
				iRowIndex2++
			}

			iRowIndex++
		}

		// Add Contacts
		iRowIndex = iRow
		for _, contact := range rec.Contacts {
			xlsx, _ = vendorsCmd_AddContactsRecords(rec.ExternalCode, contact, xlsx, iRowIndex)
			iRowIndex++
		}

		// Add DeliveryTerms
		if len(rec.DeliveryTerm.DeliveryTermCode) > 0 {
			iRowIndex = iRow
			xlsx, _ = vendorsCmd_AddDeliveryTermRecords(rec.ExternalCode, rec.DeliveryTerm, xlsx, iRowIndex)
		}

		// Add PaymentMeans
		iRowIndex = iRow
		for _, paymentMean := range rec.PaymentMeans {
			xlsx, _ = vendorsCmd_AddPaymentMeanRecords(rec.ExternalCode, paymentMean, xlsx, iRowIndex)

			// Add Financial Institutions
			for _, fa := range paymentMean.FinancialAccounts {
				xlsx, _ = vendorsCmd_AddFinancialInstitutionRecords(rec.ExternalCode, paymentMean.PaymentMeansCode, fa.FinancialInstitution, xlsx, iRowIndex2a)
				// Add Financial Accounting Identifiers
				for _, faIdentifier := range fa.FinancialAccountIdentifiers {
					xlsx, _ = vendorsCmd_AddFinancialAccountIdentifierRecords(rec.ExternalCode, paymentMean.PaymentMeansCode, faIdentifier, xlsx, iRowIndex3a)
					iRowIndex3a++
				}
				iRowIndex2a++
			}
			iRowIndex++
		}

		// Add OrderingDetail
		iRowIndex = iRow
		xlsx, _ = vendorsCmd_AddOrderingDetailRecords(rec.ExternalCode, rec.OrderingDetails, xlsx, iRowIndex)

		// Add ProcessingStatus
		iRowIndex = iRow
		xlsx, _ = vendorsCmd_AddProcessingStatusRecords(rec.ExternalCode, rec.ProcessingStatus, xlsx, iRowIndex)

		// Add Tags
		iRowIndex = iRow
		for _, tag := range rec.Tags {
			xlsx, _ = vendorsCmd_AddTagRecords(rec.ExternalCode, tag, xlsx, iRowIndex)
			iRowIndex++
		}

		// Add CustomFields
		iRowIndex = iRow
		for _, customField := range rec.CustomFields {
			xlsx, _ = vendorsCmd_AddCustomFieldRecords(rec.ExternalCode, customField, xlsx, iRowIndex)
			iRowIndex++
		}

		//* Add Companies
		iRowIndex = iRow
		for _, company := range rec.Companies {
			xlsx, _ = vendorsCmd_AddCompanyRecord(rec.ExternalCode, company, xlsx, iRowIndex)
			iRowIndex++
		}
		iRow++
	}
	return xlsx, iRow
}

func vendorsCmd_AutoResizeExcel(sheetName string, xlsx *excelize.File) (*excelize.File, error) {
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

func vendorsCmd_AddDocumentRecord(vendorItem components.VendorEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "Vendors"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.ExternalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.VendorCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.VendorParent)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.Buvid)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.Name)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.SourceSystem)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.Description)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.SupplierAssignedAccountId)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.EligibleForSourcing)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.VendorClass)

	xlsx, _ = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, vendorItem.LastUpdated)

	return xlsx, iRow
}

func vendorsCmd_AddCol(xlsx *excelize.File, sheetName string, colIndex int, rowIndex int, caption interface{}) (*excelize.File, int) {

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

func vendorsCmd_AddCompanyRecord(externalCode string, companyItem components.VendorCompanyEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "Companies"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, companyItem.CompanyCode)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, companyItem.InheritToChildUnits)

	return xlsx, iRow
}

func vendorsCmd_CreateCompaniesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from list record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "companyCode")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "inheritToChildUnits")

	return xlsx, sheetIndex
}

func vendorsCmd_OAuth02_AuthenticationToken() {
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

func vendorsCmd_OAuth02_GetVendorDocuments(continueationToken string, key string) ([]components.VendorEntity, string, error) {

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var webMethod string
	webMethod = config.EndpointUrl + "v1/vendors/" + "?pageSize=" + strconv.Itoa(config.PageSize) + "&system=P2P"
	if len(key) > 0 {
		webMethod = config.EndpointUrl + "v1/vendors/" + key
	}

	// get account
	req, err := http.NewRequest("GET", webMethod, nil)
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
		log.Println("getVendorDocuments - response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		vendorsCmd_writeJson(string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	tmpW := []components.VendorEntity{}
	if len(key) > 0 {
		_tmpW := components.VendorEntity{}
		errVendorResponse := json.Unmarshal(body, &_tmpW)
		if errVendorResponse != nil {
			if response.StatusCode == 404 {
				log.Printf("%s", "No data retrieved!")
			} else {
				log.Printf("Got error %s", errVendorResponse.Error())
			}
			return []components.VendorEntity{}, "", errVendorResponse
		} else {
			tmpW = append(tmpW, _tmpW)
		}

	} else {
		errVendorResponse := json.Unmarshal(body, &tmpW)
		if errVendorResponse != nil {
			if response.StatusCode == 404 {
				log.Printf("%s", "No data retrieved!")
			} else {
				log.Printf("Got error %s", errVendorResponse.Error())
			}
			return []components.VendorEntity{}, "", errVendorResponse
		}
	}

	return tmpW, continueationToken, nil
}

func vendorsCmd_CreateIdentifiersSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "id")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "schemeId")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "defaultPartyId")

	return xlsx, sheetIndex
}

func vendorsCmd_AddIdentifierRecords(externalCode string, item components.VendorIdentifierEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "Identifiers"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Id)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.SchemeId)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.DefaultPartyId)

	return xlsx, iRow
}

func vendorsCmd_CreateAddressesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")       // externalCode from address record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "name")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "addressType")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "addressLine1")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "addressLine2")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "addressLine3")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "cityName")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "postalZone")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "poBox")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "streetName")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "locality")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "countrySubEntity")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "countrySubEntityDescription")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "countryId")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "default")

	return xlsx, sheetIndex
}

func vendorsCmd_AddAddressFieldRecords(externalCode string, item components.VendorAddressEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "Addresses"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.ExternalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Name)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Description)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AddressType)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AddressLine1)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AddressLine2)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AddressLine3)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CityName)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.PostalZone)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.PoBox)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.StreetName)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Locality)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CountrySubEntity)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CountrySubEntityDescription)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CountryId)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Default)

	return xlsx, iRow
}

func vendorsCmd_CreateAdditionalAddressFieldsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode")  // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "addressExternalCode") // externalCode from address record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "key")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "value")

	return xlsx, sheetIndex
}

func vendorsCmd_AddAdditionalAddressFieldRecords(externalCode string, addressPartCode string, item components.VendorAddressPartEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "AdditionalAddressFields"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, addressPartCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Key)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Value)

	return xlsx, iRow
}

func vendorsCmd_CreateContactsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "name")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "telephone")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "telefax")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "email")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "website")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "role")

	return xlsx, sheetIndex
}

func vendorsCmd_AddContactsRecords(externalCode string, item components.VendorContactEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "Contacts"

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode) // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Name)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Description)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Telephone)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Telefax)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Email)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Website)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Role)

	return xlsx, iRow
}

func vendorsCmd_CreateDeliveryTermsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "deliveryTermCode")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "deliveryLocation")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")

	return xlsx, sheetIndex
}

func vendorsCmd_AddDeliveryTermRecords(externalCode string, item components.VendorDeliveryTermEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "DeliveryTerms"

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode) // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.DeliveryTermCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.DeliveryLocation)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Description)

	return xlsx, iRow
}

func vendorsCmd_CreatePaymentMeansSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentMeansCode")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCode")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "default")

	return xlsx, sheetIndex
}

func vendorsCmd_AddPaymentMeanRecords(externalCode string, item components.VendorPaymentMeanEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "PaymentMeans"

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode) // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.PaymentMeansCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Description)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CurrencyCode)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Default)

	return xlsx, iRow
}

func vendorsCmd_CreateFinancialInstitutionSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentMeansCode")   // reference back to payment means
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "name")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "id")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "schemeId")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "branchId")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "branchIdSchemeId")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "countryId")

	return xlsx, sheetIndex
}

func vendorsCmd_AddFinancialInstitutionRecords(externalCode string, paymentmeansCode string, item components.FinancialInstitutionEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "FinancialInstitutions"

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode) // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentmeansCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Name)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Id)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.SchemeId)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.BranchId)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.BranchIdSchemeId)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CountryId)

	return xlsx, iRow
}

func vendorsCmd_CreateFinancialAccountIdentifiersSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentMeansCode")   // reference back to payment means
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "schemeId")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "description")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "id")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountHolderName")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "currencyCode")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountAdditionalData1")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "accountAdditionalData2")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "default")

	return xlsx, sheetIndex
}

func vendorsCmd_AddFinancialAccountIdentifierRecords(externalCode string, paymentmeansCode string, item components.FinancialAccountIdentifierEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "FinancialAccountIdentifiers"

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode) // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, paymentmeansCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.SchemeId)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Description)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Id)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AccountHolderName)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CurrencyCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AccountAdditionalData1)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AccountAdditionalData2)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Default)

	return xlsx, iRow
}

func vendorsCmd_CreateOrderingDetailsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderingFormat")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderingMessageLanguage")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderingLanguage")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderEmail")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "orderProcessType")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "minimumOrderAllowed")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "isTaxable")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "automaticallyReceiveOnInvoice")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "automaticallyReceiveOnOrder")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "hasActiveCatalog")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "createOrderAutomatically")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "deliverOrderAutomatically")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "noFreeformItems")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "sendToNetwork")

	return xlsx, sheetIndex
}

func vendorsCmd_AddOrderingDetailRecords(externalCode string, item components.VendorOrderingDetailsEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {

	sheetName := "OrderingDetails"

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode) // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.OrderingFormat)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.OrderingMessageLanguage)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.OrderingLanguage)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.OrderEmail)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.OrderProcessType)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.MinimumOrderAllowed)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.IsTaxable)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AutomaticallyReceiveOnInvoice)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.AutomaticallyReceiveOnOrder)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.HasActiveCatalog)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.CreateOrderAutomatically)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.DeliverOrderAutomatically)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.NoFreeformItems)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.SendToNetwork)

	return xlsx, iRow
}

func vendorsCmd_CreateProcessingStatusSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "active")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "paymentDenied")

	return xlsx, sheetIndex
}

func vendorsCmd_AddProcessingStatusRecords(externalCode string, item components.VendorProcessingStatusEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "ProcessingStatus"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Active)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.PaymentDenied)

	return xlsx, iRow
}

func vendorsCmd_CreateTagsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "name")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "value")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "tagGroup")

	return xlsx, sheetIndex
}

func vendorsCmd_AddTagRecords(externalCode string, item components.VendorTagEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "Tags"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Name)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.TagGroup)

	return xlsx, iRow
}

func vendorsCmd_CreateCustomFieldsSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "vendorExternalCode") // externalCode from vendor record
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "name")
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "value")
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, 1, "groupName")

	return xlsx, sheetIndex
}

func vendorsCmd_AddCustomFieldRecords(externalCode string, item components.VendorAdditionalFieldEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "CustomFields"
	colindex := 1

	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.Name)
	xlsx, _ = vendorsCmd_AddCol(xlsx, sheetName, colindex, iRow, item.GroupName)

	return xlsx, iRow
}
