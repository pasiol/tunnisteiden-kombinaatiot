/*
Copyright © 2022 Steve Francia <spf@spf13.com>

*/
package cmd

import (
	"log"
	"riverian-tunnisteiden-kombinaatiot/internal"

	"github.com/spf13/cobra"
)

// getReportCmd represents the getReport command
var getReportCmd = &cobra.Command{
	Use:   "getReport",
	Short: "Generates a csv report.",
	Long:  `Generates a csv report not valid combinations of the bookings.`,
	Run: func(cmd *cobra.Command, args []string) {
		month, _ := cmd.Flags().GetString("month")
		log.Printf("%s", month)
		err := internal.GetReport(month)
		if err != nil {
			log.Fatalf("getting report failed: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getReportCmd)
	getReportCmd.Flags().String("month", "", "month tag for example 09.2022")
}
