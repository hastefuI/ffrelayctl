package cmd

import (
	"fmt"
	"strconv"

	"github.com/hastefuI/ffrelayctl/api"
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

var phonesUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Update a phone mask",
	Long: `Update a phone mask's configuration.

This command allows you to enable or disable call and text forwarding
for a specific phone mask without deleting and recreating it.

Examples:
  ffrelayctl phones update 1 --enabled
  ffrelayctl phones update 1 --disabled`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		var enabled *bool
		if cmd.Flags().Changed("enabled") {
			val := true
			enabled = &val
		}
		if cmd.Flags().Changed("disabled") {
			val := false
			enabled = &val
		}

		if enabled == nil {
			return fmt.Errorf("must specify either --enabled or --disabled")
		}

		req := api.UpdateRelayNumberRequest{
			Enabled: enabled,
		}

		number, err := client.UpdateRelayNumber(id, req)
		if err != nil {
			return err
		}
		return printJSON(number)
	},
}

func init() {
	rootCmd.AddCommand(phonesCmd)
	phonesCmd.AddCommand(phonesListCmd)
	phonesCmd.AddCommand(phonesUpdateCmd)

	phonesUpdateCmd.Flags().Bool("enabled", false, "Enable call/text forwarding")
	phonesUpdateCmd.Flags().Bool("disabled", false, "Disable call/text forwarding")
	phonesUpdateCmd.MarkFlagsMutuallyExclusive("enabled", "disabled")
}
