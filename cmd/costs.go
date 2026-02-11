package cmd

import (
	"time"

	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

var costsCmd = &cobra.Command{
	Use:   "costs <installation-id>",
	Short: "Get cost data for an installation",
	Long: `Retrieve cost information for a specific installation.
Whole months are considered for the time range.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installationID := args[0]

		fromFlag, _ := cmd.Flags().GetString("from")
		toFlag, _ := cmd.Flags().GetString("to")

		var from, to *time.Time
		var err error

		if fromFlag != "" {
			fromTime, err := time.Parse(time.DateOnly, fromFlag)
			cobra.CheckErr(err)
			from = &fromTime
		}
		if toFlag != "" {
			toTime, err := time.Parse(time.DateOnly, toFlag)
			cobra.CheckErr(err)
			to = &toTime
		}

		costs, err := clientInstance.GetCosts(installationID, from, to)
		cobra.CheckErr(err)

		gout.MustPrint(costs)
	},
}

func init() {
	costsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	costsCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")

	rootCmd.AddCommand(costsCmd)
}
