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
	"math"
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

var lcovReport string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parse-lcov",
	Short: "Parses an lcov report into a human readable format",
	Long: `parse-lcov is a utility to transform a report in lcov format into
a number of human readable formats.
- Formatted Text
- Delimited Files (future)
- HTML Table (future)
- Markdown (future)

Examples:
- ./parse-lcov -r report.info -- Prints the lcov to stdout as text
- ./parse-lcov -r report.info -f csv -- Prints the lcov to stdout as a csv
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Open file
		file, err := os.Open(lcovReport)
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
				// Calculate coverage percentages
				linesCovered := "-"
				functionsCovered := "-"
				branchesCovered := "-"

				if curRecord.lineCount > 0 {
					linesCovered = fmt.Sprintf("%.2f %%", math.Round((float64(curRecord.linesHit)/float64(curRecord.lineCount))*100))
				}

				if curRecord.functionCount > 0 {
					functionsCovered = fmt.Sprintf("%.2f %%", math.Round((float64(curRecord.functionsHit)/float64(curRecord.functionCount))*100))
				}

				if curRecord.branchCount > 0 {
					branchesCovered = fmt.Sprintf("%.2f %%", math.Round((float64(curRecord.branchesHit)/float64(curRecord.branchCount))*100))
				}

				tw.AppendRow(table.Row{
					curRecord.testName,
					curRecord.sourceFile,
					curRecord.linesHit,
					curRecord.lineCount,
					linesCovered,
					curRecord.functionsHit,
					curRecord.functionCount,
					functionsCovered,
					curRecord.branchesHit,
					curRecord.branchCount,
					branchesCovered,
				})
				tw.AppendSeparator()
				curRecord = record{}
			default:
			}
		}
		tw.SetIndexColumn(1)
		tw.SetTitle(lcovReport)
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
	rootCmd.Flags().StringVarP(&lcovReport, "report", "r", "", "lcov report")
	rootCmd.MarkFlagRequired("report")
}
