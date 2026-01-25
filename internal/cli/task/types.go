//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

// taskStats holds counts of completed and pending tasks.
//
// Used by separateTasks to report how many tasks were processed during
// an archive operation.
//
// Fields:
//   - completed: Number of tasks marked with [x]
//   - pending: Number of tasks marked with [ ]
type taskStats struct {
	completed int
	pending   int
}
