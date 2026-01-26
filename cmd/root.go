package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hastefuI/ffrelayctl/api"
	"github.com/hastefuI/ffrelayctl/output"
	"github.com/spf13/cobra"
)

const (
	envKeyName = "FFRELAYCTL_KEY"
)

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

type CmdConfig struct {
	APIKey       string
	BaseURL      string
	Timeout      time.Duration
	OutputFormat string
	Client       *api.Client
	Ctx          context.Context
	Cancel       context.CancelFunc
	VersionInfo  VersionInfo
}

type configKey struct{}

func GetConfig(cmd *cobra.Command) *CmdConfig {
	return cmd.Context().Value(configKey{}).(*CmdConfig)
}

var rootCmd = &cobra.Command{
	Use:          "ffrelayctl",
	Short:        "Firefox Relay CLI",
	Long:         `ffrelayctl - A CLI for Firefox Relay.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		cfg := GetConfig(cmd)

		cfg.APIKey, _ = cmd.Flags().GetString("key")
		cfg.BaseURL, _ = cmd.Flags().GetString("base-url")
		cfg.Timeout, _ = cmd.Flags().GetDuration("timeout")
		cfg.OutputFormat, _ = cmd.Flags().GetString("output")

		if !output.IsValidFormat(cfg.OutputFormat) {
			return fmt.Errorf("invalid output format %q: must be one of [text|json]", cfg.OutputFormat)
		}

		cfg.Ctx, cfg.Cancel = context.WithCancel(cmd.Context())
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			defer func() {
				signal.Stop(sigChan)
				close(sigChan)
			}()
			select {
			case sig := <-sigChan:
				fmt.Fprintf(os.Stderr, "\nReceived signal %v, cancelling request...\n", sig)
				cfg.Cancel()
			case <-cfg.Ctx.Done():
				return
			}
		}()

		if cfg.APIKey == "" {
			cfg.APIKey = os.Getenv(envKeyName)
		}

		if cfg.APIKey == "" {
			return fmt.Errorf("no API key provided.\nUse --key <API_KEY> or set the %s environment variable", envKeyName)
		}

		var opts []api.ClientOption
		if cfg.BaseURL != "" {
			opts = append(opts, api.WithBaseURL(cfg.BaseURL))
		}
		opts = append(opts, api.WithTimeout(cfg.Timeout))
		opts = append(opts, api.WithUserAgent("ffrelayctl/"+cfg.VersionInfo.Version))
		opts = append(opts, api.WithContext(cfg.Ctx))
		cfg.Client = api.NewClient(cfg.APIKey, opts...)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().String("base-url", "", fmt.Sprintf("Base URL for the API (default: %s)", api.DefaultBaseURL))
	rootCmd.PersistentFlags().String("key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringP("output", "o", output.FormatText, "Output format [text|json]")
	rootCmd.PersistentFlags().Duration("timeout", api.DefaultTimeout, "HTTP request timeout (e.g., 15s, 2m)")
}

func Execute(vi VersionInfo) {
	cfg := &CmdConfig{
		Timeout:      api.DefaultTimeout,
		OutputFormat: output.FormatText,
		VersionInfo:  vi,
	}

	ctxWithConfig := context.WithValue(context.Background(), configKey{}, cfg)
	rootCmd.SetContext(ctxWithConfig)

	rootCmd.Version = vi.Version
	rootCmd.SetVersionTemplate(fmt.Sprintf("ffrelayctl version %s\ncommit: %s\nbuilt at: %s\n", vi.Version, vi.Commit, vi.Date))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	if cfg.Cancel != nil {
		cfg.Cancel()
	}
}
