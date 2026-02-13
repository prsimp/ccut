package data

import "testing"

func TestIsHex(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"0", true},
		{"0123456789abcdef", true},
		{"0123456789ABCDEF", true},
		{"abcdef00", true},
		{"b5a543cb", true},
		{"xyz", false},
		{"abcdefgh", false},
		{"12345g78", false},
		{"abcdef0 ", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isHex(tt.input)
			if got != tt.want {
				t.Errorf("isHex(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractProjectName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "standard src path",
			input: "-home-prs-src-arc-system",
			want:  "arc-system",
		},
		{
			name:  "standard src simple project",
			input: "-home-prs-src-ccut",
			want:  "ccut",
		},
		{
			name:  "projects dir",
			input: "-home-prs-projects-myapp",
			want:  "myapp",
		},
		{
			name:  "repos dir",
			input: "-home-prs-repos-tool",
			want:  "tool",
		},
		{
			name:  "code dir",
			input: "-home-prs-code-webapp",
			want:  "webapp",
		},
		{
			name:  "go dir",
			input: "-home-prs-go-mymod",
			want:  "mymod",
		},
		{
			name:  "worktree path",
			input: "-home-prs--claude-swarm-worktrees-12345678-1234-1234-1234-123456789abc-arc-system-b5a543cb-worktree-87654321-4321-4321-4321-cba987654321",
			want:  "arc-system",
		},
		{
			name:  "fallback with home prefix",
			input: "-home-prs-stuff-thing",
			want:  "stuff-thing",
		},
		{
			name:  "fallback with users prefix",
			input: "-Users-bob-stuff-thing",
			want:  "stuff-thing",
		},
		{
			name:  "no known markers single segment",
			input: "-something",
			want:  "-something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractProjectName(tt.input)
			if got != tt.want {
				t.Errorf("extractProjectName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
