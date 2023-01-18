/*
Copyright © 2023 Wayne Manselle wayne@viridianforge.tech

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type record struct {
	testName      string
	sourceFile    string
	linesHit      int
	lineCount     int
	functionsHit  int
	functionCount int
	branchesHit   int
	branchCount   int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parse-lcov",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genReport called")
		// Open file
		filename := "/home/vforge/Downloads/parts.info"
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Could not open the file due to this %s error \n", err)
		}
		// Create buffered scanner
		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)
		var fileLines []string

		// Read in lines from file
		for fileScanner.Scan() {
			// Can probably handle parsing here
			fileLines = append(fileLines, fileScanner.Text())
		}
		if err = file.Close(); err != nil {
			fmt.Printf("Could not close the file due to this %s error \n", err)
		}

		curRecord := record{}

		tw := table.NewWriter()
		tw.AppendHeader(table.Row{
			"Test Name", "Source File", "Lines Hit", "Line Count", "Lines Covered",
			"Functions Hit", "Function Count", "Functions Covered", "Branches Hit", "Branch Count", "Branches Covered",
		})

		for _, value := range fileLines {
			contents := strings.Split(value, ":")
			switch code := contents[0]; code {
			case "TN":
				curRecord.testName = contents[1]
			case "SF":
				curRecord.sourceFile = contents[1]
			case "FNF":
				curRecord.functionCount, _ = strconv.Atoi(contents[1])
			case "FNH":
				curRecord.functionsHit, _ = strconv.Atoi(contents[1])
			case "LF":
				curRecord.lineCount, _ = strconv.Atoi(contents[1])
			case "LH":
				curRecord.linesHit, _ = strconv.Atoi(contents[1])
			case "BRF":
				curRecord.branchCount, _ = strconv.Atoi(contents[1])
			case "BRH":
				curRecord.branchesHit, _ = strconv.Atoi(contents[1])
			case "end_of_record":
				tw.AppendRow(table.Row{
					curRecord.testName,
					curRecord.sourceFile,
					curRecord.linesHit,
					curRecord.lineCount,
					curRecord.functionsHit,
					curRecord.functionCount,
					curRecord.branchesHit,
					curRecord.branchCount,
				})
				curRecord = record{}
			default:
				fmt.Println(value)
			}
		}
		tw.SetIndexColumn(1)
		tw.SetTitle("Sample Coverage Report")
		fmt.Println(tw.Render())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.parse-lcov.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
