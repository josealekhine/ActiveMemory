package config

const (
	CtxMarkerEnd           = "<!-- ctx:end -->"
	CtxMarkerStart         = "<!-- ctx:context -->"
	DirArchive             = "archive"
	DirClaude              = ".claude"
	DirClaudeHooks         = ".claude/hooks"
	DirContext             = ".context"
	DirSessions            = "sessions"
	FileAutoSave           = "auto-save-session.sh"
	FileBlockNonPathScript = "block-non-path-ctx.sh"
	FileClaudeMd           = "CLAUDE.md"
	FileSettings           = ".claude/settings.local.json"
)

// WatchAutoSaveInterval is the number of updates between auto-saves
// in watch mode.
const WatchAutoSaveInterval = 5

// Context file name constants.
const (
	FilenameConstitution  = "CONSTITUTION.md"
	FilenameTask          = "TASKS.md"
	FilenameConvention    = "CONVENTIONS.md"
	FilenameArchitecture  = "ARCHITECTURE.md"
	FilenameDecision      = "DECISIONS.md"
	FilenameLearning      = "LEARNINGS.md"
	FilenameGlossary      = "GLOSSARY.md"
	FilenameDrift         = "DRIFT.md"
	FilenameAgentPlaybook = "AGENT_PLAYBOOK.md"
	FilenameDependency    = "DEPENDENCIES.md"
)

// Update type constants for context entries.
//
// These are used in switch statements for routing add/update commands
// to the appropriate handler.
const (
	UpdateTypeTask       = "task"
	UpdateTypeDecision   = "decision"
	UpdateTypeLearning   = "learning"
	UpdateTypeConvention = "convention"
	UpdateTypeComplete   = "complete"
)

// Plural aliases for update types.
//
// Accepted as synonyms for the singular forms.
const (
	UpdateTypeTasks       = "tasks"
	UpdateTypeDecisions   = "decisions"
	UpdateTypeLearnings   = "learnings"
	UpdateTypeConventions = "conventions"
)

// FileType maps short names to actual file names.
var FileType = map[string]string{
	UpdateTypeDecision:    FilenameDecision,
	UpdateTypeDecisions:   FilenameDecision,
	UpdateTypeTask:        FilenameTask,
	UpdateTypeTasks:       FilenameTask,
	UpdateTypeLearning:    FilenameLearning,
	UpdateTypeLearnings:   FilenameLearning,
	UpdateTypeConvention:  FilenameConvention,
	UpdateTypeConventions: FilenameConvention,
}

// RequiredFiles lists the essential context files that must be present.
//
// These are the files created with `ctx init --minimal` and checked by
// drift detection for missing files.
var RequiredFiles = []string{
	FilenameConstitution,
	FilenameTask,
	FilenameDecision,
}

// FileReadOrder defines the priority order for reading context files.
//
// The order follows a logical progression for AI agents:
//
//  1. CONSTITUTION — Inviolable rules. Must be loaded first so the agent
//     knows what it cannot do before attempting anything.
//
//  2. TASKS — Current work items. What the agent should focus on.
//
//  3. CONVENTIONS — How to write code. Patterns and standards to follow.
//
//  4. ARCHITECTURE — System structure. Understanding of components and
//     boundaries before making changes.
//
//  5. DECISIONS — Historical context. Why things are the way they are,
//     to avoid re-debating settled decisions.
//
//  6. LEARNINGS — Gotchas and tips. Lessons from past work that inform
//     current implementation.
//
//  7. GLOSSARY — Reference material. Domain terms and abbreviations for
//     lookup as needed.
//
//  8. DRIFT — Staleness indicators. Lower priority since it's primarily
//     for maintenance workflows.
//
//  9. AGENT_PLAYBOOK — Meta instructions. How to use this context system.
//     Loaded last because it's about the system itself, not the work.
//     The agent should understand the content before the operating manual.
var FileReadOrder = []string{
	FilenameConstitution,
	FilenameTask,
	FilenameConvention,
	FilenameArchitecture,
	FilenameDecision,
	FilenameLearning,
	FilenameGlossary,
	FilenameDrift,
	FilenameAgentPlaybook,
}

// filePriority maps filenames to their priority (derived from FileReadOrder).
var filePriority = func() map[string]int {
	m := make(map[string]int, len(FileReadOrder))
	for i, name := range FileReadOrder {
		m[name] = i + 1
	}
	return m
}()

// Packages maps dependency manifest files to their descriptions.
//
// Used by sync to detect projects and suggest dependency documentation.
var Packages = map[string]string{
	"package.json":     "Node.js dependencies",
	"go.mod":           "Go module dependencies",
	"Cargo.toml":       "Rust dependencies",
	"requirements.txt": "Python dependencies",
	"Gemfile":          "Ruby dependencies",
}

// Pattern represents a config file pattern and its documentation topic.
type Pattern struct {
	Pattern string // Glob pattern to match (e.g., ".eslintrc*")
	Topic   string // Documentation topic (e.g., "linting conventions")
}

// Patterns lists config files that should be documented in CONVENTIONS.md.
//
// Used by sync to suggest documenting project configuration.
var Patterns = []Pattern{
	{".eslintrc*", "linting conventions"},
	{".prettierrc*", "formatting conventions"},
	{"tsconfig.json", "TypeScript configuration"},
	{".editorconfig", "editor configuration"},
	{"Makefile", "build commands"},
	{"Dockerfile", "containerization"},
}

// FilePriority returns the priority of a context file.
//
// Lower numbers indicate higher priority (1 = highest).
// Unknown files return 100.
//
// Parameters:
//   - name: Filename to look up (e.g., "TASKS.md")
//
// Returns:
//   - int: Priority value (1-9 for known files, 100 for unknown)
func FilePriority(name string) int {
	if p, ok := filePriority[name]; ok {
		return p
	}
	return 100
}
