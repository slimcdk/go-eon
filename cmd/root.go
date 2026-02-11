package cmd

import (
	"os"

	"github.com/drewstinnett/gout/v2"
	"github.com/drewstinnett/gout/v2/formats/json"
	"github.com/slimcdk/go-eon/eon"
	"github.com/spf13/cobra"
)

// clientInstance holds the Eon API client
var clientInstance eon.Client

var rootCmd = &cobra.Command{
	Use:   "eon",
	Short: "A CLI for the Eon Energy Navigator API",
	Long: `Access Eon energy data including installations, measurements, and costs.

Credentials are loaded from environment variables CLIENT_ID and CLIENT_SECRET,
or can be provided via command-line flags.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")

		if clientID != "" && clientSecret != "" {
			clientInstance = eon.NewWithCredentials(clientID, clientSecret)
		} else {
			clientInstance = eon.New()
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	gout.SetFormatter(json.Formatter{})

	rootCmd.PersistentFlags().String("client-id", "", "Eon API client ID (env: CLIENT_ID)")
	rootCmd.PersistentFlags().String("client-secret", "", "Eon API client secret (env: CLIENT_SECRET)")
}
