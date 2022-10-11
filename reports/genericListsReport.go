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

// genericListsCmd represents the genericLists command
var genericListsCmd = &cobra.Command{
	Use:   "genericLists",
	Short: "Retrieves all data from the lists API and creates Excel document of the content",
	Long: `Retrieves all data from the lists API. 
	       Data is written into an Excel document under the .multiTool folder`,
	Run: func(cmd *cobra.Command, args []string) {
		config.EndpointMethod = "lists"
		config.Scope = "lists.read"

		// setup list of genericLists
		listArray := strings.Split(config.GenericLists, ",")
		if len(listArray) > 0 {
			switch config.AuthMethod {
			case BASIC:
				getGenericListBasicData(listArray)
			case OAUTH2:
				getGenericListSecuredData(listArray)
			}
		} else {
			log.Println("Missing values from parameter api_lists")
		}
	},
}

func init() {
	reportCmd.AddCommand(genericListsCmd)
}

func getGenericListBasicData(listArray []string) {

	// var listArray []string
	var continueationToken2BasicAuth string
	var iRowBasicAuth int
	var xlsxBasicAuth *excelize.File
	var indexBasicAuth int

	for _, element := range listArray {

		iRowBasicAuth = 0
		listname := element
		continueationToken2BasicAuth = ""

		// fetch
		iRowBasicAuth = 2
		for {
			jsonBasicAuth, continueationTokenBasicAuth, errgenericListsCmdBasicAuthGetGenericListDocuments := genericListsCmd_Basic_Authentication_GetGenericListDocuments(listname, continueationToken2BasicAuth)

			if errgenericListsCmdBasicAuthGetGenericListDocuments == nil {
				// create Excel Document (report)
				if iRowBasicAuth == 2 {
					xlsxBasicAuth, indexBasicAuth = genericListsCmd_CreateExcelDocument()
				}

				//create rows
				xlsxBasicAuth, iRowBasicAuth = genericListsCmd_AddRecord(jsonBasicAuth, xlsxBasicAuth, iRowBasicAuth)

				if continueationTokenBasicAuth == "" {
					genericListsCmd_SaveExcelDocument(listname, xlsxBasicAuth, indexBasicAuth, "")
					break

				} else {

					continueationToken2BasicAuth = continueationTokenBasicAuth
					log.Printf("continueationToken => %s ", continueationToken2BasicAuth)
					time.Sleep(10 * time.Millisecond)
					log.Println("fetch next ...")

					// for debug purpose we break
					//genericListsCmd_SaveExcelDocument(listname, xlsxBasicAuth, indexBasicAuth, "")
					//break
				}
			} else {
				break
			}
		}
	}
}

func getGenericListSecuredData(listArray []string) {
	genericListsCmd_OAuth02_AuthenticationToken()
	if len(config.Token) > 0 {
		// var listArray []string
		var continueationToken2OAuth02 string
		var iRowBasicAuth int
		var xlsxBasicAuth *excelize.File
		var indexBasicAuth int

		for _, element := range listArray {

			iRowBasicAuth = 0
			listname := element
			continueationToken2OAuth02 = ""

			// fetch
			iRowBasicAuth = 2
			for {
				jsonBasicAuth, continueationTokenOAuth02, errgenericListsCmdOAuth02GetGenericListDocuments := genericListsCmd_OAuth02_GetGenericListDocuments(listname, continueationToken2OAuth02)

				if errgenericListsCmdOAuth02GetGenericListDocuments == nil {
					// create Excel Document (report)
					if iRowBasicAuth == 2 {
						xlsxBasicAuth, indexBasicAuth = genericListsCmd_CreateExcelDocument()
					}

					//create rows
					xlsxBasicAuth, iRowBasicAuth = genericListsCmd_AddRecord(jsonBasicAuth, xlsxBasicAuth, iRowBasicAuth)

					if continueationTokenOAuth02 == "" {
						genericListsCmd_SaveExcelDocument(listname, xlsxBasicAuth, indexBasicAuth, "")
						break

					} else {

						continueationToken2OAuth02 = continueationTokenOAuth02
						log.Printf("continueationToken => %s ", continueationToken2OAuth02)
						time.Sleep(10 * time.Millisecond)
						log.Println("fetch next ...")

						// for debug purpose we break
						//genericListsCmd_SaveExcelDocument(listname, xlsxBasicAuth, indexBasicAuth, "")
						//break
					}
				} else {
					break
				}
			}
		}
	}
}

func genericListsCmd_CreateExcelDocument() (*excelize.File, int) {

	var index int
	var sheetName string

	// create Excel Document (report)
	xlsx := excelize.NewFile()

	// create Sheet
	sheetName = "list"
	xlsx, _ = genericListsCmd_CreateGenericListDocumentSheet(xlsx, sheetName)

	sheetName = "companies"
	genericListsCmd_CreateCompaniesSheet(xlsx, sheetName)

	return xlsx, index
}

