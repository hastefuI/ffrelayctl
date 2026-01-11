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

var phonesContactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "List inbound contacts for phone masks",
	Long: `List contacts who have called or texted your phone masks.

This command shows all inbound contacts across all your phone masks,
including phone numbers, call/text counts, and blocking status.

Note: This feature requires a premium subscription with phone masks enabled.
If you don't have a premium subscription, you'll receive a 404 error.

Examples:
  ffrelayctl phones contacts`,
	RunE: func(cmd *cobra.Command, args []string) error {
		contacts, err := client.ListInboundContacts()
		if err != nil {
			return err
		}
		return printJSON(contacts)
	},
}

var phonesContactsUpdateCmd = &cobra.Command{
	Use:   "contacts-update <ID>",
	Short: "Update an inbound contact",
	Long: `Update an inbound contact's settings, such as blocking or unblocking them.

This command allows you to block or unblock specific phone numbers from
calling or texting your phone masks.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones contacts-update 12345 --block
  ffrelayctl phones contacts-update 12345 --unblock`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		var blocked *bool
		if cmd.Flags().Changed("block") {
			val := true
			blocked = &val
		}
		if cmd.Flags().Changed("unblock") {
			val := false
			blocked = &val
		}

		if blocked == nil {
			return fmt.Errorf("must specify either --block or --unblock")
		}

		req := api.UpdateInboundContactRequest{
			Blocked: blocked,
		}

		contact, err := client.UpdateInboundContact(id, req)
		if err != nil {
			return err
		}
		return printJSON(contact)
	},
}

func init() {
	rootCmd.AddCommand(phonesCmd)
	phonesCmd.AddCommand(phonesListCmd)
	phonesCmd.AddCommand(phonesUpdateCmd)
	phonesCmd.AddCommand(phonesContactsCmd)
	phonesCmd.AddCommand(phonesContactsUpdateCmd)

	phonesUpdateCmd.Flags().Bool("enabled", false, "Enable call/text forwarding")
	phonesUpdateCmd.Flags().Bool("disabled", false, "Disable call/text forwarding")
	phonesUpdateCmd.MarkFlagsMutuallyExclusive("enabled", "disabled")

	phonesContactsUpdateCmd.Flags().Bool("block", false, "Block this contact")
	phonesContactsUpdateCmd.Flags().Bool("unblock", false, "Unblock this contact")
	phonesContactsUpdateCmd.MarkFlagsMutuallyExclusive("block", "unblock")
}
