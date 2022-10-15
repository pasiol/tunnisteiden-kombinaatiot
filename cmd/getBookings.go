/*
Copyright © 2022 Steve Francia <spf@spf13.com>

*/
package cmd

import (
	
	"fmt"

	"riverian-tunnisteiden-kombinaatiot/internal"

	"github.com/spf13/cobra"
)

// getBookingsCmd represents the getBookings command
var getBookingsCmd = &cobra.Command{
	Use:   "getBookings",
	Short: "Importing bookings to pg database.",
	Long: `Importing bookings from Atlas db service to bookings table.
	Truncating table before importing the new data.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.TruncateBookingsTable()
		if err != nil {
			fmt.Printf("truncating bookings failed: %s", err)
		}
		internal.GetBookins()
	},
}

func init() {
	rootCmd.AddCommand(getBookingsCmd)
}
