package app

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// serviceFile is a service config file in the config repo, together with the
// docker image it deploys.
type serviceFile struct {
	// Path is the absolute path of the file.
	Path string
	// Image is the docker image the file deploys, taken from the file itself or,
	// for a service that does not declare one, inherited from the env's _base
	// file. It is empty when neither defines an image.
	Image string
}

// Name is the file name, e.g. "service-sms.json".
func (f serviceFile) Name() string {
	return path.Base(f.Path)
}

// Service is the short name people type on the command line, e.g. "sms".
func (f serviceFile) Service() string {
	name := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))

	return strings.TrimPrefix(name, "service-")
}

// IsEnvFile reports whether the file holds env-wide values (_base.json,
// _base.yaml, ...) rather than the config of a single service.
func (f serviceFile) IsEnvFile() bool {
	return strings.HasPrefix(f.Name(), "_")
}

// definesOwnVersion reports whether the file pins its own commit or image tag,
// rather than inheriting them from the env's _base file.
func (f serviceFile) definesOwnVersion() bool {
	blob, err := os.ReadFile(f.Path)
	if err != nil {
		return false
	}

	if isYaml(f.Path) {
		return commitRegAddYaml.Match(blob) || imageTagRegAddYaml.Match(blob)
	}

	return commitRegAdd.Match(blob) || imageTagRegAdd.Match(blob)
}

type yamlImageConfig struct {
	Image struct {
		Repository string `yaml:"repository"`
	} `yaml:"image"`
}

func isYaml(p string) bool {
	return strings.HasSuffix(p, ".yaml") || strings.HasSuffix(p, ".yml")
}

// ownImage returns the docker image a config file declares itself: the
// top-level "docker_image_name" of a JSON file, or the last segment of
// "image.repository" in a Helm values YAML file. A missing file, or one that
// declares no image, returns "".
func ownImage(p string) (string, error) {
	blob, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	if isYaml(p) {
		var cfg yamlImageConfig
		if err := yaml.Unmarshal(blob, &cfg); err != nil {
			return "", fmt.Errorf("parse %s: %w", p, err)
		}

		return path.Base(cfg.Image.Repository), nil
	}

	matches := imageName.FindSubmatch(blob)
	if len(matches) != 2 {
		return "", nil
	}

	return string(matches[1]), nil
}

// resolveServiceFile works out which docker image a service config file
// deploys, falling back to the env's _base file the same way the deployment
// does when the service does not declare its own image.
func resolveServiceFile(p string) (serviceFile, error) {
	f := serviceFile{Path: p}

	image, err := ownImage(p)
	if err != nil {
		return f, err
	}

	if image == "" && !f.IsEnvFile() {
		base := path.Join(path.Dir(p), "_base"+path.Ext(p))

		image, err = ownImage(base)
		if err != nil {
			return f, err
		}
	}

	if image == "." {
		image = ""
	}

	f.Image = image

	return f, nil
}

// listServiceFiles returns every config file in an env's service directory.
// Env-wide files (_base.json, ...) are only included when includeEnvFiles is
// set. A missing directory returns no files.
func listServiceFiles(dir string, includeEnvFiles bool) ([]serviceFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Warnf("no service directory at %s, skipping", dir)
			return nil, nil
		}

		return nil, err
	}

	files := []serviceFile{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".json") && !isYaml(name) {
			continue
		}

		if !includeEnvFiles && strings.HasPrefix(name, "_") {
			continue
		}

		f, err := resolveServiceFile(path.Join(dir, name))
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

// namedServiceFiles returns the config files of the named services, preferring
// the JSON file and falling back to YAML.
func namedServiceFiles(repoPath, system, env string, services []string) ([]serviceFile, error) {
	files := []serviceFile{}

	for _, service := range services {
		p := paths.GetPathForService(repoPath, system, env, service)

		if _, err := os.Stat(p); err != nil {
			p = paths.GetYamlPathForService(repoPath, system, env, service)
			if _, err2 := os.Stat(p); err2 != nil {
				log.Warn(err)
				log.Warn(err2)
			}
		}

		f, err := resolveServiceFile(p)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

// warnInheritedVersions warns when updating an env-wide _base file would also
// move services that were skipped (for --image, or because their build is
// missing), as they take their commit from that _base file rather than pinning
// their own.
func warnInheritedVersions(env string, updatedBase serviceFile, skipped []serviceFile) {
	inheriting := []serviceFile{}

	for _, f := range skipped {
		if f.IsEnvFile() || path.Ext(f.Path) != path.Ext(updatedBase.Path) {
			continue
		}

		if !f.definesOwnVersion() {
			inheriting = append(inheriting, f)
		}
	}

	if len(inheriting) == 0 {
		return
	}

	log.Warnf("%s: %s was updated and these services inherit their commit from it, so they will move too even though they were skipped: %s", env, updatedBase.Name(), joinServiceNames(inheriting))
}

func joinServiceNames(files []serviceFile) string {
	if len(files) == 0 {
		return "(none)"
	}

	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, f.Service())
	}

	sort.Strings(names)

	return strings.Join(names, ", ")
}

// buildChecker decides which services can be deployed, by checking CI has
// pushed each service's image for the commit to ECR. Most services share an
// image (go-services, js-services, ...), so each image is looked up once per run.
type buildChecker struct {
	commit string
	branch string
	force  bool
	ready  map[string]bool
	errs   map[string]error
	find   func(image, commit, branch string) (bool, error)
}

