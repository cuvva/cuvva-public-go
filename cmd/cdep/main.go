package main

import (
	"fmt"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	cmd "github.com/cuvva/cuvva-public-go/tools/cdep/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.UpdateCmd)
	rootCmd.AddCommand(cmd.UpdateWebCmd)
	rootCmd.AddCommand(cmd.UpdateDefaultCmd)

	if err := rootCmd.Execute(); err != nil {
		if v, ok := err.(cher.E); ok {
			outStr := "Error:\n"

			if humanMessage, ok := cdep.ErrorCodeMapping[v.Code]; ok {
				outStr = fmt.Sprintf("%s%s\n", outStr, humanMessage)
			}

			if a, ok := v.Meta["path"]; ok {
				outStr = fmt.Sprintf("%sPath: %s\n", outStr, a)
			}

			if image, ok := v.Meta["image"]; ok {
				if tag, ok := v.Meta["tag"]; ok {
					outStr = fmt.Sprintf("%sImage: %s:%s\n", outStr, image, tag)
				} else {
					outStr = fmt.Sprintf("%sImage: %s\n", outStr, image)
				}
			}

			if a, ok := v.Meta["error"]; ok {
				outStr = fmt.Sprintf("%sCause: %s\n", outStr, a)
			}

			if a, ok := v.Meta["allowed"]; ok {
				if allowedOptions, ok := a.([]string); ok {
					outStr = fmt.Sprintf("%s\nAllowed options:\n", outStr)

					for _, option := range allowedOptions {
						outStr = fmt.Sprintf("%s - %s\n", outStr, option)
					}
				}
			}

			fmt.Println(outStr)
		}
	}
}

var rootCmd = &cobra.Command{
	Use: "cdep",
}
