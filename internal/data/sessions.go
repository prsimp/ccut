package data

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DiscoverSessions finds all session JSONL files and extracts metadata.
// Returns at most maxSessions, sorted by most recent first.
func DiscoverSessions(maxSessions int) ([]SessionInfo, error) {
	projectsDir := filepath.Join(ClaudeDir(), "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, err
	}

	var sessions []SessionInfo
	for _, projEntry := range entries {
		if !projEntry.IsDir() {
			continue
		}
		projPath := filepath.Join(projectsDir, projEntry.Name())
		files, err := os.ReadDir(projPath)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".jsonl") {
				continue
			}
			sessionID := strings.TrimSuffix(f.Name(), ".jsonl")
			info := parseSessionMeta(filepath.Join(projPath, f.Name()), sessionID, projEntry.Name())
			if info.MessageCount > 0 {
				sessions = append(sessions, info)
			}
		}
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastTime > sessions[j].LastTime
	})

	if len(sessions) > maxSessions {
		sessions = sessions[:maxSessions]
	}
	return sessions, nil
}

func parseSessionMeta(path, sessionID, projectDir string) SessionInfo {
	info := SessionInfo{
		SessionID:   sessionID,
		ProjectDir:  projectDir,
		ProjectName: extractProjectName(projectDir),
	}

	f, err := os.Open(path)
	if err != nil {
		return info
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		var line SessionLine
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			continue
		}

		if line.Type != "user" && line.Type != "assistant" {
			continue
		}

		if line.GitBranch != "" && info.GitBranch == "" {
			info.GitBranch = line.GitBranch
		}
		if line.Timestamp != "" {
			if info.FirstTime == "" {
				info.FirstTime = line.Timestamp
			}
			info.LastTime = line.Timestamp
		}

		info.MessageCount++

		if line.Type == "assistant" && len(line.Message) > 0 {
			var msg SessionMessage
			if err := json.Unmarshal(line.Message, &msg); err == nil {
				if msg.Model != "" {
					info.Model = msg.Model
				}
				if msg.Usage != nil {
					info.OutputTokens += msg.Usage.OutputTokens
					info.TokenUsage.InputTokens += msg.Usage.InputTokens
					info.TokenUsage.OutputTokens += msg.Usage.OutputTokens
					info.TokenUsage.CacheReadInputTokens += msg.Usage.CacheReadInputTokens
					info.TokenUsage.CacheCreationInputTokens += msg.Usage.CacheCreationInputTokens
				}
			}
		}
	}
	return info
}

func extractProjectName(dir string) string {
	// Reconstruct the filesystem path from the encoded dir name.
	// Claude encodes paths by replacing '/' with '-', so
	// "-home-prs-src-arc-system" was originally "/home/prs/src/arc-system".
	// For worktree dirs like:
	//   -home-prs--claude-swarm-worktrees-UUID-arc-system-HASH-worktree-UUID
	// we want to extract "arc-system".

	// If it contains "worktree" markers, try to extract the project name
	// between the worktrees-UUID- prefix and the -HASH-worktree suffix.
	if idx := strings.Index(dir, "-worktrees-"); idx >= 0 {
		// Skip past "-worktrees-" + UUID (36 chars) + "-"
		after := dir[idx+len("-worktrees-"):]
		// Skip UUID: 8-4-4-4-12 = 36 chars
		if len(after) > 37 && after[36] == '-' {
			after = after[37:]
		}
		// after is now like "arc-system-b5a543cb-worktree-UUID"
		// Find the project name by removing the trailing hash-worktree-UUID
		if wi := strings.Index(after, "-worktree-"); wi >= 0 {
			after = after[:wi]
		}
		// Remove trailing short hash (8 hex chars)
		if len(after) > 9 && after[len(after)-9] == '-' && isHex(after[len(after)-8:]) {
			after = after[:len(after)-9]
		}
		// If the worktree name is generic (e.g. "agent"), extract the real
		// project name from the path before the worktree marker.
		// e.g. "-Users-rob-dev-kajabi-products--claude-worktrees-agent-HASH"
		//   prefix before "--claude-worktrees" is "-Users-rob-dev-kajabi-products"
		if after == "" || after == "agent" || after == "swarm" {
			prefix := dir[:idx]
			// Strip the "--claude" suffix from the prefix if present
			if ci := strings.Index(prefix, "--claude"); ci >= 0 {
				prefix = prefix[:ci]
			}
			return extractProjectName(prefix)
		}
		return after
	}

	// Standard path: "-home-prs-src-arc-system"
	// Find the last known parent directory marker and take everything after it.
	markers := []string{"-src-", "-projects-", "-repos-", "-code-", "-go-"}
	for _, m := range markers {
		if idx := strings.LastIndex(dir, m); idx >= 0 {
			name := dir[idx+len(m):]
			if name != "" {
				return name
			}
		}
	}

	// Fallback: take the last path-like segment.
	// The dir starts with "-", split on "-" and rejoin the tail.
	parts := strings.Split(strings.TrimLeft(dir, "-"), "-")
	if len(parts) >= 2 {
		// Skip username-like segments (short, common names)
		skip := map[string]bool{"home": true, "users": true, "var": true}
		start := 0
		for i, p := range parts {
			if skip[strings.ToLower(p)] {
				start = i + 2 // skip the username too
			}
		}
		if start < len(parts) {
			return strings.Join(parts[start:], "-")
		}
	}
	return dir
}

// ComputeDailyCosts scans all session JSONL files and aggregates estimated
// cost by date. Returns a map of "YYYY-MM-DD" → cost in USD.
func ComputeDailyCosts() map[string]float64 {
	costs := make(map[string]float64)
	projectsDir := filepath.Join(ClaudeDir(), "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return costs
	}

	for _, projEntry := range entries {
		if !projEntry.IsDir() {
			continue
		}
		projPath := filepath.Join(projectsDir, projEntry.Name())
		files, err := os.ReadDir(projPath)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".jsonl") {
				continue
			}
			aggregateSessionCosts(filepath.Join(projPath, f.Name()), costs)
		}
	}
	return costs
}

func aggregateSessionCosts(path string, costs map[string]float64) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		var line SessionLine
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			continue
		}
		if line.Type != "assistant" || len(line.Message) == 0 {
			continue
		}
		var msg SessionMessage
		if err := json.Unmarshal(line.Message, &msg); err != nil || msg.Usage == nil {
			continue
		}
		if line.Timestamp == "" {
			continue
		}
		// Extract date from ISO timestamp
		date := line.Timestamp
		if len(date) >= 10 {
			date = date[:10]
		}
		usage := ModelUsage{
			InputTokens:              msg.Usage.InputTokens,
			OutputTokens:             msg.Usage.OutputTokens,
			CacheReadInputTokens:     msg.Usage.CacheReadInputTokens,
			CacheCreationInputTokens: msg.Usage.CacheCreationInputTokens,
		}
		costs[date] += CalculateCost(msg.Model, usage)
	}
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return len(s) > 0
}
