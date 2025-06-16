package helpers

import (
	"os"
	"regexp"
)

// dockerImageNameRegex regex pattern to extract docker_image_name from service config files
// Handles both JSON and YAML formats with optional quotes (single or double)
var dockerImageNameRegex = regexp.MustCompile(`["']?docker_image_name["']?\s*:\s*["']?([a-zA-Z\d_-]+)["']?`)

// ExtractDockerImageName extracts the docker_image_name from a service configuration file.
// Returns the docker image name, or an empty string if not found.
func ExtractDockerImageName(filePath string) (string, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	matches := dockerImageNameRegex.FindSubmatch(fileContents)
	if len(matches) != 2 {
		return "", nil // Not found, but not an error
	}

	return string(matches[1]), nil
}
