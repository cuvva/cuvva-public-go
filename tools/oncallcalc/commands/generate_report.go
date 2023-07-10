package commands

import (
	"context"
	"errors"
	"time"

	"github.com/cuvva/cuvva-public-go/tools/oncallcalc/app"
	"github.com/cuvva/cuvva-public-go/tools/oncallcalc/config"
	"github.com/cuvva/cuvva-public-go/tools/oncallcalc/lib/govuk"
	"github.com/spf13/cobra"
)

var ScheduleIDs []string
var TimeIn string
var Verbose bool

func init() {
	GenerateReportCmd.Flags().StringArrayVarP(&ScheduleIDs, "schedule_id", "s", []string{"PKICNIO", "PKW2AKK"}, "Schedule IDs from PagerDuty")
	GenerateReportCmd.Flags().StringVarP(&TimeIn, "time", "t", "Jan 2020", "Which month and year should we look at?")
	GenerateReportCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode?")
}

var GenerateReportCmd = &cobra.Command{
	Use:     "generate-report",
	Aliases: []string{"gr"},
	Short:   "Generate reports on a month",
	Long:    "Given a month, generate a report on the payouts for anyone who was oncall.",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := config.BuildPDClient()
		if err != nil {
			return err
		}

		govUK := govuk.New()

		app := app.New(client, govUK)

		timeIn, err := convMonthToMonth(TimeIn)
		if err != nil {
			if Verbose {
				cmd.Printf("%+v\n", timeIn)
			}
			return err
		}

		if err := app.GetBankHolidays(context.Background(), timeIn.Year()); err != nil {
			return err
		}

		rota, debug, err := app.GenerateRota(ScheduleIDs, timeIn.Year(), timeIn.Month())
		if err != nil {
			if Verbose {
				cmd.Printf("%+v\n", debug)
			}
			return err
		}

		cmd.Printf("\n")
		cmd.Println(app.StringifyRota(rota))

		return nil
	},
}

func convMonthToMonth(in string) (time.Time, error) {
	if in == "" {
		return time.Now(), errors.New("you must provide a --month param")
	}

	t, err := time.Parse("Jan 2006", in)
	if err != nil {
		return time.Now(), err
	}

	return t, nil
}
