package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hastefuI/ffrelayctl/api"
	"github.com/hastefuI/ffrelayctl/output"
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
		cfg := GetConfig(cmd)
		numbers, err := cfg.Client.ListRelayNumbers()
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,numbers)
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
		cfg := GetConfig(cmd)
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

		number, err := cfg.Client.UpdateRelayNumber(id, req)
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,number)
	},
}

var phonesDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover available phone numbers",
	Long: `Discover available phone numbers for creating a new phone mask.

This command returns available phone numbers organized by categories:
- Same prefix: Numbers with the same area code and prefix as your real number
- Same area: Numbers with the same area code as your real number
- Other areas: Numbers from different area codes
- Random: Random phone numbers from various locations

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones discover`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		suggestions, err := cfg.Client.GetRelayNumberSuggestions()
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,suggestions)
	},
}

var phonesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for available phone masks by area code",
	Long: `Search for available phone masks filtered by a specific area code.

This command returns an array of available phone numbers matching the
specified area code.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones search --areacode 430`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		areaCode, err := cmd.Flags().GetString("areacode")
		if err != nil {
			return fmt.Errorf("failed to get areacode flag: %w", err)
		}

		if areaCode == "" {
			return fmt.Errorf("--areacode flag is required")
		}

		numbers, err := cfg.Client.SearchRelayNumbers(areaCode)
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,numbers)
	},
}

var phonesForwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "Manage forwarding number (real phone number)",
	Long:  `Register, verify, and manage your real phone number used for forwarding calls and texts from phone masks (premium).`,
}

var phonesForwardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered forwarding numbers",
	Long: `List all registered real phone numbers used for phone mask forwarding.

This command shows your registered phone numbers and their verification status.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones forward list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		phones, err := cfg.Client.GetRealPhone()
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat,phones)
	},
}

var phonesForwardGetCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "Get a specific forwarding number",
	Long: `Retrieve a specific registered real phone number by ID.

This command shows details of a specific forwarding number.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones forward get 1234`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		phones, err := cfg.Client.GetRealPhone()
		if err != nil {
			return err
		}

		for _, phone := range phones {
			if phone.ID == id {
				return output.Print(cfg.OutputFormat,phone)
			}
		}

		return fmt.Errorf("forwarding number with ID %d not found", id)
	},
}

var phonesForwardRegisterCmd = &cobra.Command{
	Use:   "register <phone_number>",
	Short: "Register a forwarding number",
	Long: `Register your real phone number for receiving forwarded calls and texts from phone masks.

After registration, you'll receive an SMS with a verification code that you'll need
to verify using the 'phones forward verify' command.

The phone number must be in E.164 format (e.g., +15551234567).

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones forward register +15551234567`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		req := api.RegisterRealPhoneRequest{
			Number: args[0],
		}
		phone, err := cfg.Client.RegisterRealPhone(req)
		if err != nil {
			if apiErr, ok := err.(*api.APIError); ok {
				fmt.Fprintln(cmd.OutOrStdout(), apiErr.Body)
				return nil
			}
			return err
		}
		return output.Print(cfg.OutputFormat,phone)
	},
}

var phonesForwardVerifyCmd = &cobra.Command{
	Use:   "verify <ID> <phone_number> <verification_code>",
	Short: "Verify a forwarding number",
	Long: `Verify your real phone number using the SMS verification code sent during registration.

The phone number must be in E.164 format (e.g., +15551234567).
The verification code is a 6-digit code sent via SMS.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones forward verify 12040 +15551234567 123456`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		req := api.VerifyRealPhoneRequest{
			Number:           args[1],
			VerificationCode: args[2],
		}
		phone, err := cfg.Client.VerifyRealPhone(id, req)
		if err != nil {
			if apiErr, ok := err.(*api.APIError); ok {
				fmt.Fprintln(cmd.OutOrStdout(), apiErr.Body)
				return nil
			}
			return err
		}
		return output.Print(cfg.OutputFormat,phone)
	},
}

var phonesForwardDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Delete a forwarding number",
	Long: `Remove the registered real phone number from your account.

This will disconnect your real phone number from phone mask forwarding.

Note: This feature requires a premium subscription with phone masks enabled.

Examples:
  ffrelayctl phones forward delete 1                # Delete with confirmation
  ffrelayctl phones forward delete 1 --force        # Delete without confirmation`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return fmt.Errorf("failed to get force flag: %w", err)
		}

		if !force {
			fmt.Printf("Are you sure you want to delete forwarding number %d? This cannot be undone. [y/N]: ", id)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read confirmation: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		if err := cfg.Client.DeleteRealPhone(id); err != nil {
			return err
		}
		fmt.Printf("Forwarding number %d deleted successfully.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(phonesCmd)
	phonesCmd.AddCommand(phonesListCmd)
	phonesCmd.AddCommand(phonesUpdateCmd)
	phonesCmd.AddCommand(phonesDiscoverCmd)
	phonesCmd.AddCommand(phonesSearchCmd)
	phonesCmd.AddCommand(phonesForwardCmd)

	phonesForwardCmd.AddCommand(phonesForwardListCmd)
	phonesForwardCmd.AddCommand(phonesForwardGetCmd)
	phonesForwardCmd.AddCommand(phonesForwardRegisterCmd)
	phonesForwardCmd.AddCommand(phonesForwardVerifyCmd)
	phonesForwardCmd.AddCommand(phonesForwardDeleteCmd)

	phonesUpdateCmd.Flags().Bool("enabled", false, "Enable call/text forwarding")
	phonesUpdateCmd.Flags().Bool("disabled", false, "Disable call/text forwarding")
	phonesUpdateCmd.MarkFlagsMutuallyExclusive("enabled", "disabled")

	phonesSearchCmd.Flags().String("areacode", "", "Area code to filter phone numbers")
	phonesSearchCmd.MarkFlagRequired("areacode")

	phonesForwardDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
