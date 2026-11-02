package cdep

import "strings"

// Languages a service filter can select. Every service is built by either the
// Go pipeline or one of the JS pipelines, so a filter picks one side.
const (
	LanguageGo = "go"
	LanguageJS = "js"
)

// LanguageForImage reports which pipeline builds a docker image. The Go
// pipeline publishes go-services and its variants (e.g.
// go-services-playwright); every other image (js-services, js-quote,
// web-underwriter, ...) comes from a JS pipeline.
func LanguageForImage(image string) string {
	if strings.HasPrefix(image, "go-services") || strings.HasPrefix(image, "go_services") {
		return LanguageGo
	}

	return LanguageJS
}
