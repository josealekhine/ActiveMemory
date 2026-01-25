package task

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
)

// tasksFilePath returns the path to TASKS.md.
func tasksFilePath() string {
	return filepath.Join(config.DirContext, config.FilenameTask)
}

// archiveDirPath returns the path to the archive directory.
func archiveDirPath() string {
	return filepath.Join(config.DirContext, config.DirArchive)
}
