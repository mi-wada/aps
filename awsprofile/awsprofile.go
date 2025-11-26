package awsprofile

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// List returns all AWS profiles from ~/.aws/config and ~/.aws/credentials
// sorted alphabetically.
func List() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".aws", "config")
	credentialsPath := filepath.Join(homeDir, ".aws", "credentials")

	profilesMap := make(map[string]bool)

	for profile := range parseProfilesFromFile(configPath, true) {
		profilesMap[profile] = true
	}

	for profile := range parseProfilesFromFile(credentialsPath, false) {
		profilesMap[profile] = true
	}

	profiles := make([]string, 0, len(profilesMap))
	for profile := range profilesMap {
		profiles = append(profiles, profile)
	}

	if len(profiles) == 0 {
		profiles = []string{"default"}
	}

	sort.Strings(profiles)

	return profiles, nil
}

// Current returns the current AWS profile from the AWS_PROFILE environment variable
// Returns "default" if not set
func Current() string {
	currentProfile := os.Getenv("AWS_PROFILE")
	if currentProfile == "" {
		return "default"
	}
	return currentProfile
}

// parseProfilesFromFile reads profile names from an AWS config file
func parseProfilesFromFile(filePath string, stripPrefix bool) map[string]bool {
	profiles := make(map[string]bool)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return profiles
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := strings.TrimSpace(line[1 : len(line)-1])
			if stripPrefix {
				profile = strings.TrimPrefix(profile, "profile ")
				profile = strings.TrimSpace(profile)
			}
			profiles[profile] = true
		}
	}

	return profiles
}
