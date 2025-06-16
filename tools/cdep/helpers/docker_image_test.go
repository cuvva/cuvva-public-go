package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractDockerImageName(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "cdep_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		fileContent    string
		expectedImage  string
		expectedError  bool
		description    string
	}{
		{
			name: "go_services_json",
			fileContent: `{
	"docker_image_name": "go_services",
	"branch": "master",
	"commit": "abc123"
}`,
			expectedImage: "go_services",
			expectedError: false,
			description:   "Standard Go service JSON config (underscore)",
		},
		{
			name: "go_services_hyphen_json",
			fileContent: `{
	"docker_image_name": "go-services",
	"branch": "master",
	"commit": "abc123"
}`,
			expectedImage: "go-services",
			expectedError: false,
			description:   "Standard Go service JSON config (hyphen)",
		},
		{
			name: "js_service_json",
			fileContent: `{
	"docker_image_name": "js_frontend",
	"branch": "master",
	"commit": "def456"
}`,
			expectedImage: "js_frontend",
			expectedError: false,
			description:   "Standard JS service JSON config",
		},
		{
			name: "quoted_image_name",
			fileContent: `{
	"docker_image_name": "my-service-name",
	"other_field": "value"
}`,
			expectedImage: "my-service-name",
			expectedError: false,
			description:   "Image name with hyphens in quotes",
		},
		{
			name: "unquoted_image_name",
			fileContent: `{
	docker_image_name: my_service_123,
	other_field: "value"
}`,
			expectedImage: "my_service_123",
			expectedError: false,
			description:   "Unquoted image name (YAML-style)",
		},
		{
			name: "image_name_with_spaces",
			fileContent: `{
	"docker_image_name"  :  "spaced_service"  ,
	"other_field": "value"
}`,
			expectedImage: "spaced_service",
			expectedError: false,
			description:   "Image name with extra spaces around colons",
		},
		{
			name: "yaml_format",
			fileContent: `docker_image_name: yaml_service
branch: master
commit: xyz789`,
			expectedImage: "yaml_service",
			expectedError: false,
			description:   "YAML format configuration",
		},
		{
			name: "yaml_quoted",
			fileContent: `docker_image_name: "yaml_quoted_service"
branch: "master"`,
			expectedImage: "yaml_quoted_service",
			expectedError: false,
			description:   "YAML format with quoted values",
		},
		{
			name: "no_docker_image_name",
			fileContent: `{
	"branch": "master",
	"commit": "abc123",
	"other_field": "value"
}`,
			expectedImage: "",
			expectedError: false,
			description:   "Config without docker_image_name field",
		},
		{
			name: "empty_docker_image_name",
			fileContent: `{
	"docker_image_name": "",
	"branch": "master"
}`,
			expectedImage: "",
			expectedError: false,
			description:   "Config with empty docker_image_name",
		},
		{
			name: "multiple_occurrences",
			fileContent: `{
	"docker_image_name": "first_service",
	"description": "This service uses docker_image_name for deployment",
	"backup_docker_image_name": "backup_service"
}`,
			expectedImage: "first_service",
			expectedError: false,
			description:   "Multiple occurrences - should match first",
		},
		{
			name: "malformed_json",
			fileContent: `{
	"docker_image_name": "service_name"
	"missing_comma": "value"
}`,
			expectedImage: "service_name",
			expectedError: false,
			description:   "Malformed JSON but valid docker_image_name",
		},
		{
			name: "nested_structure",
			fileContent: `{
	"config": {
		"docker_image_name": "nested_service"
	},
	"docker_image_name": "top_level_service"
}`,
			expectedImage: "nested_service",
			expectedError: false,
			description:   "Nested structure - should match first occurrence",
		},
		{
			name: "special_characters",
			fileContent: `{
	"docker_image_name": "service-with_special123",
	"branch": "master"
}`,
			expectedImage: "service-with_special123",
			expectedError: false,
			description:   "Image name with allowed special characters",
		},
		{
			name: "empty_file",
			fileContent: "",
			expectedImage: "",
			expectedError: false,
			description:   "Empty file",
		},
		{
			name: "only_whitespace",
			fileContent: "   \n\t  \n  ",
			expectedImage: "",
			expectedError: false,
			description:   "File with only whitespace",
		},
		{
			name: "base_config_file",
			fileContent: `{
	"docker_image_name": "base_template_service",
	"branch": "master",
	"commit": "template123"
}`,
			expectedImage: "base_template_service",
			expectedError: false,
			description:   "Base configuration file (should extract image name but be filtered out at higher level)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".json")
			err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test the function
			result, err := ExtractDockerImageName(testFile)

			// Check error expectation
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check result
			if result != tt.expectedImage {
				t.Errorf("Expected image name '%s', got '%s'", tt.expectedImage, result)
			}

			t.Logf("✓ %s: Expected '%s', got '%s'", tt.description, tt.expectedImage, result)
		})
	}
}

func TestExtractDockerImageName_FileErrors(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		description string
	}{
		{
			name:        "non_existent_file",
			filePath:    "/tmp/non_existent_file_12345.json",
			description: "Non-existent file should return error",
		},
		{
			name:        "empty_path",
			filePath:    "",
			description: "Empty file path should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractDockerImageName(tt.filePath)

			// Should return error
			if err == nil {
				t.Errorf("Expected error for %s but got none", tt.description)
			}

			// Result should be empty on error
			if result != "" {
				t.Errorf("Expected empty result on error, got '%s'", result)
			}

			t.Logf("✓ %s: Got expected error: %v", tt.description, err)
		})
	}
}

func TestExtractDockerImageName_EdgeCases(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cdep_edge_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		fileContent   string
		expectedImage string
		description   string
	}{
		{
			name: "very_long_image_name",
			fileContent: `{
	"docker_image_name": "very_long_service_name_with_many_characters_and_numbers_123456789",
	"branch": "master"
}`,
			expectedImage: "very_long_service_name_with_many_characters_and_numbers_123456789",
			description:   "Very long image name",
		},
		{
			name: "single_character_name",
			fileContent: `{
	"docker_image_name": "a",
	"branch": "master"
}`,
			expectedImage: "a",
			description:   "Single character image name",
		},
		{
			name: "numbers_only",
			fileContent: `{
	"docker_image_name": "123456",
	"branch": "master"
}`,
			expectedImage: "123456",
			description:   "Numbers only image name",
		},
		{
			name: "mixed_quotes",
			fileContent: `{
	"docker_image_name": 'mixed_quotes_service',
	"branch": "master"
}`,
			expectedImage: "mixed_quotes_service",
			description:   "Mixed quote styles around values",
		},
		{
			name: "case_sensitivity",
			fileContent: `{
	"DOCKER_IMAGE_NAME": "uppercase_field",
	"docker_image_name": "lowercase_field"
}`,
			expectedImage: "lowercase_field",
			description:   "Case sensitivity - should match exact field name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.name+".json")
			err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := ExtractDockerImageName(testFile)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expectedImage {
				t.Errorf("Expected image name '%s', got '%s'", tt.expectedImage, result)
			}

			t.Logf("✓ %s: Expected '%s', got '%s'", tt.description, tt.expectedImage, result)
		})
	}
}

// Benchmark test to ensure the function performs well
func BenchmarkExtractDockerImageName(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "cdep_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testContent := `{
	"docker_image_name": "benchmark_service",
	"branch": "master",
	"commit": "abc123",
	"other_field": "value",
	"nested": {
		"field": "value"
	}
}`
	testFile := filepath.Join(tempDir, "benchmark.json")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractDockerImageName(testFile)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
