package helpers

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
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

// ShouldUpdateService determines if a service should be updated based on the filter flags
func ShouldUpdateService(filePath string, goOnly, jsOnly bool) (bool, error) {
	// If no filters are set, update all services
	if !goOnly && !jsOnly {
		return true, nil
	}

	// Extract docker_image_name from the service configuration file
	dockerImageName, err := ExtractDockerImageName(filePath)
	if err != nil {
		return false, err
	}

	if dockerImageName == "" {
		// If we can't find docker_image_name, skip this service with a warning
		log.Warnf("Could not find docker_image_name in %s, skipping", filePath)
		return false, nil
	}

	// Apply filters based on docker_image_name
	if goOnly {
		// Only update Go services (docker_image_name == "go_services" or "go-services")
		return dockerImageName == "go_services" || dockerImageName == "go-services", nil
	}

	if jsOnly {
		// Only update JS services (docker_image_name != "go_services" and != "go-services")
		return dockerImageName != "go_services" && dockerImageName != "go-services", nil
	}

	return true, nil
}
