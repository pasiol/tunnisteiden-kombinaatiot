/*
Copyright © 2022 Steve Francia <spf@spf13.com>

*/
package cmd

import (
	"fmt"

	"riverian-tunnisteiden-kombinaatiot/internal"

	"github.com/spf13/cobra"
)

// updateIdentifiersCmd represents the updateIdentifiers command
var updateIdentifiersCmd = &cobra.Command{
	Use:   "updateIdentifiers",
	Short: "Updating identifiers",
	Long: `Exporting identifiers from csv file to table identifiers.
	Truncating table before importing the new data.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		data := []internal.AccountingIdentifiers{}
		data, err := internal.ReadAccountingIdentifiers(filename)
		if err != nil {
			fmt.Printf("reading identifiers failed: %s", err)
		}
		err = internal.TruncateIdentifiersTable()
		if err != nil {
			fmt.Printf("truncating identifiers failed: %s", err)
		}
		err = internal.UpdateAccountingIdentifiers(data)
		if err != nil {
			fmt.Printf("updating identifiers failed: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateIdentifiersCmd)
	updateIdentifiersCmd.Flags().String("filename", "", "csv-file of accounting identifiers")
}
