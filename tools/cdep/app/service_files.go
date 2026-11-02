package app

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/cuvva/cuvva-public-go/tools/cdep"
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

// Language reports which pipeline builds the file's image, or "" when the image
// is unknown.
func (f serviceFile) Language() string {
	if f.Image == "" {
		return ""
	}

	return cdep.LanguageForImage(f.Image)
}

// Matches reports whether the file should be updated for the given language
// filter. An empty filter matches every file.
func (f serviceFile) Matches(only string) bool {
	return only == "" || f.Language() == only
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

// filterServiceFiles splits files into those matching the language filter and
// those skipped, and tells the user which is which.
func filterServiceFiles(env, only string, files []serviceFile) (selected, skipped []serviceFile) {
	for _, f := range files {
		if f.Matches(only) {
			selected = append(selected, f)
		} else {
			skipped = append(skipped, f)
		}
	}

	if only == "" {
		return selected, skipped
	}

	log.Infof("%s: updating %d %s service(s): %s", env, len(selected), only, joinServiceNames(selected))

	if len(skipped) > 0 {
		log.Infof("%s: skipping %d service(s) not built by %s: %s", env, len(skipped), only, joinServiceNames(skipped))
	}

	for _, f := range skipped {
		if f.Image == "" {
			log.Warnf("%s: could not work out the docker image of %s, so it was skipped; check its docker_image_name or image.repository", env, f.Name())
		}
	}

	return selected, skipped
}

// warnInheritedVersions warns when updating an env-wide _base file would also
// move services the language filter skipped, because they take their commit
// from that _base file rather than pinning their own.
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

	log.Warnf("%s: %s was updated and these skipped services inherit their commit from it, so they will move too: %s", env, updatedBase.Name(), joinServiceNames(inheriting))
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

// ecrChecker checks each docker image has a build for the commit being
// deployed. Most services share an image (go-services, js-services, ...), so
// each image is only looked up once per run.
type ecrChecker struct {
	commit  string
	branch  string
	checked map[string]struct{}
	find    func(image, commit, branch string) error
}

func newECRChecker(commit, branch string) *ecrChecker {
	return &ecrChecker{
		commit:  commit,
		branch:  branch,
		checked: map[string]struct{}{},
		find:    findBuildInECR,
	}
}

func (c *ecrChecker) Check(f serviceFile) {
	if f.Image == "" {
		if !f.IsEnvFile() {
			log.Warnf("could not work out the docker image of %s, skipping the ECR check", f.Name())
		}

		return
	}

	if _, ok := c.checked[f.Image]; ok {
		return
	}

	c.checked[f.Image] = struct{}{}

	if err := c.find(f.Image, c.commit, c.branch); err != nil {
		log.Warnf("ecr: could not check for a %s build of %s: %s", f.Image, c.commit, err)
	}
}
