/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Fetch data from basware open API and store data into Excel",
	Long: `Fetch data from basware open API and store data into Excel:

- multiTool.exe --config <config>.yaml report genericList
- multiTool.exe --config <config>.yaml report requestStatus
- multiTool.exe --config <config>.yaml report vendors
- multiTool.exe --config <config>.yaml report accountingDocuments
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("report called")
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	_ = viper.BindPFlag("authmethod", rootCmd.PersistentFlags().Lookup("authmethod"))
	_ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("api_username", rootCmd.PersistentFlags().Lookup("api_username"))
	_ = viper.BindPFlag("api_clientid", rootCmd.PersistentFlags().Lookup("api_clientid"))
	_ = viper.BindPFlag("api_clientsecret", rootCmd.PersistentFlags().Lookup("api_clientsecret"))
	_ = viper.BindPFlag("api_token", rootCmd.PersistentFlags().Lookup("api_token"))
	_ = viper.BindPFlag("api_endpointurl", rootCmd.PersistentFlags().Lookup("api_endpointurl"))
	_ = viper.BindPFlag("api_password", rootCmd.PersistentFlags().Lookup("api_password"))
	_ = viper.BindPFlag("api_swagger", rootCmd.PersistentFlags().Lookup("api_swagger"))
	_ = viper.BindPFlag("api_pagesize", rootCmd.PersistentFlags().Lookup("api_pagesize"))
	_ = viper.BindPFlag("api_lists", rootCmd.PersistentFlags().Lookup("api_lists"))
}
