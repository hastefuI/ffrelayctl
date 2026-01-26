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

var (
	randomMask *bool
)

type commonUpdateFields struct {
	enabled         *bool
	description     *string
	blockListEmails *bool
}

func parseCommonUpdateFlags(cmd *cobra.Command) (commonUpdateFields, error) {
	var fields commonUpdateFields

	if cmd.Flags().Changed("enabled") {
		enabled, err := cmd.Flags().GetBool("enabled")
		if err != nil {
			return fields, fmt.Errorf("failed to get enabled flag: %w", err)
		}
		fields.enabled = &enabled
	}
	if cmd.Flags().Changed("disabled") {
		disabled := false
		fields.enabled = &disabled
	}
	if cmd.Flags().Changed("description") {
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return fields, fmt.Errorf("failed to get description flag: %w", err)
		}
		fields.description = &description
	}
	if cmd.Flags().Changed("block-list") {
		blockList := true
		fields.blockListEmails = &blockList
	}
	if cmd.Flags().Changed("no-block-list") {
		blockList := false
		fields.blockListEmails = &blockList
	}

	return fields, nil
}

var masksCmd = &cobra.Command{
	Use:   "masks",
	Short: "Manage email masks (both random and custom domain)",
	Long: `Manage your Firefox Relay email masks.

The --random flag controls which mask types to work with:
  (no flag)       - List: shows all masks; Get: tries both types; Create/Update/Delete: random masks
  --random=true   - Work with random relay addresses only
  --random=false  - Work with custom domain addresses only (Premium required)`,
}

var masksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all masks",
	Long: `List all email masks.

Examples:
  ffrelayctl masks list                # List all masks (both random and custom domain)
  ffrelayctl masks list --random=true  # List only random masks
  ffrelayctl masks list --random=false # List only custom domain masks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		if randomMask == nil {
			relayAddresses, err := cfg.Client.ListRelayAddresses()
			if err != nil {
				return err
			}
			domainAddresses, err := cfg.Client.ListDomainAddresses()
			if err != nil {
				return err
			}

			combined := make([]output.CombinedMask, 0)
			for _, addr := range relayAddresses {
				combined = append(combined, output.CombinedMask{Type: "random", Mask: addr})
			}
			for _, addr := range domainAddresses {
				combined = append(combined, output.CombinedMask{Type: "custom", Mask: addr})
			}

			return output.Print(cfg.OutputFormat, combined)
		}

		if *randomMask {
			addresses, err := cfg.Client.ListRelayAddresses()
			if err != nil {
				return err
			}
			return output.Print(cfg.OutputFormat, addresses)
		} else {
			addresses, err := cfg.Client.ListDomainAddresses()
			if err != nil {
				return err
			}
			return output.Print(cfg.OutputFormat, addresses)
		}
	},
}

var masksGetCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "Get a specific mask",
	Long: `Get details of a specific email mask by ID.

When --random flag is not specified, automatically checks both random and custom domain masks.

Examples:
  ffrelayctl masks get 12345                # Try random first, then custom domain if premium
  ffrelayctl masks get 12345 --random=true  # Get random mask only
  ffrelayctl masks get 12345 --random=false # Get custom domain mask only`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		if randomMask != nil {
			if *randomMask {
				address, err := cfg.Client.GetRelayAddress(id)
				if err != nil {
					return err
				}
				return output.Print(cfg.OutputFormat, address)
			} else {
				address, err := cfg.Client.GetDomainAddress(id)
				if err != nil {
					return err
				}
				return output.Print(cfg.OutputFormat, address)
			}
		}

		address, err := cfg.Client.GetRelayAddress(id)
		if err == nil {
			return output.Print(cfg.OutputFormat, address)
		}

		profiles, profileErr := cfg.Client.GetProfiles()
		if profileErr != nil {
			return err
		}

		if len(profiles) > 0 && profiles[0].HasPremium {
			domainAddress, domainErr := cfg.Client.GetDomainAddress(id)
			if domainErr == nil {
				return output.Print(cfg.OutputFormat, domainAddress)
			}
			return err
		}
		return err
	},
}

var masksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new mask",
	Long: `Create a new email mask.

For random masks (--random=true, default):
  ffrelayctl masks create --description "Shopping" --generated-for "amazon.com"

For custom domain masks (--random=false, Premium required):
  ffrelayctl masks create --random=false --address "shopping" --description "Shopping sites"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return fmt.Errorf("failed to get description flag: %w", err)
		}
		blockList, err := cmd.Flags().GetBool("block-list")
		if err != nil {
			return fmt.Errorf("failed to get block-list flag: %w", err)
		}
		disabled, err := cmd.Flags().GetBool("disabled")
		if err != nil {
			return fmt.Errorf("failed to get disabled flag: %w", err)
		}

		if randomMask == nil || *randomMask {
			generatedFor, err := cmd.Flags().GetString("generated-for")
			if err != nil {
				return fmt.Errorf("failed to get generated-for flag: %w", err)
			}
			usedOn, err := cmd.Flags().GetString("used-on")
			if err != nil {
				return fmt.Errorf("failed to get used-on flag: %w", err)
			}

			req := api.CreateRelayAddressRequest{
				Enabled:         !disabled,
				Description:     description,
				GeneratedFor:    generatedFor,
				UsedOn:          usedOn,
				BlockListEmails: blockList,
			}

			address, err := cfg.Client.CreateRelayAddress(req)
			if err != nil {
				return err
			}
			return output.Print(cfg.OutputFormat, address)
		} else {
			address, err := cmd.Flags().GetString("address")
			if err != nil {
				return fmt.Errorf("failed to get address flag: %w", err)
			}
			if address == "" {
				return fmt.Errorf("--address is required for custom domain masks (--random=false)")
			}

			req := api.CreateDomainAddressRequest{
				Address:         address,
				Enabled:         !disabled,
				Description:     description,
				BlockListEmails: blockList,
			}

			domainAddress, err := cfg.Client.CreateDomainAddress(req)
			if err != nil {
				return err
			}
			return output.Print(cfg.OutputFormat, domainAddress)
		}
	},
}

var masksUpdateCmd = &cobra.Command{
	Use:   "update <ID>",
	Short: "Update a mask",
	Long: `Update an existing email mask.

