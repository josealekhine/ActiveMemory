#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Auto-save session transcript on exit (including Ctrl+C)
# This hook is triggered by Claude Code's SessionEnd event
#
# WHAT THIS DOES:
# - Captures the full session transcript when Claude Code exits
# - Saves it to .context/sessions/ with timestamp
# - Works even on Ctrl+C (SIGINT)
#
# SETUP:
# Add to .claude/settings.local.json:
# {
#   "hooks": {
#     "SessionEnd": [{
#       "hooks": [{"type": "command", "command": ".claude/hooks/auto-save-session.sh"}]
#     }]
#   }
# }

# Read hook input from stdin (JSON)
HOOK_INPUT=$(cat)

# Extract transcript_path and session info
TRANSCRIPT_PATH=$(echo "$HOOK_INPUT" | jq -r '.transcript_path // empty')
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id // "unknown"')
REASON=$(echo "$HOOK_INPUT" | jq -r '.reason // "unknown"')
PROJECT_DIR=$(echo "$HOOK_INPUT" | jq -r '.cwd // "."')

# Only proceed if we have a transcript path
if [ -z "$TRANSCRIPT_PATH" ] || [ ! -f "$TRANSCRIPT_PATH" ]; then
  exit 0
fi

# Create sessions directory if it doesn't exist
SESSIONS_DIR="$PROJECT_DIR/.context/sessions"
mkdir -p "$SESSIONS_DIR"

# Generate filename with timestamp: YYYY-MM-DD-HHMMSS-session-<reason>.jsonl
TIMESTAMP=$(date +%Y-%m-%d-%H%M%S)
FILENAME="$SESSIONS_DIR/${TIMESTAMP}-session-${REASON}.jsonl"

# Copy the transcript
cp "$TRANSCRIPT_PATH" "$FILENAME"

# Also create a human-readable summary
SUMMARY_FILE="$SESSIONS_DIR/${TIMESTAMP}-session-${REASON}-summary.md"
cat > "$SUMMARY_FILE" << EOF
# Session Auto-Save

**Saved**: $(date -Iseconds)
**Reason**: $REASON
**Session ID**: $SESSION_ID
**Transcript**: ${TIMESTAMP}-session-${REASON}.jsonl

This session was auto-saved by the SessionEnd hook.
To analyze the full transcript, read the .jsonl file.
EOF

exit 0
