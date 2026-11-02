package parsers

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

func TestParseAll(t *testing.T) {
	p, err := Parse([]string{"service", "avocado", "all"}, false)
	if err != nil {
		t.Fatal(err)
	}

	if !p.AllItems || len(p.Items) != 0 {
		t.Fatalf("got AllItems=%v Items=%v, want all services", p.AllItems, p.Items)
	}

	p.Branch = "master"
	p.Only = "go"

	if got, want := p.String("update"), "update service avocado all --go-only"; got != want {
		t.Fatalf("String: got %q, want %q", got, want)
	}
}

func TestParseAllErrors(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"service", "avocado", "all", "sms"}, "all_must_be_alone"},
		{[]string{"lambda", "avocado", "all"}, "all_requires_services"},
	}

	for _, tt := range tests {
		_, err := Parse(tt.args, false)
		if code := cher.Coerce(err).Code; code != tt.want {
			t.Errorf("%v: got %q, want %q (err: %v)", tt.args, code, tt.want, err)
		}
	}
}

func TestParseNamedServices(t *testing.T) {
	p, err := Parse([]string{"service", "avocado", "sms", "web-underwriter"}, false)
	if err != nil {
		t.Fatal(err)
	}

	if p.AllItems || len(p.Items) != 2 || p.Items[0] != "service-sms" || p.Items[1] != "web-underwriter" {
		t.Fatalf("got %+v", p)
	}
}
