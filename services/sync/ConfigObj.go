package sync_services

import (
	"path/filepath"
	"strings"
)

type ConfigObj struct {
	KeyName        string
	FileName       string // with extension
	Content        []byte
	DeterminedPath DeterminedPath
}

type DeterminedPath struct {
	Path        string
	FromHomeDir bool
}

func (p *DeterminedPath) BuildPath(homeDir string) string {
	if p.FromHomeDir {
		return filepath.Join(homeDir, p.Path)
	} else {
		return p.Path
	}
}

func NewDeterminedPath(path, homeDir string) DeterminedPath {
	return DeterminedPath{
		Path:        strings.TrimPrefix(path, homeDir),
		FromHomeDir: strings.HasSuffix(path, homeDir),
	}
}
