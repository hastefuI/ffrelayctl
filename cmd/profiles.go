package cmd

import (
	"github.com/spf13/cobra"
	"go.hasteful.org/ffrelayctl/output"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage Firefox Relay profiles",
	Long:  `View and manage your Firefox Relay profiles.`,
}

var profilesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List user profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		profiles, err := cfg.Client.GetProfiles(cfg.Ctx)
		if err != nil {
			return err
		}
		return output.Print(cfg.OutputFormat, profiles)
	},
}

func init() {
	rootCmd.AddCommand(profilesCmd)
	profilesCmd.AddCommand(profilesListCmd)
}
