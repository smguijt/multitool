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

// requestStatusCmd represents the requestStatus command
var requestStatusCmd = &cobra.Command{
	Use:   "requestStatus",
	Short: "Retrieves all data from the requestStatus API and creates Excel document of the content",
	Long: `Retrieves all data from the requestStatus API. 
	Data is written into an Excel document under the .multiTool folder`,
	Run: func(cmd *cobra.Command, args []string) {
		// set api_endpointmethod
		config.EndpointMethod = "requestStatus"
		config.Scope = "requestStatus.read"

		switch config.AuthMethod {
		case BASIC:
			var continueationToken2BasicAuth string
			iRowBasicAuth := 0
			var xlsxBasicAuth *excelize.File
			var indexBasicAuth int

			// fetch
			iRowBasicAuth = 2
			for {
				jsonBasicAuth, continueationTokenBasicAuth, errRequestStatusCmdBasicAuthGetRequestStatus := requestStatusCmd_Basic_Authentication_GetRequestStatus(continueationToken2BasicAuth)
				if errRequestStatusCmdBasicAuthGetRequestStatus == nil {
					// create Excel Document (report)
					if iRowBasicAuth == 2 {
						xlsxBasicAuth, indexBasicAuth = requestStatusCmd_CreateExcelDocument()
					}

					//create rows
					xlsxBasicAuth, iRowBasicAuth = requestStatusCmd_AddRecord(jsonBasicAuth, xlsxBasicAuth, iRowBasicAuth)

					if continueationTokenBasicAuth == "" {
						// set active sheet and save document
						requestStatusCmd_SaveExcelDocument(xlsxBasicAuth, indexBasicAuth, "")

						// end loop
						break

					} else {
						continueationToken2BasicAuth = continueationTokenBasicAuth
						log.Printf("continueationToken => %s ", continueationToken2BasicAuth)
						log.Println("fetch next ...")
					}
				} else {
					log.Fatalf(errRequestStatusCmdBasicAuthGetRequestStatus.Error())
				}
			}

		case OAUTH2:
			requestStatusCmd_OAuth02_Authentication()
			if len(config.Token) > 0 {
				// init
				var continueationToken2 string
				iRow := 0
				var xlsx *excelize.File
				var index int

				// fetch
				iRow = 2
				for {
					json, continueationToken, errRequestStatusCmdOAuth02GetRequestStatus := requestStatusCmd_OAuth02_GetRequestStatus(continueationToken2)
					if errRequestStatusCmdOAuth02GetRequestStatus == nil {
						// create Excel Document (report)
						if iRow == 2 {
							xlsx, index = requestStatusCmd_CreateExcelDocument()
						}

						//create rows
						xlsx, iRow = requestStatusCmd_AddRecord(json, xlsx, iRow)

						if continueationToken == "" {
							// set active sheet and save document
							requestStatusCmd_SaveExcelDocument(xlsx, index, "")

							// end loop
							break

						} else {
							continueationToken2 = continueationToken
							log.Printf("continueationToken => %s ", continueationToken2)
							log.Println("fetch next ...")

							//if config.Debug {
							//	SaveExcelDocument(xlsx, index, continueationToken2)
							//}
							//log.Println()
						}
					}
				}
			}
		}
	},
}

func init() {
	reportCmd.AddCommand(requestStatusCmd)
}

func requestStatusCmd_Basic_Authentication_GetRequestStatus(continueationToken string) (components.RequestStatusResponse, string, error) {
	log.Println("Authentication method => BASIC")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Scope => %v", config.Scope)
	log.Printf("Endpoint method => %v", config.EndpointMethod)
	log.Printf("pageSize => %v", config.PageSize)
	log.Printf("entityType => %v", config.EntityType)

	//entityType=TaxCode
	tmpEntityType := config.EntityType
	if len(tmpEntityType) > 0 {
		tmpEntityType = "&entityType=" + tmpEntityType
	}

	tmpUrl := config.EndpointUrl + "v1/requestStatus?pageSize=" + strconv.Itoa(config.PageSize) + "&system=P2P" + tmpEntityType
	log.Printf("url to execute: => %v", tmpUrl)

	req, err := http.NewRequest("GET", tmpUrl, nil)
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
		requestStatusCmd_writeJson(string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")

	// unmarshal response to object
	tmpW := components.RequestStatusResponse{}
	errRequestStatusResponse := json.Unmarshal(body, &tmpW)
	if errRequestStatusResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errRequestStatusResponse.Error())
		}
		return components.RequestStatusResponse{}, "", errRequestStatusResponse
	}

	return tmpW, continueationToken, nil
}

func requestStatusCmd_OAuth02_Authentication() {
	log.Println("Authentication method => OAUTH2")

	// setup client
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

func requestStatusCmd_OAuth02_GetRequestStatus(continueationToken string) (components.RequestStatusResponse, string, error) {

	// config.Debug = true

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// get account
	req, err := http.NewRequest("GET", config.EndpointUrl+"v1/requestStatus?pageSize="+strconv.Itoa(config.PageSize)+"&system=P2P", nil)
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
		log.Println("getRequestStatus - response status: ", response.StatusCode)
	}
	defer response.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(response.Body)
	if config.Debug {
		requestStatusCmd_writeJson(string([]byte(body)), continueationToken)
	}

	// get header
	header := response.Header
	continueationToken = header.Get("X-Amz-Meta-Continuationtoken")
	// log.Printf("CT: %v", continueationToken)

	// unmarshal response to object
	tmpW := components.RequestStatusResponse{}
	errRequestStatusResponse := json.Unmarshal(body, &tmpW)
	if errRequestStatusResponse != nil {
		if response.StatusCode == 404 {
			log.Printf("%s", "No data retrieved!")
		} else {
			log.Fatalf("Got error %s", errRequestStatusResponse.Error())
		}
		return components.RequestStatusResponse{}, "", errRequestStatusResponse
	}

	return tmpW, continueationToken, nil
}

