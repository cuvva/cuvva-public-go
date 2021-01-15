package commands

import (
	"go/ast"
	"log"
	"os"

	"github.com/cuvva/cuvva-public-go/tools/ctxcheck"
	"github.com/spf13/cobra"
)

// ScanCmd is the cobra defintion for the "scan" command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all files to find context mis-use",
	RunE: func(cmd *cobra.Command, args []string) error {
		searchDir, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		// Below predicates check for a function called 'Go' where the first argument is 'ctx'
		var functionPredicate ctxcheck.FunctionPredicate = func(i *ast.Ident) bool {
			return i.Name == "Go"
		}

		var argumentPredicate ctxcheck.ArgPredicate = func(index int, ident *ast.Ident) *ast.Ident {
			if index == 0 && ident.Name == "ctx" {
				return ident
			}

			return nil
		}

		v := ctxcheck.NewVisitor(functionPredicate, argumentPredicate)

		err = ctxcheck.Walk(searchDir, v)
		if err != nil {
			return err
		}

		for _, match := range v.Matches {
			log.Printf("gctx misuse in file %s%s at pos %d\n", searchDir, match.Path, match.Pos)
		}

		if len(v.Matches) > 0 {
			os.Exit(0)
		}

		return nil
	},
}
