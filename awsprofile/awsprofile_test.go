package awsprofile_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mi-wada/aps/awsprofile"
)

func TestCurrent(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns default when AWS_PROFILE is not set",
			envValue: "",
			want:     "default",
		},
		{
			name:     "returns the value of AWS_PROFILE when set",
			envValue: "production",
			want:     "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalValue := os.Getenv("AWS_PROFILE")
			defer os.Setenv("AWS_PROFILE", originalValue)

			// Set test value
			if tt.envValue == "" {
				os.Unsetenv("AWS_PROFILE")
			} else {
				os.Setenv("AWS_PROFILE", tt.envValue)
			}

			got := awsprofile.Current()
			if got != tt.want {
				t.Errorf("Current() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	t.Run("returns profiles from config and credentials sorted alphabetically", func(t *testing.T) {
		// Create temporary home directory
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		// Create test config file
		configContent := `[default]
region = us-east-1
[profile staging]
region = us-west-1
[profile production]
region = eu-west-1`
		configPath := filepath.Join(awsDir, "config")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		// Create test credentials file
		credentialsContent := `[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
[development]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE2`
		credentialsPath := filepath.Join(awsDir, "credentials")
		err = os.WriteFile(credentialsPath, []byte(credentialsContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create credentials file: %v", err)
		}

		// Save original HOME and restore after test
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		// Test List()
		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		// Expected profiles sorted alphabetically
		want := []string{"default", "development", "production", "staging"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})

	t.Run("returns default when no config files exist", func(t *testing.T) {
		// Create temporary home directory without .aws folder
		tmpDir := t.TempDir()

		// Save original HOME and restore after test
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		// Test List()
		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		// Should return default profile when no config files exist
		want := []string{"default"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})
}
