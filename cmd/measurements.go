package cmd

import (
	"strconv"
	"time"

	"github.com/drewstinnett/gout/v2"
	"github.com/slimcdk/go-eon/eon"
	"github.com/spf13/cobra"
)

var measurementsCmd = &cobra.Command{
	Use:   "measurements <series-id>",
	Short: "Get measurement data for a measurement series",
	Long: `Retrieve measurement values for a specific measurement series ID.

Resolution options:
  - quarter: 15-minute intervals (requires from/to, max 3 months)
  - hour: Hourly values (requires from/to, max 1 year)
  - day: Daily values
  - month: Monthly values`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seriesID, err := strconv.Atoi(args[0])
		cobra.CheckErr(err)

		fromFlag, _ := cmd.Flags().GetString("from")
		toFlag, _ := cmd.Flags().GetString("to")
		resolution, _ := cmd.Flags().GetString("resolution")
		includeMissing, _ := cmd.Flags().GetBool("include-missing")

		var from, to time.Time
		if fromFlag != "" {
			from, err = time.Parse(time.DateOnly, fromFlag)
			cobra.CheckErr(err)
		}
		if toFlag != "" {
			to, err = time.Parse(time.DateOnly, toFlag)
			cobra.CheckErr(err)
		}

		measurements, err := clientInstance.GetMeasurements(
			seriesID,
			eon.Resolution(resolution),
			from,
			to,
			includeMissing,
		)
		cobra.CheckErr(err)

		gout.MustPrint(measurements)
	},
}

func init() {
	measurementsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	measurementsCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	measurementsCmd.Flags().String("resolution", "hour", "Resolution: quarter, hour, day, month")
	measurementsCmd.Flags().Bool("include-missing", false, "Fill in missing values")

	rootCmd.AddCommand(measurementsCmd)
}
