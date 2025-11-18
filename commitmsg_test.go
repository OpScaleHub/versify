package main

import "testing"

func TestAnalyzeCommits_SpecialCharsAndEmojis(t *testing.T) {
	tests := []struct {
		name    string
		commits []string
		want    SemVerBump
	}{
		{"feat simple", []string{"feat: add new feature"}, BumpMinor},
		{"fix simple", []string{"fix: fix a bug"}, BumpPatch},
		{"chore only", []string{"chore: housekeeping"}, BumpNone},
		{"breaking in body", []string{"chore: update\n\nBREAKING CHANGE: alters API"}, BumpMajor},
		{"feat with emoji", []string{"feat: add ✨ new feature"}, BumpMinor},
		{"fix with emoji and special chars", []string{"fix(parser): handle 🐛 & ☠️"}, BumpPatch},
		{"exclamation breaking", []string{"feat!: big change that breaks things"}, BumpMajor},
		{"lowercase breaking change", []string{"refactor: reshape\n\nbreaking change: lower-case marker"}, BumpMajor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzeCommits(tt.commits)
			if got != tt.want {
				t.Fatalf("analyzeCommits() = %v, want %v for commits: %#v", got, tt.want, tt.commits)
			}
		})
	}
}
