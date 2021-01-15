package main

import (
	"fmt"
	"os"

	cmd "github.com/cuvva/cuvva-public-go/tools/ctxcheck/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(cmd.ScanCmd)

	cmd.ScanCmd.Flags().StringP("path", "p", "", "Determine the base path to scan for context misuse")
	cmd.ScanCmd.MarkFlagRequired("path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ctxcheck",
	Short: "Tool to find misuse of context",
	Long:  "A CLI tool to help find misuse of context.Context within errgroup",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