func newBuildChecker(commit, branch string, force bool) *buildChecker {
	return &buildChecker{
		commit: commit,
		branch: branch,
		force:  force,
		ready:  map[string]bool{},
		errs:   map[string]error{},
		find:   findBuildInECR,
	}
}

// Ready reports whether the image has been built for the commit.
func (c *buildChecker) Ready(image string) (bool, error) {
	if ready, ok := c.ready[image]; ok {
		return ready, c.errs[image]
	}

	ready, err := c.find(image, c.commit, c.branch)
	if err != nil {
		c.ready[image] = false
		c.errs[image] = err

		if c.force {
			log.Warnf("could not check ECR for %s, deploying anyway as --force was given: %s", image, err)
		}

		return false, err
	}

	if ready {
		log.Infof("found %s:%s in ECR", image, imageTag(c.branch, c.commit))
	} else {
		log.Warnf("%s:%s is not in ECR yet", image, imageTag(c.branch, c.commit))
	}

	c.ready[image] = ready

	return ready, nil
}

// Select returns the files that can be deployed and those skipped because
// their image has not been built for the commit, or has no image cdep can find.
// Env-wide files without an image of their own (e.g. a _base.yaml holding only
// Helm settings) are always selected.
//
// When strict is set, as for services named on the command line, a missing
// build is an error rather than a skip. With force, every file is selected and
// problems are only warned about.
func (c *buildChecker) Select(env string, files []serviceFile, strict bool) (selected, skipped []serviceFile, err error) {
	notBuilt := map[string][]serviceFile{}
	unknown := []serviceFile{}

	for _, f := range files {
		if f.Image == "" {
			if f.IsEnvFile() || c.force {
				selected = append(selected, f)
				continue
			}

			if strict {
				log.Warnf("could not work out the docker image of %s, so cdep can't check it has been built", f.Name())
				selected = append(selected, f)
				continue
			}

			unknown = append(unknown, f)
			skipped = append(skipped, f)

			continue
		}

		ready, err := c.Ready(f.Image)
		if err != nil {
			if !c.force {
				return nil, nil, cher.New("ecr_unavailable", cher.M{"image": f.Image, "error": err.Error()})
			}
		}

		if ready || c.force {
			selected = append(selected, f)
			continue
		}

		if strict {
			return nil, nil, cher.New("build_not_found", cher.M{
				"path":  f.Path,
				"image": f.Image,
				"tag":   imageTag(c.branch, c.commit),
			})
		}

		notBuilt[f.Image] = append(notBuilt[f.Image], f)
		skipped = append(skipped, f)
	}

	images := make([]string, 0, len(notBuilt))
	for image := range notBuilt {
		images = append(images, image)
	}

	sort.Strings(images)

	for _, image := range images {
		log.Warnf("%s: skipping %d service(s) until %s:%s is built, run cdep again once CI has finished: %s",
			env, len(notBuilt[image]), image, imageTag(c.branch, c.commit), joinServiceNames(notBuilt[image]))
	}

	if len(unknown) > 0 {
		log.Warnf("%s: skipping %d service(s) with no docker image cdep can find (check docker_image_name or image.repository, or pass --force): %s",
			env, len(unknown), joinServiceNames(unknown))
	}

	return selected, skipped, nil
}

// servicePlan is the service files to update in one env, and those skipped.
type servicePlan struct {
	Selected []serviceFile
	Skipped  []serviceFile
}

// filterImages keeps the files that run one of images, and tells the user which
// services were left out. An empty list keeps every file.
func filterImages(env string, files []serviceFile, images []string) (kept, skipped []serviceFile) {
	if len(images) == 0 {
		return files, nil
	}

	for _, f := range files {
		if slices.Contains(images, f.Image) {
			kept = append(kept, f)
		} else {
			skipped = append(skipped, f)
		}
	}

	log.Infof("%s: updating %d service(s) running %s: %s", env, len(kept), strings.Join(images, ", "), joinServiceNames(kept))

	return kept, skipped
}

// Plan selects the files to update in every env before any are written, so a
// missing build, an ECR failure or a mistyped --image stops the run without
// leaving the config repo half edited.
//
// images limits the run to services running those docker images; empty means
// every service.
func (c *buildChecker) Plan(envs map[string]struct{}, images []string, list func(env string) ([]serviceFile, error), strict bool) (map[string]servicePlan, error) {
	names := make([]string, 0, len(envs))
	for env := range envs {
		names = append(names, env)
	}

	sort.Strings(names)

	listed := map[string][]serviceFile{}
	seenImages := map[string]struct{}{}

	for _, env := range names {
		files, err := list(env)
		if err != nil {
			return nil, fmt.Errorf("list services in %s: %w", env, err)
		}

		for _, f := range files {
			if f.Image != "" {
				seenImages[f.Image] = struct{}{}
			}
		}

		listed[env] = files
	}

	// Catch a mistyped --image before looking anything up in ECR.
	for _, image := range images {
		if _, ok := seenImages[image]; !ok {
			allowed := make([]string, 0, len(seenImages))
			for seen := range seenImages {
				allowed = append(allowed, seen)
			}

			sort.Strings(allowed)

			return nil, cher.New("unknown_image", cher.M{"image": image, "allowed": allowed})
		}
	}

	plans := map[string]servicePlan{}

	for _, env := range names {
		files, notRequested := filterImages(env, listed[env], images)

		selected, notBuilt, err := c.Select(env, files, strict)
		if err != nil {
			return nil, err
		}

		plans[env] = servicePlan{Selected: selected, Skipped: append(notRequested, notBuilt...)}
	}

	return plans, nil
}