func requestStatusCmd_writeJson(response string, continueationToken string) {
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

func requestStatusCmd_CreateExcelDocument() (*excelize.File, int) {
	// create Excel Document (report)
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("requestStatus")

	// create main
	xlsx.SetCellValue("requestStatus", "A1", "RequestId")
	xlsx.SetCellValue("requestStatus", "B1", "EntityType")
	xlsx.SetCellValue("requestStatus", "C1", "RequestStatus")
	xlsx.SetCellValue("requestStatus", "D1", "LastUpdated")
	// create system
	xlsx.SetCellValue("requestStatus", "E1", "SystemStatus")
	xlsx.SetCellValue("requestStatus", "F1", "SystemLastUpdated")
	// create items
	xlsx.SetCellValue("requestStatus", "G1", "Status")
	xlsx.SetCellValue("requestStatus", "H1", "ExternalCode")
	// create errors
	xlsx.SetCellValue("requestStatus", "I1", "Type")
	xlsx.SetCellValue("requestStatus", "J1", "Category")
	xlsx.SetCellValue("requestStatus", "K1", "Message")

	return xlsx, index

}

func requestStatusCmd_AutoResizeExcel(sheetName string, xlsx *excelize.File) (*excelize.File, error) {
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

func requestStatusCmd_SaveExcelDocument(xlsx *excelize.File, index int, continueationToken string) {

	home, _ := homedir.Dir()
	homeDir := home + "/.multiTool/"
	var filename string

	if continueationToken == "" {
		filename = "requestStatus.xlsx"
	} else {
		filename = "requestStatus_" + continueationToken + ".xlsx"
	}
	prefix := config.Prefix
	if len(prefix) > 0 {
		filename = prefix + "_" + filename
	}

	// set active sheet
	xlsx.SetActiveSheet(index)
	xlsx.DeleteSheet("Sheet1")

	// resize requestStatus
	var errAutoResize error
	xlsx, errAutoResize = requestStatusCmd_AutoResizeExcel("requestStatus", xlsx)
	if errAutoResize != nil {
		log.Fatalf("Failed to autoResize sheet %s. Reason: %s", "requestStatus", errAutoResize.Error())
	}

	// save document
	errRequestStatusExcel := xlsx.SaveAs(homeDir + filename)
	if errRequestStatusExcel != nil {
		log.Println(errRequestStatusExcel)
	} else {
		log.Println("Excel document generated!!")
	}
}

func requestStatusCmd_AddRecord(json components.RequestStatusResponse, xlsx *excelize.File, iRow int) (*excelize.File, int) {
	for _, rec := range json.RequestStatus {
		for _, rec2 := range rec.Systems {
			for _, rec3 := range rec2.Items {
				if rec3.Status != "Success" {
					for _, rec4 := range rec3.Errors {
						// Main
						xlsx.SetCellValue("requestStatus", "A"+strconv.Itoa(iRow), rec.RequestId)
						xlsx.SetCellValue("requestStatus", "B"+strconv.Itoa(iRow), rec.EntityType)
						xlsx.SetCellValue("requestStatus", "C"+strconv.Itoa(iRow), rec.RequestStatus)
						xlsx.SetCellValue("requestStatus", "D"+strconv.Itoa(iRow), rec.LastUpdated)
						// Systems
						xlsx.SetCellValue("requestStatus", "E"+strconv.Itoa(iRow), rec2.SystemStatus)
						xlsx.SetCellValue("requestStatus", "F"+strconv.Itoa(iRow), rec2.LastUpdated)
						// Items
						xlsx.SetCellValue("requestStatus", "G"+strconv.Itoa(iRow), rec3.Status)
						xlsx.SetCellValue("requestStatus", "H"+strconv.Itoa(iRow), rec3.ExternalCode)
						// Errors
						xlsx.SetCellValue("requestStatus", "I"+strconv.Itoa(iRow), rec4.Type)
						xlsx.SetCellValue("requestStatus", "J"+strconv.Itoa(iRow), rec4.Category)
						xlsx.SetCellValue("requestStatus", "K"+strconv.Itoa(iRow), rec4.Message)
						iRow++
					}
				} else {
					// Main
					xlsx.SetCellValue("requestStatus", "A"+strconv.Itoa(iRow), rec.RequestId)
					xlsx.SetCellValue("requestStatus", "B"+strconv.Itoa(iRow), rec.EntityType)
					xlsx.SetCellValue("requestStatus", "C"+strconv.Itoa(iRow), rec.RequestStatus)
					xlsx.SetCellValue("requestStatus", "D"+strconv.Itoa(iRow), rec.LastUpdated)
					// Systems
					xlsx.SetCellValue("requestStatus", "E"+strconv.Itoa(iRow), rec2.SystemStatus)
					xlsx.SetCellValue("requestStatus", "F"+strconv.Itoa(iRow), rec2.LastUpdated)
					// Items
					xlsx.SetCellValue("requestStatus", "G"+strconv.Itoa(iRow), rec3.Status)
					xlsx.SetCellValue("requestStatus", "H"+strconv.Itoa(iRow), rec3.ExternalCode)
					iRow++
				}
			}
		}
	}

	return xlsx, iRow
}
