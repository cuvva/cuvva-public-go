package cdep

import "testing"

func TestLanguageForImage(t *testing.T) {
	tests := map[string]string{
		"go-services":            LanguageGo,
		"go_services":            LanguageGo,
		"go-services-playwright": LanguageGo,
		"js-services":            LanguageJS,
		"js-quote":               LanguageJS,
		"web-underwriter":        LanguageJS,
	}

	for image, want := range tests {
		if got := LanguageForImage(image); got != want {
			t.Errorf("%s: got %q, want %q", image, got, want)
		}
	}
}
