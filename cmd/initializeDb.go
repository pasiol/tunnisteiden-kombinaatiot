/*
Copyright © 2022 Steve Francia <spf@spf13.com>

*/
package cmd

import (
	"riverian-tunnisteiden-kombinaatiot/internal"

	"github.com/spf13/cobra"
)

// initializeDbCmd represents the initializeDb command
var initializeDbCmd = &cobra.Command{
	Use:   "initializeDb",
	Short: "Initializing db",
	Long: `Initializing db creating tables bookings and identifiers.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.InitializeDb()
	},
}

func init() {
	rootCmd.AddCommand(initializeDbCmd)
}
