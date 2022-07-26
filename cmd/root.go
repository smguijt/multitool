/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AuthenticationMethod int

const (
	BASIC AuthenticationMethod = iota + 1
	OAUTH2
)

type AccessTokenResponse struct {
	Access_Token string `mapstructure:"access_token"`
	Scope        string `mapstructure:"scope"`
	Expires_In   int    `mapstructure:"expires_in"`
	Token_Type   string `mapstructure:"token_type"`
}

type Config struct {
	Version        string               `mapstructure:"version"`
	Debug          bool                 `mapstructure:"debug"`
	AuthMethod     AuthenticationMethod `mapstructure:"authmethod"`
	NTLogin        string               `mapstructure:"username"`
	Username       string               `mapstructure:"api_username"`
	ClientId       string               `mapstructure:"api_clientId"`
	ClientSecret   string               `mapstructure:"api_clientSecret"`
	Token          string               `mapstructure:"api_token"`
	EndpointUrl    string               `mapstructure:"api_endpointurl"`
	Password       string               `mapstructure:"api_password"`
	ApiSwaggerUrl  string               `mapstructure:"api_swagger"`
	Scope          string               `mapstructure:"api_scope"`
	EndpointMethod string               `mapstructure:"api_endpointmethod"`
	PageSize       int                  `mapstructure:"api_pageSize"`
	ExportFormat   string               `mapstructure:"exportformat"`
	GenericLists   string               `mapstructure:"api_lists"`
	Prefix         string               `mapstructure:"prefix"`
	EntityType     string               `mapstructure:"api_entity_type"`
}

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multiTool.exe",
	Short: "cli for multiple purpose usage",
	Long:  `cli for multiple purpose usage`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("Inside rootCmd Run with args: %v\n", args)
		// fmt.Println("viper keys: ", viper.AllKeys())
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
	},
}

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// create .multiTool folder
	home, _ := homedir.Dir()
	log.Printf("multiTool::Home directory: %v", home)
	os.Mkdir(home+"/.multiTool", 0755)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.multiTool/multiTool.yaml)")
	// rootCmd.PersistentFlags().BoolP("debug", "d", false, "prints debug messages")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Bind viper to these flags so viper can read flag values along with config, env, etc.
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
	_ = viper.BindPFlag("exportformat", rootCmd.PersistentFlags().Lookup("exportformat"))

}

func initConfig() {
	// set version
	version := "1.0.3"
	config.Version = version
	config.ExportFormat = "excel"

	// fmt.Println("", viper.GetViper().Get("authmethod"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home + "/.multiTool")
		viper.SetConfigName("multiTool")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %v", viper.ConfigFileUsed())
		log.Printf("version: %v", config.Version)
		log.Printf("Copyright: %s ", "2022, Sietse Guijt")
		log.Printf("Export Format: %s", config.ExportFormat)

		// set prefix
		config.Prefix = strings.ReplaceAll(viper.ConfigFileUsed(), ".yaml", "")
		if config.Prefix == "multiTool" {
			config.Prefix = ""
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("%v", err)
		}

	} else {
		// file does not exist, so create default
		log.Printf("Missing configuration file. Creating default template...")

		config.Debug = false
		config.AuthMethod = OAUTH2
		config.EndpointUrl = "https://test-api.basware.com/"

		err = viper.WriteConfig()
		if err != nil {
			log.Printf("%v", err)
		}
		viper.SafeWriteConfig()

		//initial file is created so exit
		log.Printf("Configuration file created! Please update and restart the program!")
		os.Exit(1)
	}

	// runtime example: multiTool --authmethod 0 --config debug.yaml
	// parsing --authmethod parses the AutMethod parameter to the config instead of using
	// the one from the config file itself.
	//
	// parsing --config debug.yaml overides the debug file to use from multiTool.yaml to debug.yaml
	// fmt.Println("config:authmethod::", config.AuthMethod)

}
