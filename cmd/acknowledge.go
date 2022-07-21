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

// acknowledgeCmd represents the acknowledge command
var acknowledgeCmd = &cobra.Command{
	Use:   "acknowledge",
	Short: "Send acknowledge request to remove item from the given API result list",
	Long: `Send acknowledge request to remove item from the given API result list:

    --multiTool --config [config].yaml acknowledge --type requestStatus --id "2d2127dd-ed88-44d9-a6ad-9c88e0eef394"
    --multiTool --config [config].yaml acknowledge --type requestStatus --id "2d2127dd-ed88-44d9-a6ad-9c88e0eef394", "3d2127dd-ed88-44d9-a6ad-9c88e0eef395"
	--multiTool --config [config].yaml acknowledge --type accountingDocument --id "2d2127dd94"`,
	Run: func(cmd *cobra.Command, args []string) {

		// Check entityType
		_entityType, errFlagsEntityType := cmd.Flags().GetString("type")
		if errFlagsEntityType != nil {
			log.Fatalf("%s", errFlagsEntityType.Error())
		} else if len(_entityType) == 0 {
			log.Fatalf("%s", "required flag \"type\" not set")
		} else {
			log.Printf("type: %s", _entityType)
		}

		// Check id
		_id, errFlagsId := cmd.Flags().GetStringSlice("id")
		if errFlagsId != nil {
			log.Fatalf("%s", errFlagsId.Error())
		} else if len(_id) == 0 {
			log.Fatalf("%s", "required flag \"_id\" not set")
		} else {
			log.Printf("_id: %s", _id)
		}

		// Process
		if _entityType == "requestStatus" {
			config.EndpointMethod = "requestStatus"
			config.Scope = "requestStatus.write"

			// Post data via correct authentication method
			switch config.AuthMethod {
			case BASIC:
				acknowledgeCmd_postAcknowledge(true, _entityType, _id)
			case OAUTH2:
				acknowledgeCmd_OAuth02_AuthenticationToken()
				acknowledgeCmd_postAcknowledge(false, _entityType, _id)
			}

		} else if _entityType == "accountingDocument" {
			config.EndpointMethod = "accountingDocuments"
			config.Scope = "accountingDocuments.write"

			for _, key := range _id {
				var keysList []string
				keysList = append(keysList, key)
				// Post data via correct authentication method
				switch config.AuthMethod {
				case BASIC:
					acknowledgeCmd_postAcknowledge(true, _entityType, keysList)
				case OAUTH2:
					acknowledgeCmd_OAuth02_AuthenticationToken()
					acknowledgeCmd_postAcknowledge(false, _entityType, keysList)
				}
			}
		} else {
			log.Fatalf("Unsupported type provided. Allowed values are \"requestStatus\", \"accountingDocument\". ")
		}

	},
}

func init() {
	rootCmd.AddCommand(acknowledgeCmd)

	// set PersistentFlag for entityType and mark it as mandatory!
	acknowledgeCmd.PersistentFlags().String("type", "", "--type requestStatus / --type accountingDocument")
	acknowledgeCmd.MarkPersistentFlagRequired("type")

	acknowledgeCmd.PersistentFlags().StringSlice("id", nil, "--id value1,value2")
	acknowledgeCmd.MarkPersistentFlagRequired("id")
}

func acknowledgeCmd_postAcknowledge(useBasicAuthentication bool, entityType string, keyList []string) {
	// activate debug so that output json is written
	// config.Debug = true

	// setup client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// create request object
	var req *http.Request
	var errHttpRequest error
	if entityType == "accountingDocument" {

		// create new http request message
		if len(keyList) > 0 {
			req, errHttpRequest = http.NewRequest("POST", config.EndpointUrl+"v1/accountingDocuments/"+keyList[0]+"/acknowledge", nil)
			if errHttpRequest != nil {
				log.Fatalf("Got error %s", errHttpRequest.Error())
			}
		} else {
			log.Fatalf("Missing Id?")
		}

	} else if entityType == "requestStatus" {
		var payload components.ProcessingStatusGetBatchRequest
		payload.RequestIds = keyList

		// convert struct to json
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(payload)

		// create new http request message
		req, errHttpRequest = http.NewRequest("POST", config.EndpointUrl+"v1/requestStatus/acknowledge", payloadBuf)
		if errHttpRequest != nil {
			log.Fatalf("Got error %s", errHttpRequest.Error())
		}

	} else {
		log.Fatalf("Unsupported type provided. Allowed values are \"requestStatus\", \"accountingDocument\". ")
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
		acknowledgeCmd_writeJson(string(body), "")
	}
}

func acknowledgeCmd_OAuth02_AuthenticationToken() {
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

func acknowledgeCmd_writeJson(response string, continueationToken string) {
	// create file
	var fileName string
	home, _ := homedir.Dir()
	if continueationToken == "" {
		fileName = "acknowledge_" + config.EndpointMethod + ".json"
	} else {
		fileName = "acknowledge_" + config.EndpointMethod + "_" + continueationToken + ".json"
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
