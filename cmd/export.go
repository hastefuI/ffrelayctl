package cmd

import (
	"fmt"
	"sync"

	"github.com/hastefuI/ffrelayctl/api"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all Firefox Relay account data",
	Long: `Export all data from a Firefox Relay account.

This command fetches all masks, phones, profiles, and contacts from a
Firefox Relay account for backup purposes.

Examples:
  ffrelayctl export
  ffrelayctl export > ffrelay-backup.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)

		type combinedMask struct {
			Type string      `json:"type"`
			Mask interface{} `json:"mask"`
		}

		type exportData struct {
			Masks    []combinedMask       `json:"masks"`
			Phones   []api.RelayNumber    `json:"phones"`
			Profiles []api.Profile        `json:"profiles"`
			Contacts []api.InboundContact `json:"contacts"`
		}

		var (
			wg     sync.WaitGroup
			mu     sync.Mutex
			errors []error
			result exportData
		)

		wg.Add(4)

		go func() {
			defer wg.Done()
			relayAddresses, err := cfg.Client.ListRelayAddresses()
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch relay addresses: %w", err))
				mu.Unlock()
				return
			}
			domainAddresses, err := cfg.Client.ListDomainAddresses()
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch domain addresses: %w", err))
				mu.Unlock()
				return
			}

			combined := make([]combinedMask, 0, len(relayAddresses)+len(domainAddresses))
			for _, addr := range relayAddresses {
				combined = append(combined, combinedMask{Type: "random", Mask: addr})
			}
			for _, addr := range domainAddresses {
				combined = append(combined, combinedMask{Type: "custom", Mask: addr})
			}

			mu.Lock()
			result.Masks = combined
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			numbers, err := cfg.Client.ListRelayNumbers()
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch relay numbers: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			result.Phones = numbers
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			profiles, err := cfg.Client.GetProfiles()
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch profiles: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			result.Profiles = profiles
			mu.Unlock()
		}()

		go func() {
			defer wg.Done()
			contacts, err := cfg.Client.ListInboundContacts()
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch inbound contacts: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			result.Contacts = contacts
			mu.Unlock()
		}()

		wg.Wait()

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			}
			return fmt.Errorf("failed to export data: %d error(s) occurred", len(errors))
		}

		return printJSON(result)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
