package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"strconv"
)

// Conventional Commit regex to extract type and look for BREAKING CHANGE
var commitRegex = regexp.MustCompile(`^(feat|fix|chore|docs|style|refactor|perf|test|build|ci)(\([\w\-]+\))?(!?): (.+)`)

// SemVerBump represents the highest version bump found in commits.
type SemVerBump int

const (
	BumpNone SemVerBump = iota
	BumpPatch
	BumpMinor
	BumpMajor
)

// getCurrentVersion reads the latest Git tag (vX.Y.Z) and returns the version components.
func getCurrentVersion() (int, int, int, error) {
	// Find the latest annotated or lightweight tag matching vX.Y.Z
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", "--match", "v[0-9]*.[0-9]*.[0-9]*")
	tag, err := cmd.Output()

	if err != nil {
		// If no tags exist or no commits, start at v0.0.0
		errorOutput := err.Error()
		if strings.Contains(errorOutput, "No tags found") || strings.Contains(errorOutput, "No names found") {
			fmt.Println("No SemVer tags found or no commits. Starting from v0.0.0.")
			return 0, 0, 0, nil
		}
		return 0, 0, 0, fmt.Errorf("error getting latest tag: %w", err)
	}

	// Clean and parse the tag (e.g., "v1.2.3\n" -> "1.2.3")
	versionStr := strings.TrimPrefix(strings.TrimSpace(string(tag)), "v")
	parts := strings.Split(versionStr, ".")

	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid SemVer tag format: %s", versionStr)
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	fmt.Printf("Last released version: v%d.%d.%d\n", major, minor, patch)
	return major, minor, patch, nil
}

// getCommitsSinceLastTag fetches all commit messages since the last Git tag.
func getCommitsSinceLastTag() ([]string, error) {
	// Get the last tag again to determine the commit range
	cmdTag := exec.Command("git", "describe", "--tags", "--abbrev=0")
	lastTag, err := cmdTag.Output()

	var cmd *exec.Cmd
	if err != nil || len(lastTag) == 0 {
		// If no tags, get all commits from the beginning
		cmd = exec.Command("git", "log", "--pretty=format:%s%n%b")
		fmt.Println("Analyzing all commits (no previous tag found).")
	} else {
		// Get commits since the last tag, including subject and body
		tag := strings.TrimSpace(string(lastTag))
		cmd = exec.Command("git", "log", fmt.Sprintf("%s..HEAD", tag), "--pretty=format:%s%n%b")
		fmt.Printf("Analyzing commits since %s...\n", tag)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git log: %w", err)
	}

	// Split the output into individual commit messages (subject + body)
	commits := strings.Split(strings.TrimSpace(string(output)), "\n\n")
	return commits, nil
}

// analyzeCommits determines the required SemVer bump based on Conventional Commits.
func analyzeCommits(commits []string) SemVerBump {
	requiredBump := BumpNone

	for _, commit := range commits {
		// Check for explicit 'BREAKING CHANGE' in the commit body (case-insensitive)
		if strings.Contains(strings.ToUpper(commit), "BREAKING CHANGE") {
			return BumpMajor // Major is the highest possible bump, so we can return early
		}

		// Check the commit header for '!' (breaking change syntax) and type
		match := commitRegex.FindStringSubmatch(commit)
		if len(match) > 0 {
			commitType := match[1]
			isBreaking := match[3] == "!"

			if isBreaking {
				return BumpMajor
			}

			if commitType == "feat" {
				if requiredBump < BumpMinor {
					requiredBump = BumpMinor
				}
			} else if commitType == "fix" {
				if requiredBump < BumpPatch {
					requiredBump = BumpPatch
				}
			}
			// Other types (chore, docs, test, etc.) default to BumpNone unless configured otherwise
		}
	}

	return requiredBump
}

func main() {
	fmt.Println("--- SemVer Version Bumper (Conventional Commits) ---")

	currentMajor, currentMinor, currentPatch, err := getCurrentVersion()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	commits, err := getCommitsSinceLastTag()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
    
    // If no new commits, just return the current version
	if len(commits) == 0 || (len(commits) == 1 && commits[0] == "") {
        fmt.Println("\nNo new conventional commits since last tag.")
        fmt.Printf("Next suggested version: v%d.%d.%d (No Change)\n", currentMajor, currentMinor, currentPatch)
        return
    }
    
	bump := analyzeCommits(commits)

	newMajor, newMinor, newPatch := currentMajor, currentMinor, currentPatch

	switch bump {
	case BumpMajor:
		newMajor++
		newMinor = 0
		newPatch = 0
		fmt.Println("\nDetermined BUMP: MAJOR (Found BREAKING CHANGE or feat!/fix! commit)")
	case BumpMinor:
		newMinor++
		newPatch = 0
		fmt.Println("\nDetermined BUMP: MINOR (Found 'feat:' commit)")
	case BumpPatch:
		newPatch++
		fmt.Println("\nDetermined BUMP: PATCH (Found 'fix:' commit)")
	case BumpNone:
		fmt.Println("\nDetermined BUMP: NONE (Only 'chore:', 'docs:', etc. commits found)")
	}

	// Note: We avoid bumping 0.Y.Z to 1.0.0 automatically here for simplicity,
	// but production tools handle this based on first MAJOR release rule.
	
	if bump == BumpNone {
		fmt.Printf("Next suggested version: v%d.%d.%d (No Change)\n", newMajor, newMinor, newPatch)
	} else {
		fmt.Printf("Next suggested version: v%d.%d.%d\n", newMajor, newMinor, newPatch)
	}
}

// Note: This Go program relies on the 'git' command-line tool being installed and 
// the code being run inside a Git repository.
