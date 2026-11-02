package app

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()

	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	return dir
}

const goCommit = "7fbcf728fe0682e62f3c5e179ebf2e846a99d3c2"

func TestResolveServiceFile(t *testing.T) {
	dir := writeFiles(t, map[string]string{
		"_base.json":             `{"docker_image_name": "go-services", "branch": "master", "commit": "` + goCommit + `"}`,
		"service-own-go.json":    `{"docker_image_name": "go-services-playwright"}`,
		"service-own-js.json":    `{"docker_image_name": "js-services", "commit": "` + goCommit + `"}`,
		"service-inherits.json":  `{"internal_keys": []}`,
		"_base.yaml":             "podDisruptionBudget:\n  enabled: false\n",
		"service-helm.yaml":      "image:\n  repository: 005717268539.dkr.ecr.eu-west-1.amazonaws.com/js-quote\n  tag: master-abc\n",
		"service-helm-bare.yaml": "replicas: 2\n",
	})

	tests := []struct {
		file      string
		wantImage string
		wantLang  string
	}{
		{"_base.json", "go-services", "go"},
		{"service-own-go.json", "go-services-playwright", "go"},
		{"service-own-js.json", "js-services", "js"},
		{"service-inherits.json", "go-services", "go"},
		{"service-helm.yaml", "js-quote", "js"},
		{"service-helm-bare.yaml", "", ""},
		{"_base.yaml", "", ""},
		{"service-missing.json", "go-services", "go"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, err := resolveServiceFile(filepath.Join(dir, tt.file))
			if err != nil {
				t.Fatal(err)
			}

			if f.Image != tt.wantImage {
				t.Errorf("image: got %q, want %q", f.Image, tt.wantImage)
			}

			if f.Language() != tt.wantLang {
				t.Errorf("language: got %q, want %q", f.Language(), tt.wantLang)
			}
		})
	}
}

func TestListAndFilterServiceFiles(t *testing.T) {
	dir := writeFiles(t, map[string]string{
		"_base.json":         `{"docker_image_name": "go-services", "commit": "` + goCommit + `"}`,
		"service-sms.json":   `{}`,
		"service-quote.json": `{"docker_image_name": "js-quote"}`,
		"service-web.yaml":   "image:\n  repository: 005717268539.dkr.ecr.eu-west-1.amazonaws.com/go-services\n",
		"README.md":          "not a config file",
	})

	files, err := listServiceFiles(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("got %d files without env files, want 3: %v", len(files), files)
	}

	files, err = listServiceFiles(dir, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 4 {
		t.Fatalf("got %d files with env files, want 4: %v", len(files), files)
	}

	tests := []struct {
		only         string
		wantSelected string
		wantSkipped  string
	}{
		{"", "_base, quote, sms, web", "(none)"},
		{"go", "_base, sms, web", "quote"},
		{"js", "quote", "_base, sms, web"},
	}

	for _, tt := range tests {
		t.Run("only="+tt.only, func(t *testing.T) {
			selected, skipped := filterServiceFiles("avocado", tt.only, files)

			if got := joinServiceNames(selected); got != tt.wantSelected {
				t.Errorf("selected: got %q, want %q", got, tt.wantSelected)
			}

			if got := joinServiceNames(skipped); got != tt.wantSkipped {
				t.Errorf("skipped: got %q, want %q", got, tt.wantSkipped)
			}
		})
	}
}

func TestListServiceFilesMissingDir(t *testing.T) {
	files, err := listServiceFiles(filepath.Join(t.TempDir(), "nope"), true)
	if err != nil || len(files) != 0 {
		t.Fatalf("got %v, %v; want no files and no error", files, err)
	}
}

func TestDefinesOwnVersion(t *testing.T) {
	dir := writeFiles(t, map[string]string{
		"pinned.json":   `{"commit": "` + goCommit + `"}`,
		"tagged.json":   `{"tag": "master-` + goCommit + `"}`,
		"inherits.json": `{"docker_image_name": "js-services"}`,
		"pinned.yaml":   "image:\n  tag: master-abc\n",
		"inherits.yaml": "replicas: 1\n",
	})

	tests := map[string]bool{
		"pinned.json":   true,
		"tagged.json":   true,
		"inherits.json": false,
		"pinned.yaml":   true,
		"inherits.yaml": false,
	}

	for name, want := range tests {
		if got := (serviceFile{Path: filepath.Join(dir, name)}).definesOwnVersion(); got != want {
			t.Errorf("%s: got %v, want %v", name, got, want)
		}
	}
}

func TestECRCheckerLooksUpEachImageOnce(t *testing.T) {
	lookups := map[string]int{}

	c := newECRChecker(goCommit, "master")
	c.find = func(image, commit, branch string) error {
		lookups[image]++
		return nil
	}

	for _, f := range []serviceFile{
		{Path: "a.json", Image: "go-services"},
		{Path: "b.json", Image: "go-services"},
		{Path: "c.json", Image: "js-services"},
		{Path: "d.json", Image: ""},
	} {
		c.Check(f)
	}

	if lookups["go-services"] != 1 || lookups["js-services"] != 1 || len(lookups) != 2 {
		t.Fatalf("unexpected lookups: %v", lookups)
	}
}
