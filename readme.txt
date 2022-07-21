When run the tool for the first time, the tool generates a folder under your homedir. The name of the folder is .multiTool
Within this folder a file called multiTool.yaml is created.  
Copy the file located in this .zip document and replace it with the one under the .multiTool folder

The file contains the following information

api_clientid: Basware API OAUTH02 authentication username
api_clientsecret: Basware API OAUTH02 authentication password
api_endpointmethod: requestStatus
api_endpointurl: location for the basware API.
api_pagesize: number of items to fetch at the time. When not set, default = 500
api_scope: requestStatus.read
authmethod: 2
debug: enables debug mode. When active .json files are written for the request within the .multiTool folder

to run: open command line and execute command.  "multiTool --config <config.yaml> report requestStatus"
the file will write Excel document requestStatus.xlsx under the .multiTool folder
for more information see the documents under the documentation folder

to compile an executable please run from the terminal: go build multiTool.go
