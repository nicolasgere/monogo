package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAffectedRootDirectories(compareBranch string, dir string) ([]string, error) {
	// Get the list of changed files
	cmd := exec.Command("git", "diff", "--name-only", compareBranch)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing git command: %v", err)
	}

	// Split the output into individual file paths
	changedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Create a map to store unique root directories
	affectedRootDirs := make(map[string]bool)

	// Extract root directories from file paths
	for _, file := range changedFiles {
		rootDir := extractRootDirectory(file)
		if rootDir != "" {
			affectedRootDirs[rootDir] = true
		}
	}

	// Convert map keys to slice
	result := make([]string, 0, len(affectedRootDirs))
	for dir := range affectedRootDirs {
		result = append(result, dir)
	}

	return result, nil
}

func extractRootDirectory(filePath string) string {
	parts := strings.SplitN(filePath, string(filepath.Separator), 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
