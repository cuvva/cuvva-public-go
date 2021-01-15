package main

import (
	"fmt"
	"os"

	cmd "github.com/cuvva/cuvva-public-go/tools/oncallcalc/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(cmd.GenerateReportCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "oncallcalc",
	Short: "CLI tool to help manage the on call rota",
	Long:  "CLI tool to help manage the on call rota",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
