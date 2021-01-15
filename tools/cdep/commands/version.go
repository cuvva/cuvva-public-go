package commands

import (
	"fmt"

	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cdep",
	Long:  "All software has versions.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cdep v%s\n", cdep.Version)
	},
}
