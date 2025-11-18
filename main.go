package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// getCurrentVersion reads the latest Git tag (e.g., vX.Y.Z) or uses a baseline, and returns the version components.
func getCurrentVersion(prefix, baseline string) (int, int, int, error) {
	var versionStr string
	if baseline != "" {
		fmt.Fprintf(os.Stderr, "Using baseline version: %s\n", baseline)
		versionStr = strings.TrimPrefix(baseline, prefix)
	} else {
		// Find the latest tag on the current branch
		tag, err := getLatestTagOnBranch(prefix)
		if err != nil {
			errorOutput := err.Error()
			if strings.Contains(errorOutput, "No tags found") || strings.Contains(errorOutput, "No names found") {
				fmt.Fprintf(os.Stderr, "No SemVer tags with prefix '%s' found. Starting from %s0.0.0.\n", prefix, prefix)
				return 0, 0, 0, nil
			}
			return 0, 0, 0, fmt.Errorf("error getting latest tag: %w", err)
		}
		versionStr = strings.TrimPrefix(tag, prefix)
	}

	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid SemVer tag format: %s", versionStr)
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	fmt.Fprintf(os.Stderr, "Last released version: %s%d.%d.%d\n", prefix, major, minor, patch)
	return major, minor, patch, nil
}

// getLatestTagOnBranch finds the latest tag on the current branch matching a prefix.
func getLatestTagOnBranch(prefix string) (string, error) {
	// Get the current branch name
	cmdBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchName, err := cmdBranch.Output()
	if err != nil {
		return "", fmt.Errorf("error getting current branch name: %w", err)
	}
	branch := strings.TrimSpace(string(branchName))

	// Find the latest tag on the current branch
	matchPattern := fmt.Sprintf("%s[0-9]*.[0-9]*.[0-9]*", prefix)
	cmdTag := exec.Command("git", "describe", "--tags", "--abbrev=0", "--match", matchPattern, branch)
	tag, err := cmdTag.Output()
	if err != nil {
		return "", err // Return error to be handled by the caller
	}

	return strings.TrimSpace(string(tag)), nil
}

// getCommitsSinceLastTag fetches all commit messages since the last Git tag matching the prefix.
func getCommitsSinceLastTag(prefix string) ([]string, error) {
	lastTag, err := getLatestTagOnBranch(prefix)

	var cmd *exec.Cmd
	if err != nil || lastTag == "" {
		cmd = exec.Command("git", "log", "--pretty=format:%s%n%b")
		fmt.Fprintln(os.Stderr, "Analyzing all commits (no previous tag found).")
	} else {
		cmd = exec.Command("git", "log", fmt.Sprintf("%s..HEAD", lastTag), "--pretty=format:%s%n%b")
		fmt.Fprintf(os.Stderr, "Analyzing commits since %s...\n", lastTag)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git log: %w", err)
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n\n")
	return commits, nil
}

// analyzeCommits determines the required SemVer bump based on Conventional Commits.
func analyzeCommits(commits []string) SemVerBump {
	requiredBump := BumpNone

	for _, commit := range commits {
		if strings.Contains(strings.ToUpper(commit), "BREAKING CHANGE") {
			return BumpMajor
		}

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
		}
	}

	return requiredBump
}

func main() {
	prefix := flag.String("prefix", "v", "The prefix for the version tag (e.g., 'v', 'k8s')")
	baseline := flag.String("baseline", "", "The baseline version to use instead of the latest tag (e.g., '1.2.3')")
	addSuffix := flag.Bool("add-suffix", false, "Always add a suffix to the version (e.g., for pre-releases or builds)")
	suffixFormat := flag.String("suffix-format", "short-hash", "The format for the suffix (e.g., 'short-hash', 'datetime')")
	flag.Parse()

	fmt.Fprintln(os.Stderr, "--- SemVer Version Bumper (Conventional Commits) ---")

	currentMajor, currentMinor, currentPatch, err := getCurrentVersion(*prefix, *baseline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	commits, err := getCommitsSinceLastTag(*prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if len(commits) == 0 || (len(commits) == 1 && commits[0] == "") {
		fmt.Fprintln(os.Stderr, "\nNo new conventional commits since last tag.")
		versionString := fmt.Sprintf("%s%d.%d.%d", *prefix, currentMajor, currentMinor, currentPatch)
		fmt.Println(versionString)
		return
	}

	bump := analyzeCommits(commits)

	newMajor, newMinor, newPatch := currentMajor, currentMinor, currentPatch

	switch bump {
	case BumpMajor:
		newMajor++
		newMinor = 0
		newPatch = 0
		fmt.Fprintln(os.Stderr, "\nDetermined BUMP: MAJOR (Found BREAKING CHANGE or feat!/fix! commit)")
	case BumpMinor:
		newMinor++
		newPatch = 0
		fmt.Fprintln(os.Stderr, "\nDetermined BUMP: MINOR (Found 'feat:' commit)")
	case BumpPatch:
		newPatch++
		fmt.Fprintln(os.Stderr, "\nDetermined BUMP: PATCH (Found 'fix:' commit)")
	case BumpNone:
		fmt.Fprintln(os.Stderr, "\nDetermined BUMP: NONE (Only 'chore:', 'docs:', etc. commits found)")
	}

	versionString := fmt.Sprintf("%s%d.%d.%d", *prefix, newMajor, newMinor, newPatch)

	if *addSuffix {
		var suffix string
		if *suffixFormat == "datetime" {
			suffix = time.Now().UTC().Format("20060102150405")
		} else { // default to short-hash
			cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
			shortHash, err := cmd.Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting short commit hash: %v\n", err)
			} else {
				suffix = strings.TrimSpace(string(shortHash))
			}
		}
		if suffix != "" {
			versionString = fmt.Sprintf("%s-%s", versionString, suffix)
		}
		fmt.Fprintln(os.Stderr, "Suffix added.")
	}

	if bump == BumpNone && !*addSuffix {
		fmt.Fprintln(os.Stderr, "No change detected.")
	}

	fmt.Println(versionString)
}

// Note: This Go program relies on the 'git' command-line tool being installed and
// the code being run inside a Git repository.
