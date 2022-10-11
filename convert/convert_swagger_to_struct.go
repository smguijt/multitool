/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Swagger.json to GO stuct",
	Long: `reads open api specifications (swagger.json) and
	transforms the schema definition to GO struct objects
	conversion data stored under /%user%/.convert/components directory`,
	Run: func(cmd *cobra.Command, args []string) {

		// transform swagger.json to map
		list := transform_swagger_to_struct()

		// sort list
		keys := make([]string, 0, len(list))
		for k := range list {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// retrieve homedir. No  need to verify errors as directory is already validated
		home, _ := homedir.Dir()
		log.Printf("convert::Home directory: %v", home)

		// create .multiTool folder
		os.Mkdir(home+"/.multiTool", 0755)

		// create .convert folder
		os.Mkdir(home+"/.multiTool/.convert", 0755)
		// create components sub folder
		os.Mkdir(home+"/.multiTool/.convert/components", 0755)
		// set output dir
		componentsDir := home + "/.multiTool/.convert/components/"

		// process sorted list
		for _, k := range keys {
			fmt.Println(k)

			// create output file
			writeElement(componentsDir, k+".go", list[k])
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

type ComponentSchemaElement struct {
	Required    bool
	Name        string
	Type        string
	MaxLength   float64
	MinLength   float64
	Description string
	Example     string
	Format      string
	Items       []string
	Nullable    bool
	Enum        []string
	Pattern     string
	Minimum     float64
	Maximum     float64
	Reference   string
	ReadOnly    bool
} // var componentEl map[string][]ComponentSchemaElement

/*
	Transforms swagger.json document to structure
	out: map of ComponentSchemaElements
*/
func transform_swagger_to_struct() map[string][]ComponentSchemaElement {
	jsonParsed := readSwagger()
	schemaFieldList, totalSchemaItems := getSchemaElements(jsonParsed)

	fmt.Printf("retrieved %v schema items\n", totalSchemaItems)
	componentEl := make(map[string][]ComponentSchemaElement)
	for i := 0; i < totalSchemaItems; i++ {
		// fmt.Printf("%v\n", schemaFieldList[i])
		componentEl[schemaFieldList[i]] = nil
		//if schemaFieldList[i] == "AccountEntity" {
		retValue := getElementFields(schemaFieldList[i], jsonParsed)
		componentEl[schemaFieldList[i]] = retValue[schemaFieldList[i]]
		//}
	}

	// return sorted map
	return componentEl
}

/* check if item exists in array
    in: s = array to check
        str = value to validate
	out: bool
*/
func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

/* Reads Swagger.JSON document page either from disc or online
   returns gabs.Container
*/
func readSwagger() *gabs.Container {

	// check if local page exists
	// check for swagger.json
	e1Json, err := ioutil.ReadFile("swagger.json")
	// init vars
	var jsonParsed *gabs.Container
	var jsonErr error

	// if err then swagger.json file does not exists so fetch online
	if err != nil {

		// create httpClient
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		// fetch content
		req, err := http.NewRequest("GET", config.ApiSwaggerUrl, nil)
		if err != nil {
			log.Fatalf("Got error %s", err.Error())
		}

		// add header params (form urlencoded and length)
		data := url.Values{}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		// send request
		response, err := client.Do(req)
		if err != nil {
			log.Fatalf("Got error %s", err.Error())
		} else {
			log.Println("response status: ", response.StatusCode)
		}
		defer response.Body.Close()

		// process response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		jsonParsed, jsonErr = gabs.ParseJSON([]byte(string(body)))
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	} else {
		// swagger.json exists so conver to gabs object
		jsonParsed, jsonErr = gabs.ParseJSON([]byte(e1Json))
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	}

	return jsonParsed
}

/* retrieve element names under the component/swagger tree
   in: *gabs.Container
   out: []string, int
*/
func getSchemaElements(data *gabs.Container) ([]string, int) {
	var list []string
	x1, _ := data.S("components", "schemas").ChildrenMap()
	for el := range x1 {
		// fmt.Printf("element: %v\n", el)
		list = append(list, el)
	}
	return list, len(list)
}

/* retrieve element names under the retrieved schema element list
    in: elementName: name of the element to search for
	    *gabs.Container
    out: []string, int
*/
func getElementFields(elementName string, data *gabs.Container) map[string][]ComponentSchemaElement {

	// create map object
	//var componentEl map[string][]ComponentSchemaElement
	// init map object
	componentEl := make(map[string][]ComponentSchemaElement)
	componentEl[elementName] = nil

	// var list []string
	var listRequiredFields []string
	x1, _ := data.S("components", "schemas", elementName).ChildrenMap()

	// retrieve required fields
	for el, obj := range x1 {
		switch el {
		case "required":
			x2, _ := obj.Children()
			for _, childData := range x2 {
				x2Value := childData.Data().(string)
				listRequiredFields = append(listRequiredFields, x2Value)
			}
		}
	}

	// retrieve properties
	for el, obj := range x1 {
		switch el {
		case "properties":
			x2, _ := obj.ChildrenMap()
			for childEl, childData := range x2 {
				// ComponentSchemaElement
				var property_element ComponentSchemaElement

				// get property name
				property_name := childEl
				property_element.Name = strings.ToUpper(property_name[0:1]) + property_name[1:]
				// fmt.Printf("%v\n", property_name)

				// get property elements
				property_elements, _ := childData.ChildrenMap()
				for propertyItem, propertyItemValue := range property_elements {
					switch propertyItem {
					case "type":
						property_element.Type = propertyItemValue.Data().(string)
					case "maxLength":
						property_element.MaxLength = propertyItemValue.Data().(float64)
					case "minLength":
						property_element.MinLength = propertyItemValue.Data().(float64)
					case "description":
						property_element.Description = propertyItemValue.Data().(string)
					case "example":
						tmpExample := propertyItemValue.Data()
						// fmt.Println(reflect.TypeOf(tmpExample).String())

						switch reflect.TypeOf(tmpExample).String() {
						case "float64":
							property_element.Example = fmt.Sprint(tmpExample)
						case "bool":
							property_element.Example = fmt.Sprint(tmpExample)
						default:
							property_element.Example = propertyItemValue.Data().(string)
						}

					case "items":
						// only to be populated if type = array
						_items, _ := propertyItemValue.ChildrenMap()
						for _, _itemItemValue := range _items {
							_tmpItemValueData := _itemItemValue.Data().(string)
							_tmpItemValueData = strings.ReplaceAll(_tmpItemValueData, "#/components/schemas/", "")
							property_element.Items = append(property_element.Items, _tmpItemValueData)
						}

					case "format":
						property_element.Format = propertyItemValue.Data().(string)
					case "nullable":
						property_element.Nullable = propertyItemValue.Data().(bool)
					case "enum":
						_items, _ := propertyItemValue.ChildrenMap()
						for _, _itemItemValue := range _items {
							_tmpItemValueData := _itemItemValue.Data().(string)
							property_element.Enum = append(property_element.Enum, _tmpItemValueData)
						}
					case "pattern":
						property_element.Pattern = propertyItemValue.Data().(string)
					case "minimum":
						property_element.Minimum = propertyItemValue.Data().(float64)
					case "maximum":
						property_element.Maximum = propertyItemValue.Data().(float64)
					case "$ref":
						tmpReference := propertyItemValue.Data().(string)
						tmpReference = strings.ReplaceAll(tmpReference, "#/components/schemas/", "")
						property_element.Reference = tmpReference

						if len(property_element.Type) == 0 {
							property_element.Type = fmt.Sprintf("%s%s", strings.ToUpper(tmpReference[0:1]), tmpReference[1:])
						}

					case "readOnly":
						property_element.ReadOnly = propertyItemValue.Data().(bool)
					default:
						fmt.Printf("%v - %v\n", propertyItem, propertyItemValue)
					}
				}
				if isElementExist(listRequiredFields, property_element.Name) {
					property_element.Required = true
				}

				// append object to map
				componentEl[elementName] = append(componentEl[elementName], property_element)
			}

		}
	}

	return componentEl
}

/*
	writeElement: Create file on disc
*/
func writeElement(path string, filename string, data []ComponentSchemaElement) {

	// create file
	file, err := os.Create(path + filename)
	if err != nil {
		log.Fatal(err)
	}

	// init buffer
	writer := bufio.NewWriterSize(file, 10)

	// build
	var stringList []string
	stringList = append(stringList, "package components")
	stringList = append(stringList, "")
	stringList = append(stringList, fmt.Sprintf("type %s struct {", strings.Split(filename, ".")[0]))
	for _, d := range data {
		var buildStr string
		switch d.Type {
		case "array":
			iArrayLength := 20 - len(d.Items[0]) + 10
			if len(d.Items[0]) < 10 {
				iArrayLength = iArrayLength - 2
			} else {
				iArrayLength = iArrayLength + 2
			}
			buildStr = fmt.Sprintf("	%-25s []%-10s %-"+strconv.Itoa(iArrayLength)+"s %s", d.Name, d.Items[0], "", json_build_string(d))
			//
		case "enum":
			//

		case "boolean":
			iArrayLength := 20 - len(d.Type) + 11
			buildStr = fmt.Sprintf("	%-25s %-10s %-"+strconv.Itoa(iArrayLength)+"s %s", d.Name, "bool", "", json_build_string(d))
		case "integer":
			iArrayLength := 20 - len(d.Type) + 11
			buildStr = fmt.Sprintf("	%-25s %-10s %-"+strconv.Itoa(iArrayLength)+"s %s", d.Name, "int", "", json_build_string(d))
		case "string":
			if len(d.Type) == 0 {
				fmt.Println("")
			}
			iArrayLength := 20 - len(d.Type) + 10
			buildStr = fmt.Sprintf("	%-25s %-10s %-"+strconv.Itoa(iArrayLength)+"s %s", d.Name, d.Type, "", json_build_string(d))
		case "number":
			iArrayLength := 20 - len(d.Type) + 10
			buildStr = fmt.Sprintf("	%-25s %-10s %-"+strconv.Itoa(iArrayLength)+"s %s", d.Name, "float64", "", json_build_string(d))
		default:
			iArrayLength := 20 - len(d.Type) + 14
			//if len(d.Type) < 10 {
			//	iArrayLength = iArrayLength - 2
			//} else {
			//	iArrayLength = iArrayLength + 2
			//}
			buildStr = fmt.Sprintf("	%-25s %-10s %-"+strconv.Itoa(iArrayLength)+"s %s", d.Name, d.Type, "", json_build_string(d))
			//
		}

		stringList = append(stringList, buildStr)
	}
	stringList = append(stringList, "}")

	fmt.Sprintln(stringList)

	// ** create content **
	for _, line := range stringList {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}
	// ** end create content **

	// flush and close file
	writer.Flush()
}

func json_build_string(data ComponentSchemaElement) string {
	var jsonStr string

	// `json:"externalCode" validate:"required"`
	// `json:"externalCode"`
	// `json:"externalCode" validate:"required,min=3,max=12"`
	// `json:"externalCode" validate:"required,numeric"`
	// `json:"externalCode,omitempty" validate:"required,numeric"`

	// set validation rules
	if data.Required {
		jsonStr += "required"
	}

	if data.MinLength > 0 {
		jsonStr += fmt.Sprintf(",min=%s", strconv.Itoa(int(data.MinLength)))
	}

	if data.MaxLength > 0 {
		jsonStr += fmt.Sprintf(",max=%s", strconv.Itoa(int(data.MaxLength)))
	}

	if len(jsonStr) > 0 {
		jsonStr = " validate:\"" + jsonStr[1:] + "\""
	}

	// build actual string
	if data.Nullable {
		jsonStr = "json:\"" + data.Name + ",omitempty\"" + jsonStr
	} else {
		jsonStr = "json:\"" + data.Name + "\"" + jsonStr
	}

	// return json Str
	return "`" + jsonStr + "`"
}
