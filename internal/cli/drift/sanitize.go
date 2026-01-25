//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

// formatCheckName converts internal check identifiers to human-readable names.
//
// Parameters:
//   - name: Internal check identifier
//     (e.g., "path_references", "staleness_check")
//
// Returns:
//   - string: Human-readable description of the check, or the original name
//     if unknown
func formatCheckName(name string) string {
	switch name {
	case "path_references":
		return "Path references are valid"
	case "staleness_check":
		return "No staleness indicators"
	case "constitution_check":
		return "Constitution rules respected"
	case "required_files":
		return "All required files present"
	default:
		return name
	}
}
