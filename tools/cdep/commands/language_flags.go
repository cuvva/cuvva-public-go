package commands

import (
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/spf13/cobra"
)

// addLanguageFlags registers the flags that limit an update to the services
// built by one pipeline.
func addLanguageFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("go-only", false, "Only update services running a Go image (go-services, go-services-playwright)")
	cmd.Flags().Bool("js-only", false, "Only update services running a JS image (js-services, js-quote, web-underwriter, ...)")
}

// languageFilter returns the language picked with --go-only or --js-only, or ""
// when neither was given.
func languageFilter(cmd *cobra.Command) (string, error) {
	goOnly, err := cmd.Flags().GetBool("go-only")
	if err != nil {
		return "", err
	}

	jsOnly, err := cmd.Flags().GetBool("js-only")
	if err != nil {
		return "", err
	}

	switch {
	case goOnly && jsOnly:
		return "", cher.New("conflicting_flags", nil)
	case goOnly:
		return cdep.LanguageGo, nil
	case jsOnly:
		return cdep.LanguageJS, nil
	default:
		return "", nil
	}
}
