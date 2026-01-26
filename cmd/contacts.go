package cmd

import (
	"fmt"
	"strconv"

	"github.com/hastefuI/ffrelayctl/api"
	"github.com/hastefuI/ffrelayctl/output"
	"github.com/spf13/cobra"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage inbound contacts for phone masks",
	Long:  `View and manage contacts who have called or texted your phone masks (Premium).`,
}

var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List inbound contacts for phone masks",
	Long: `List contacts who have called or texted your phone masks.

This command shows all inbound contacts across all your phone masks,
including phone numbers, call/text counts, and blocking status.

Note: This feature requires a premium subscription with phone masks enabled.
If you don't have a premium subscription, you'll receive a 404 error.

Examples:
  ffrelayctl contacts list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		contacts, err := cfg.Client.ListInboundContacts()
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,contacts)
	},
}

var contactsUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Update an inbound contact",
	Long: `Update an inbound contact's settings, such as blocking or unblocking them.

This command allows you to block or unblock specific phone numbers from
calling or texting your phone masks.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl contacts update 12345 --block
  ffrelayctl contacts update 12345 --unblock`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
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

		contact, err := cfg.Client.UpdateInboundContact(id, req)
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,contact)
	},
}

func init() {
	rootCmd.AddCommand(contactsCmd)
	contactsCmd.AddCommand(contactsListCmd)
	contactsCmd.AddCommand(contactsUpdateCmd)

	contactsUpdateCmd.Flags().Bool("block", false, "Block this contact")
	contactsUpdateCmd.Flags().Bool("unblock", false, "Unblock this contact")
	contactsUpdateCmd.MarkFlagsMutuallyExclusive("block", "unblock")
}
