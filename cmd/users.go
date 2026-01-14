package cmd

import (
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `View and manage Firefox Relay users.`,
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long: `List all users associated with your Firefox Relay account.

Examples:
  ffrelayctl users list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		users, err := cfg.Client.ListUsers()
		if err != nil {
			return err
		}
		return printJSON(users)
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersListCmd)
}