Examples:
  ffrelayctl masks update 12345 --disabled
  ffrelayctl masks update 12345 --description "New description"
  ffrelayctl masks update 12345 --random=false --enabled`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		fields, err := parseCommonUpdateFlags(cmd)
		if err != nil {
			return err
		}

		if randomMask == nil || *randomMask {
			req := api.UpdateRelayAddressRequest{
				Enabled:         fields.enabled,
				Description:     fields.description,
				BlockListEmails: fields.blockListEmails,
			}
			if cmd.Flags().Changed("used-on") {
				usedOn, err := cmd.Flags().GetString("used-on")
				if err != nil {
					return fmt.Errorf("failed to get used-on flag: %w", err)
				}
				req.UsedOn = &usedOn
			}

			address, err := cfg.Client.UpdateRelayAddress(id, req)
			if err != nil {
				return err
			}
			return output.Print(cfg.OutputFormat, address)
		} else {
			req := api.UpdateDomainAddressRequest{
				Enabled:         fields.enabled,
				Description:     fields.description,
				BlockListEmails: fields.blockListEmails,
			}

			address, err := cfg.Client.UpdateDomainAddress(id, req)
			if err != nil {
				return err
			}
			return output.Print(cfg.OutputFormat, address)
		}
	},
}

var masksDeleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Delete a mask",
	Long: `Delete an email mask.

Examples:
  ffrelayctl masks delete 12345                      # Delete with confirmation
  ffrelayctl masks delete 12345 --force              # Delete without confirmation
  ffrelayctl masks delete 12345 --random=false       # Delete custom domain mask`,
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

		maskType := "random mask"
		if randomMask != nil && !*randomMask {
			maskType = "custom domain mask"
		}

		if !force {
			fmt.Printf("Are you sure you want to delete %s %d? This cannot be undone. [y/N]: ", maskType, id)
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

		if randomMask == nil || *randomMask {
			if err := cfg.Client.DeleteRelayAddress(id); err != nil {
				return err
			}
			fmt.Printf("Random mask %d deleted successfully.\n", id)
		} else {
			if err := cfg.Client.DeleteDomainAddress(id); err != nil {
				return err
			}
			fmt.Printf("Custom domain mask %d deleted successfully.\n", id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(masksCmd)
	masksCmd.AddCommand(masksListCmd)
	masksCmd.AddCommand(masksGetCmd)
	masksCmd.AddCommand(masksCreateCmd)
	masksCmd.AddCommand(masksUpdateCmd)
	masksCmd.AddCommand(masksDeleteCmd)
	masksCmd.PersistentFlags().Bool("random", false, "Filter by mask type: true for random masks, false for custom domain masks")

	defaultPreRunEList := masksListCmd.PreRunE
	masksListCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("random") {
			val, err := cmd.Flags().GetBool("random")
			if err != nil {
				return fmt.Errorf("failed to get random flag: %w", err)
			}
			randomMask = &val
		}
		if defaultPreRunEList != nil {
			return defaultPreRunEList(cmd, args)
		}
		return nil
	}

	defaultPreRunEGet := masksGetCmd.PreRunE
	masksGetCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("random") {
			val, err := cmd.Flags().GetBool("random")
			if err != nil {
				return fmt.Errorf("failed to get random flag: %w", err)
			}
			randomMask = &val
		}
		if defaultPreRunEGet != nil {
			return defaultPreRunEGet(cmd, args)
		}
		return nil
	}

	for _, subCmd := range []*cobra.Command{masksCreateCmd, masksUpdateCmd, masksDeleteCmd} {
		defaultPreRunE := subCmd.PreRunE
		subCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("random") {
				val, err := cmd.Flags().GetBool("random")
				if err != nil {
					return fmt.Errorf("failed to get random flag: %w", err)
				}
				randomMask = &val
			} else {
				defaultVal := true
				randomMask = &defaultVal
			}
			if defaultPreRunE != nil {
				return defaultPreRunE(cmd, args)
			}
			return nil
		}
	}

	masksCreateCmd.Flags().String("description", "", "Description for the mask")
	masksCreateCmd.Flags().String("generated-for", "", "Site the mask was generated for (random masks only)")
	masksCreateCmd.Flags().String("used-on", "", "Site the mask is used on (random masks only)")
	masksCreateCmd.Flags().String("address", "", "Local part of the address (custom domain masks only, required)")
	masksCreateCmd.Flags().Bool("block-list", false, "Block promotional emails")
	masksCreateCmd.Flags().Bool("disabled", false, "Create in disabled state")
	masksUpdateCmd.Flags().Bool("enabled", false, "Enable the mask")
	masksUpdateCmd.Flags().Bool("disabled", false, "Disable the mask")
	masksUpdateCmd.Flags().String("description", "", "Update description")
	masksUpdateCmd.Flags().String("used-on", "", "Update used on (random masks only)")
	masksUpdateCmd.Flags().Bool("block-list", false, "Block promotional emails")
	masksUpdateCmd.Flags().Bool("no-block-list", false, "Don't block promotional emails")
	masksUpdateCmd.MarkFlagsMutuallyExclusive("enabled", "disabled")
	masksUpdateCmd.MarkFlagsMutuallyExclusive("block-list", "no-block-list")
	masksDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
