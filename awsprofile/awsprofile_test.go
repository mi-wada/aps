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
		{
			name:     "returns the value of AWS_PROFILE for development",
			envValue: "development",
			want:     "development",
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

	t.Run("strips profile prefix from config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		configContent := `[default]
[profile production]
[profile development]`
		configPath := filepath.Join(awsDir, "config")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		want := []string{"default", "development", "production"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})

	t.Run("does not strip prefix from credentials file", func(t *testing.T) {
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		credentialsContent := `[default]
[production]
[development]`
		credentialsPath := filepath.Join(awsDir, "credentials")
		err = os.WriteFile(credentialsPath, []byte(credentialsContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create credentials file: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		want := []string{"default", "development", "production"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})

	t.Run("handles extra whitespace in profile names", func(t *testing.T) {
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		configContent := `  [default]
  [ profile production ]
[profile  development]`
		configPath := filepath.Join(awsDir, "config")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		want := []string{"default", "development", "production"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})

	t.Run("ignores non-profile lines", func(t *testing.T) {
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		configContent := `[default]
region = us-east-1
output = json
[profile production]
region = us-west-2
some random line
[profile development]`
		configPath := filepath.Join(awsDir, "config")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		want := []string{"default", "development", "production"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})

	t.Run("deduplicates profiles from config and credentials", func(t *testing.T) {
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		configContent := `[default]
[profile production]`
		configPath := filepath.Join(awsDir, "config")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		credentialsContent := `[default]
[production]
[development]`
		credentialsPath := filepath.Join(awsDir, "credentials")
		err = os.WriteFile(credentialsPath, []byte(credentialsContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create credentials file: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		// Should have each profile only once despite appearing in both files
		want := []string{"default", "development", "production"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})

	t.Run("returns profiles in alphabetical order regardless of file order", func(t *testing.T) {
		tmpDir := t.TempDir()
		awsDir := filepath.Join(tmpDir, ".aws")
		err := os.MkdirAll(awsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create .aws directory: %v", err)
		}

		// Profiles in reverse alphabetical order
		configContent := `[profile zulu]
[profile yankee]
[profile alpha]
[default]`
		configPath := filepath.Join(awsDir, "config")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		got, err := awsprofile.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		// Should be sorted alphabetically
		want := []string{"alpha", "default", "yankee", "zulu"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("List() = %v, want %v", got, want)
		}
	})
}
