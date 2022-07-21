/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tasksCmd represents the tasks command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Retrieves all data from the taskStatus API and creates Excel document of the content",
	Long: `Retrieves all data from the taskStatus API and creates Excel document of the content. 
	
	taskStatus gives an overview of all running background tasks that is beiing performed by the open API platform. 
	An example is delete records.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("taskStatus not yet implemented!!")
	},
}

func init() {
	reportCmd.AddCommand(tasksCmd)
}
