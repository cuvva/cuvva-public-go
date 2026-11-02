package app

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
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
	}{
		{"_base.json", "go-services"},
		{"service-own-go.json", "go-services-playwright"},
		{"service-own-js.json", "js-services"},
		{"service-inherits.json", "go-services"},
		{"service-helm.yaml", "js-quote"},
		{"service-helm-bare.yaml", ""},
		{"_base.yaml", ""},
		{"service-missing.json", "go-services"},
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
		})
	}
}

func TestListServiceFiles(t *testing.T) {
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

	if got, want := joinServiceNames(files), "quote, sms, web"; got != want {
		t.Fatalf("without env files: got %q, want %q", got, want)
	}

	files, err = listServiceFiles(dir, true)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := joinServiceNames(files), "_base, quote, sms, web"; got != want {
		t.Fatalf("with env files: got %q, want %q", got, want)
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

// fakeBuilds returns a buildChecker that knows which images have been built,
// and counts how often each image is looked up.
func fakeBuilds(force bool, built map[string]bool, lookups map[string]int) *buildChecker {
	c := newBuildChecker(goCommit, "master", force)
	c.find = func(image, commit, branch string) (bool, error) {
		lookups[image]++

		if image == "broken" {
			return false, errors.New("ecr is down")
		}

		return built[image], nil
	}

	return c
}

var selectFiles = []serviceFile{
	{Path: "_base.json", Image: "go-services"},
	{Path: "_base.yaml"},
	{Path: "service-sms.json", Image: "go-services"},
	{Path: "service-auth.json", Image: "go-services"},
	{Path: "service-agent.json", Image: "go-services-playwright"},
	{Path: "service-quote.json", Image: "js-quote"},
	{Path: "service-mystery.json"},
}

func TestSelectSkipsServicesWithoutABuild(t *testing.T) {
	lookups := map[string]int{}
	c := fakeBuilds(false, map[string]bool{"go-services": true, "js-quote": true}, lookups)

	selected, skipped, err := c.Select("avocado", selectFiles, false)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := joinServiceNames(selected), "_base, _base, auth, quote, sms"; got != want {
		t.Errorf("selected: got %q, want %q", got, want)
	}

	if got, want := joinServiceNames(skipped), "agent, mystery"; got != want {
		t.Errorf("skipped: got %q, want %q", got, want)
	}

	// each image is looked up once, however many services share it
	if lookups["go-services"] != 1 || lookups["js-quote"] != 1 || lookups["go-services-playwright"] != 1 || len(lookups) != 3 {
		t.Errorf("unexpected lookups: %v", lookups)
	}
}

func TestSelectStrictFailsOnMissingBuild(t *testing.T) {
	c := fakeBuilds(false, map[string]bool{"go-services": true}, map[string]int{})

	_, _, err := c.Select("avocado", []serviceFile{{Path: "service-agent.json", Image: "go-services-playwright"}}, true)
	if code := cher.Coerce(err).Code; code != "build_not_found" {
		t.Fatalf("got %q, want build_not_found (err: %v)", code, err)
	}
}

func TestSelectFailsWhenECRIsUnavailable(t *testing.T) {
	c := fakeBuilds(false, nil, map[string]int{})

	_, _, err := c.Select("avocado", []serviceFile{{Path: "service-x.json", Image: "broken"}}, false)
	if code := cher.Coerce(err).Code; code != "ecr_unavailable" {
		t.Fatalf("got %q, want ecr_unavailable (err: %v)", code, err)
	}
}

func TestSelectForceSelectsEverything(t *testing.T) {
	files := append(selectFiles, serviceFile{Path: "service-x.json", Image: "broken"}, serviceFile{Path: "service-y.json", Image: "broken"})
	lookups := map[string]int{}
	c := fakeBuilds(true, nil, lookups)

	selected, skipped, err := c.Select("avocado", files, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(selected) != len(files) || len(skipped) != 0 {
		t.Fatalf("got %d selected and %d skipped, want all %d selected", len(selected), len(skipped), len(files))
	}

	if lookups["broken"] != 1 {
		t.Fatalf("failed lookup retried %d times, want 1", lookups["broken"])
	}
}

func TestPlanChecksEveryEnvBeforeWriting(t *testing.T) {
	c := fakeBuilds(false, map[string]bool{"go-services": true}, map[string]int{})

	envs := map[string]struct{}{"avocado": {}, "basil": {}}
	listed := []string{}

	_, err := c.Plan(envs, nil, func(env string) ([]serviceFile, error) {
		listed = append(listed, env)

		if env == "basil" {
			return []serviceFile{{Path: "service-quote.json", Image: "js-quote"}}, nil
		}

		return []serviceFile{{Path: "service-sms.json", Image: "go-services"}}, nil
	}, true)

	if code := cher.Coerce(err).Code; code != "build_not_found" {
		t.Fatalf("got %q, want build_not_found (err: %v)", code, err)
	}

	if len(listed) != 2 {
		t.Fatalf("listed %v, want both envs", listed)
	}
}

func TestPlanFiltersByImage(t *testing.T) {
	lookups := map[string]int{}
	c := fakeBuilds(false, map[string]bool{"go-services": true, "js-services": true, "js-quote": true}, lookups)

	list := func(env string) ([]serviceFile, error) {
		return []serviceFile{
			{Path: "_base.json", Image: "go-services"},
			{Path: "service-sms.json", Image: "go-services"},
			{Path: "service-web.json", Image: "js-services"},
			{Path: "service-quote.json", Image: "js-quote"},
		}, nil
	}

	envs := map[string]struct{}{"avocado": {}}

	plans, err := c.Plan(envs, []string{"js-services", "js-quote"}, list, false)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := joinServiceNames(plans["avocado"].Selected), "quote, web"; got != want {
		t.Errorf("selected: got %q, want %q", got, want)
	}

	if got, want := joinServiceNames(plans["avocado"].Skipped), "_base, sms"; got != want {
		t.Errorf("skipped: got %q, want %q", got, want)
	}

	// only the requested images are looked up
	if lookups["go-services"] != 0 || lookups["js-services"] != 1 || lookups["js-quote"] != 1 {
		t.Errorf("unexpected lookups: %v", lookups)
	}

	c = fakeBuilds(false, nil, lookups)
	lookups["js-quote"] = 0

	_, err = c.Plan(envs, []string{"go-servics"}, list, false)
	if e := cher.Coerce(err); e.Code != "unknown_image" {
		t.Fatalf("got %q, want unknown_image (err: %v)", e.Code, err)
	}

	if lookups["js-quote"] != 0 {
		t.Errorf("ECR was checked before the mistyped image was caught: %v", lookups)
	}
}
