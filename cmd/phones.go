package cmd

import (
	"github.com/spf13/cobra"
)

var phonesCmd = &cobra.Command{
	Use:   "phones",
	Short: "Manage phone masks (relay numbers)",
	Long:  `View and manage your Firefox Relay phone masks (Premium).`,
}

var phonesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all phone masks",
	Long: `List all phone masks.

Phone masks are a premium feature that provides virtual phone numbers
that forward calls and texts to your real phone number.

Examples:
  ffrelayctl phones list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		numbers, err := client.ListRelayNumbers()
		if err != nil {
			return err
		}
		return printJSON(numbers)
	},
}

func init() {
	rootCmd.AddCommand(phonesCmd)
	phonesCmd.AddCommand(phonesListCmd)
}
