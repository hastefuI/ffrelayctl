package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"go.hasteful.org/ffrelayctl/api"
	"go.hasteful.org/ffrelayctl/output"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all Firefox Relay account data",
	Long: `Export all data from a Firefox Relay account.

This command fetches all masks, phones, profiles, and contacts from a
Firefox Relay account for backup purposes.

The account's API token is left out unless --include-api-token is passed, so a
backup file does not carry a credential for the whole account by default.
Including it asks for confirmation first, which --force skips.

Examples:
  ffrelayctl export
  ffrelayctl export > ffrelay-backup.json
  ffrelayctl export --include-api-token > ffrelay-backup.json
  ffrelayctl export --include-api-token --force > ffrelay-backup.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)

		includeAPIToken, err := cmd.Flags().GetBool("include-api-token")
		if err != nil {
			return fmt.Errorf("failed to get include-api-token flag: %w", err)
		}
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return fmt.Errorf("failed to get force flag: %w", err)
		}

		// The prompt goes to stderr so it cannot land in a redirected export.
		if includeAPIToken && !force {
			confirmed, err := confirm(cmd.ErrOrStderr(), "The export will include your API token, which grants full access to your account. Continue? [y/N]: ")
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(cmd.ErrOrStderr(), "Export cancelled.")
				return nil
			}
		}

		type exportData struct {
			Masks    []output.CombinedMask `json:"masks"`
			Phones   []api.RelayNumber     `json:"phones"`
			Profiles []api.Profile         `json:"profiles"`
			Contacts []api.InboundContact  `json:"contacts"`
			Users    []api.User            `json:"users"`
		}

		var (
			wg     sync.WaitGroup
			mu     sync.Mutex
			errors []error
			result exportData
		)

		wg.Go(func() {
			relayAddresses, err := cfg.Client.ListRelayAddresses(cfg.Ctx)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch relay addresses: %w", err))
				mu.Unlock()
				return
			}
			domainAddresses, err := cfg.Client.ListDomainAddresses(cfg.Ctx)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch domain addresses: %w", err))
				mu.Unlock()
				return
			}

			combined := make([]output.CombinedMask, 0, len(relayAddresses)+len(domainAddresses))
			for _, addr := range relayAddresses {
				combined = append(combined, output.CombinedMask{Type: "random", Mask: addr})
			}
			for _, addr := range domainAddresses {
				combined = append(combined, output.CombinedMask{Type: "custom", Mask: addr})
			}

			mu.Lock()
			result.Masks = combined
			mu.Unlock()
		})

		wg.Go(func() {
			numbers, err := cfg.Client.ListRelayNumbers(cfg.Ctx)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch relay numbers: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			result.Phones = numbers
			mu.Unlock()
		})

		wg.Go(func() {
			profiles, err := cfg.Client.GetProfiles(cfg.Ctx)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch profiles: %w", err))
				mu.Unlock()
				return
			}
			if !includeAPIToken {
				for i := range profiles {
					profiles[i].APIToken = ""
				}
			}

			mu.Lock()
			result.Profiles = profiles
			mu.Unlock()
		})

		wg.Go(func() {
			contacts, err := cfg.Client.ListInboundContacts(cfg.Ctx)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch inbound contacts: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			result.Contacts = contacts
			mu.Unlock()
		})

		wg.Go(func() {
			users, err := cfg.Client.ListUsers(cfg.Ctx)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch users: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			result.Users = users
			mu.Unlock()
		})

		wg.Wait()

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			}
			return fmt.Errorf("failed to export data: %d error(s) occurred", len(errors))
		}

		return output.Print(cfg.OutputFormat, result)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().Bool("include-api-token", false, "Include the account's API token in the export")
	exportCmd.Flags().Bool("force", false, "Skip the confirmation prompt for --include-api-token")
}
