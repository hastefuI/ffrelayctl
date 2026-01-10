package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hastefuI/ffrelayctl/api"
	"github.com/spf13/cobra"
)

const (
	version    = "0.6.2"
	envKeyName = "FFRELAYCTL_KEY"
)

var (
	apiKey  string
	baseURL string
	timeout time.Duration = api.DefaultTimeout
	client  *api.Client
	ctx     context.Context
	cancel  context.CancelFunc
)

var rootCmd = &cobra.Command{
	Use:          "ffrelayctl",
	Short:        "Firefox Relay CLI",
	Long:         `ffrelayctl - A CLI for Firefox Relay.`,
	Version:      version,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		ctx, cancel = context.WithCancel(context.Background())
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
				cancel()
			case <-ctx.Done():
				return
			}
		}()

		if apiKey == "" {
			apiKey = os.Getenv(envKeyName)
		}

		if apiKey == "" {
			return fmt.Errorf("no API key provided.\nUse --key <API_KEY> or set the %s environment variable", envKeyName)
		}

		var opts []api.ClientOption
		if baseURL != "" {
			opts = append(opts, api.WithBaseURL(baseURL))
		}
		opts = append(opts, api.WithTimeout(timeout))
		opts = append(opts, api.WithUserAgent("ffrelayctl/"+version))
		opts = append(opts, api.WithContext(ctx))
		client = api.NewClient(apiKey, opts...)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	if cancel != nil {
		cancel()
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", fmt.Sprintf("Base URL for the API (default: %s)", api.DefaultBaseURL))
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", api.DefaultTimeout, "HTTP request timeout (e.g., 15s, 2m)")
}
