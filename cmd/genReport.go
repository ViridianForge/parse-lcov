/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// genReportCmd represents the genReport command
var genReportCmd = &cobra.Command{
	Use:   "genReport",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genReport called")
		// Open file
		filename := "/home/vforge/Downloads/parts.info"
		// Read in chunks
		file, err := os.Open(filename)
		// Print something
		if err != nil {
			fmt.Printf("Could not open the file due to this %s error \n", err)
		}
		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)
		var fileLines []string

		for fileScanner.Scan() {
			fileLines = append(fileLines, fileScanner.Text())
		}
		if err = file.Close(); err != nil {
			fmt.Printf("Could not close the file due to this %s error \n", err)
		}

		for key, value := range fileLines {
			fmt.Printf("line %v : %s \n", key, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(genReportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genReportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genReportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
