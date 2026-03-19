package cher

import (
	"testing"
)

func TestSanitizeMeta(t *testing.T) {
	tests := []struct {
		name     string
		input    M
		expected M
	}{
		{
			name:     "nil metadata",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty metadata",
			input:    M{},
			expected: M{},
		},
		{
			name: "non-PII string",
			input: M{
				"file_id": "ABC123",
				"count":   42,
			},
			expected: M{
				"file_id": "ABC123",
				"count":   42,
			},
		},
		{
			name: "PII string",
			input: M{
				"email": "john.doe@example.com",
				"count": 42,
			},
			expected: M{
				"email": "[REDACTED]",
				"count": 42,
			},
		},
		{
			name: "mixed PII and non-PII",
			input: M{
				"file_id": "ABC123",
				"email":   "john.doe@example.com",
				"phone":   "+1234567890",
				"count":   42,
			},
			expected: M{
				"file_id": "ABC123",
				"email":   "[REDACTED]",
				"phone":   "[REDACTED]",
				"count":   42,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeMeta(tt.input)

			if tt.expected == nil && result != nil {
				t.Errorf("expected nil, got %v", result)
				return
			}

			if tt.expected != nil && result == nil {
				t.Errorf("expected %v, got nil", tt.expected)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(result))
				return
			}

			for key, expectedValue := range tt.expected {
				if resultValue, exists := result[key]; !exists {
					t.Errorf("missing key: %s", key)
				} else if resultValue != expectedValue {
					t.Errorf("key %s: expected %v, got %v", key, expectedValue, resultValue)
				}
			}
		})
	}
}

func TestNewWithPIISanitization(t *testing.T) {
	// Test that New() automatically sanitizes PII
	err := New("test_error", M{
		"file_id": "ABC123",
		"email":   "john.doe@example.com",
		"count":   42,
	})

	expectedMeta := M{
		"file_id": "ABC123",
		"email":   "[REDACTED]",
		"count":   42,
	}

	if len(err.Meta) != len(expectedMeta) {
		t.Errorf("expected %d keys, got %d", len(expectedMeta), len(err.Meta))
		return
	}

	for key, expectedValue := range expectedMeta {
		if resultValue, exists := err.Meta[key]; !exists {
			t.Errorf("missing key: %s", key)
		} else if resultValue != expectedValue {
			t.Errorf("key %s: expected %v, got %v", key, expectedValue, resultValue)
		}
	}
}

func TestErrorfWithPIISanitization(t *testing.T) {
	// Test that Errorf() automatically sanitizes PII
	err := Errorf("test_error", M{
		"file_id": "ABC123",
		"email":   "john.doe@example.com",
	}, "test message")

	expectedMeta := M{
		"file_id": "ABC123",
		"email":   "[REDACTED]",
		"message": "test message",
	}

	if len(err.Meta) != len(expectedMeta) {
		t.Errorf("expected %d keys, got %d", len(expectedMeta), len(err.Meta))
		return
	}

	for key, expectedValue := range expectedMeta {
		if resultValue, exists := err.Meta[key]; !exists {
			t.Errorf("missing key: %s", key)
		} else if resultValue != expectedValue {
			t.Errorf("key %s: expected %v, got %v", key, expectedValue, resultValue)
		}
	}
}