func genericListsCmd_SaveExcelDocument(listName string, xlsx *excelize.File, index int, continueationToken string) {

	home, _ := homedir.Dir()
	homeDir := home + "/.multiTool/"
	var filename string

	if continueationToken == "" {
		filename = listName + ".xlsx"
	} else {
		filename = listName + "_" + continueationToken + ".xlsx"
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
	xlsx, errAutoResize = genericListsCmd_AutoResizeExcel("list", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "list", errAutoResize.Error())
	}

	// resize companies sheet
	xlsx, errAutoResize = genericListsCmd_AutoResizeExcel("companies", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "companies", errAutoResize.Error())
	}

	// save document
	errDocumentsExcel := xlsx.SaveAs(homeDir + filename)
	if errDocumentsExcel != nil {
		log.Println(errDocumentsExcel)
	} else {
		log.Println("Excel document generated!!")
	}
}

func genericListsCmd_Basic_Authentication_GetGenericListDocuments(listName string, continueationToken string) (components.GenericListResponse, string, error) {
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

	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/lists/"+listName+"?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
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
		genericListsCmd_writeJson(listName, string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	tmpW := components.GenericListResponse{}
	errGenericListResponse := json.Unmarshal(body, &tmpW)
	if errGenericListResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Printf("Got error %s", errGenericListResponse.Error())
		}
		return components.GenericListResponse{}, "", errGenericListResponse
	}

	return tmpW, continueationToken, nil
}

func genericListsCmd_writeJson(listName string, response string, continueationToken string) {
	// create file
	var fileName string
	home, _ := homedir.Dir()
	if continueationToken == "" {
		fileName = config.EndpointMethod + "_" + listName + ".json"
	} else {
		fileName = config.EndpointMethod + "_" + listName + "_" + continueationToken + ".json"
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

func genericListsCmd_CreateGenericListDocumentSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_1")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_2")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_3")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_4")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_5")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_6")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_7")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_8")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_9")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_10")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_11")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_12")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_13")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_14")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_15")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_16")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_17")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_18")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_19")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_20")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_21")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_22")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_23")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_24")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "text_25")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_1")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_2")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_3")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_4")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_5")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_6")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_7")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_8")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_9")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_10")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_11")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_12")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_13")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_14")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "number_15")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_1")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_2")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_3")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_4")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_5")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_6")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_7")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_8")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_9")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_10")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_11")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_12")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_13")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_14")
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "dateValue_15")

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "lastUpdated")
	xlsx, _ = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "listContent")

	return xlsx, sheetIndex
}

func genericListsCmd_CreateCompaniesSheet(xlsx *excelize.File, sheetName string) (*excelize.File, int) {

	// create excelsheet
	sheetIndex := xlsx.NewSheet(sheetName)

	// create header fields
	colindex := 1

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "externalCode") // externalCode from list record
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "companyCode")
	xlsx, _ = genericListsCmd_AddCol(xlsx, sheetName, colindex, 1, "active")

	return xlsx, sheetIndex
}

func genericListsCmd_AddCol(xlsx *excelize.File, sheetName string, colIndex int, rowIndex int, caption interface{}) (*excelize.File, int) {

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

func genericListsCmd_AutoResizeExcel(sheetName string, xlsx *excelize.File) (*excelize.File, error) {
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

func genericListsCmd_AddRecord(json components.GenericListResponse, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	var iRowIndex int
	for _, rec := range json.ListItems {
		iRowIndex = iRow
		xlsx, _ = genericListsCmd_AddDocumentRecord(rec, xlsx, iRowIndex)

		//* Processing Order numbers
		iRowIndex = iRow
		for _, company := range rec.Companies {
			xlsx, _ = genericListsCmd_AddCompanyRecord(rec.ExternalCode, company, xlsx, iRowIndex)
			iRowIndex++
		}
		iRow++
	}
	return xlsx, iRow
}

func genericListsCmd_AddDocumentRecord(listItem components.GenericListEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "list"
	colindex := 1

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.ExternalCode)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_1)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_2)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_3)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_4)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_5)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_6)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_7)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_8)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_9)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_10)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_11)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_12)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_13)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_14)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_15)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_16)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_17)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_18)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_19)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_20)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_21)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_22)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_23)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_24)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Text_25)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_1)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_2)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_3)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_4)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_5)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_6)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_7)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_8)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_9)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_10)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_11)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_12)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_13)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_14)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.Number_15)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_1)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_2)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_3)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_4)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_5)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_6)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_7)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_8)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_9)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_10)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_11)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_12)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_13)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_14)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.DateValue_15)

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.LastUpdated)

	xlsx, _ = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, listItem.ListContent)

	return xlsx, iRow
}

func genericListsCmd_AddCompanyRecord(externalCode string, companyItem components.GenericListCompanyEntity, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	sheetName := "companies"
	colindex := 1

	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, externalCode)
	xlsx, colindex = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, companyItem.CompanyCode)

	xlsx, _ = genericListsCmd_AddCol(xlsx, sheetName, colindex, iRow, companyItem.Active)

	return xlsx, iRow
}

func genericListsCmd_OAuth02_AuthenticationToken() {
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

func genericListsCmd_OAuth02_GetGenericListDocuments(listName string, continueationToken string) (components.GenericListResponse, string, error) {

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// get account
	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/lists/"+listName+"?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
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
		log.Println("getGenericListDocuments - response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		genericListsCmd_writeJson(listName, string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	tmpW := components.GenericListResponse{}
	errGenericListResponse := json.Unmarshal(body, &tmpW)
	if errGenericListResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Printf("Got error %s", errGenericListResponse.Error())
		}
		return components.GenericListResponse{}, "", errGenericListResponse
	}

	return tmpW, continueationToken, nil
}
