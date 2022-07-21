/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"bytes"
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

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reprocessCmd represents the reprocess command
// https://baswarecorp.sharepoint.com/sites/BWCN/BCN/P2P/ConfigurationManual/Current/Basware%20API%20-%20How%20to%20delete%20and%20redistribute%20data.aspx
var reprocessCmd = &cobra.Command{
	Use:   "reprocess",
	Short: "resends a given API post request from the API storage to the available subscribed services",
	Long: `resends a given API post request from the API storage to the available subscribed services.
	
	- multiTool.exe --config <config>.yaml reprocess --entityType Vendors --externalCode "code1,code2..."
	- multiTool.exe --config <config>.yaml reprocess --entityType GenericList --externalCode "code1,code2..." --listKey "INV_LIST_1"
	- multiTool.exe --config <config>.yaml reprocess --entityType Project --externalCode "code1,code2..." --version "V2"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check entityType
		_entityType, errFlagsEntityType := cmd.Flags().GetString("entityType")
		if errFlagsEntityType != nil {
			log.Fatalf("%s", errFlagsEntityType.Error())
		} else if len(_entityType) == 0 {
			log.Fatalf("%s", "required flag \"entityType\" not set")
		} else {
			log.Printf("entityType: %s", _entityType)
		}

		// Check externalCode
		_externalCodes, errFlagsExternalCode := cmd.Flags().GetStringSlice("externalCode")
		if errFlagsExternalCode != nil {
			log.Fatalf("%s", errFlagsExternalCode.Error())
		} else if len(_externalCodes) == 0 {
			log.Fatalf("%s", "required flag \"externalCode\" not set")
		}

		// Check listKey
		var _listKey string
		var errFlagsListKey error
		if _entityType == "GenericList" {
			_listKey, errFlagsListKey = cmd.Flags().GetString("listKey")
			if errFlagsListKey != nil {
				log.Fatalf("%s", errFlagsListKey.Error())
			} else if len(_listKey) == 0 {
				log.Fatalf("%s", "required flag \"listKey\" not set")
			}
		}

		// Check version
		var _version string
		_version, _ = cmd.Flags().GetString("version")
		log.Printf("version: %s", _version)
		if _entityType == "Project" {
			if _version != "V2" {
				log.Fatalf("%s", "required flag \"version\" must be set to V2")
			}
		}

		// Check subscribed Service
		var _subscribedService string
		_subscribedService, _ = cmd.Flags().GetString("subscribedService")
		log.Printf("subscribedService: %s", _subscribedService)
		if _subscribedService != "All" {
			log.Fatalf("%s", "flag \"subscribedService\" only supports the value \"All\" at the moment")
		}

		// Process
		config.EndpointMethod = "redistribute"
		config.Scope = "redistribute.write"

		//for _, key := range _externalCodes {
		//log.Printf("processing key): %s\n", key)

		switch config.AuthMethod {
		case BASIC:
			reprocessCmd_postRedistribute(true, _entityType, _externalCodes, _listKey, _version, _subscribedService)
		case OAUTH2:
			reprocessCmd_OAuth02_AuthenticationToken()
			reprocessCmd_postRedistribute(false, _entityType, _externalCodes, _listKey, _version, _subscribedService)
		}
		//}

	},
}

func init() {
	rootCmd.AddCommand(reprocessCmd)

	// set PersistentFlag for entityType and mark it as mandatory!
	reprocessCmd.PersistentFlags().String("entityType", "", "")
	reprocessCmd.MarkPersistentFlagRequired("entityType")

	// set PersistentFlag for exteralCode and mark it as mandatory!
	reprocessCmd.PersistentFlags().StringSlice("externalCode", nil, "")
	reprocessCmd.MarkPersistentFlagRequired("externalCode")

	reprocessCmd.PersistentFlags().String("listKey", "", "")
	reprocessCmd.PersistentFlags().String("version", "V1", "")
	reprocessCmd.PersistentFlags().String("subscribedService", "All", "")
}

func reprocessCmd_postRedistribute(useBasicAuthentication bool, entityType string, externalCode []string, listKey string, version string, subscribedService string) {

	// activate debug so that output json is written
	config.Debug = true

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// create request object
	var payload components.RedistributeRequest
	payload.SubscribedService = subscribedService
	payload.Version = version
	payload.EntityType = entityType
	payload.ExternalCodes = externalCode

	if entityType == "GenericList" {
		payload.ListKey = listKey
	}

	// convert struct to json
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	// create new http request message
	req, err := http.NewRequest("POST", config.EndpointUrl+"v1/redistribute", payloadBuf)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}

	// add content-Type
	req.Header.Set("Content-Type", "application/json")

	// add authentication header
	if useBasicAuthentication {
		req.SetBasicAuth(config.Username, config.Password)
	} else {
		bearer := "Bearer " + config.Token
		req.Header.Add("Authorization", bearer)
	}

	// send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	} else {
		log.Println("response status: ", resp.StatusCode)
		//if config.Debug {
		//log.Printf("response header: %s\n", resp.Header)
		//}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if config.Debug {
		//log.Printf("%s\n", string(body))
		reprocessCmd_writeJson(string(body), "", entityType, listKey)
	}
}

func reprocessCmd_OAuth02_AuthenticationToken() {
	log.Println("Authentication method => OAUTH2")

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Scope => %v", config.Scope)
	log.Printf("Endpoint method => %v", config.EndpointMethod)

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

func reprocessCmd_writeJson(response string, continueationToken string, entityType string, listKey string) {
	// create file
	var fileName string
	home, _ := homedir.Dir()
	if continueationToken == "" {
		fileName = config.EndpointMethod + "_" + entityType
		if len(listKey) > 0 {
			fileName += "_" + listKey
		}
		fileName += ".json"
	} else {
		fileName = config.EndpointMethod + "_" + entityType
		if len(listKey) > 0 {
			fileName += "_" + listKey
		}
		fileName += "_" + continueationToken + ".json"
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
