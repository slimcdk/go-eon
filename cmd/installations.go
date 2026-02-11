package cmd

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

var installationsCmd = &cobra.Command{
	Use:   "installations",
	Short: "Get installations with metadata",
	Long: `Retrieve a list of installations connected to your account.
Optionally filter by specific installation IDs.`,
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetStringSlice("filter")

		installations, err := clientInstance.GetInstallations(filter)
		cobra.CheckErr(err)

		gout.MustPrint(installations)
	},
}

var measurementSeriesCmd = &cobra.Command{
	Use:   "measurement-series",
	Short: "Get measurement series for all installations",
	Long: `Retrieve a list of available measurement series for each installation.
Use the series IDs to fetch actual measurement data.`,
	Run: func(cmd *cobra.Command, args []string) {
		series, err := clientInstance.GetMeasurementSeries()
		cobra.CheckErr(err)

		gout.MustPrint(series)
	},
}

func init() {
	installationsCmd.Flags().StringSlice("filter", nil, "Filter by installation IDs")

	rootCmd.AddCommand(installationsCmd)
	rootCmd.AddCommand(measurementSeriesCmd)
}
