//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"
	"time"
)

// formatTimeAgo returns a human-readable relative time string.
//
// Examples: "just now", "5 minutes ago", "2 hours ago", "3 days ago",
// or a formatted date for times older than a week.
//
// Parameters:
//   - t: The time to format relative to now
//
// Returns:
//   - string: Human-readable relative time
func formatTimeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}

// formatNumber returns a number with thousand separators.
//
// Examples: 500 -> "500", 1500 -> "1,500", 12345 -> "12,345"
//
// Parameters:
//   - n: The number to format
//
// Returns:
//   - string: Formatted number with commas
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}

// formatBytes returns a human-readable byte size string.
//
// Uses binary units (1024-based): B, KB, MB, GB, etc.
//
// Examples: 500 -> "500 B", 1536 -> "1.5 KB", 1048576 -> "1.0 MB"
//
// Parameters:
//   - b: The byte count to format
//
// Returns:
//   - string: Human-readable size with unit
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
