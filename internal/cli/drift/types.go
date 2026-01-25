//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import "github.com/ActiveMemory/ctx/internal/drift"

// JsonOutput represents the JSON structure for machine-readable drift output.
//
// Fields:
//   - Timestamp: RFC3339-formatted UTC time when the report was generated
//   - Status: Overall drift status ("ok", "warning", or "violation")
//   - Warnings: Issues that should be addressed but don't block
//   - Violations: Constitution violations that must be fixed
//   - Passed: Names of checks that passed successfully
type JsonOutput struct {
	Timestamp  string        `json:"timestamp"`
	Status     string        `json:"status"`
	Warnings   []drift.Issue `json:"warnings"`
	Violations []drift.Issue `json:"violations"`
	Passed     []string      `json:"passed"`
}
