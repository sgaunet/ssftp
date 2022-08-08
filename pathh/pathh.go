package pathh

import (
	"path/filepath"
	"regexp"
	"strings"
)

type Path interface {
	IsRemote() bool
	IsLocal() bool
	GetUser() string
	GetServer() string
	GetFilePath() string
}

// struct to handle path manipulation (local or remote : user@server:dir)
type path struct {
	path string
}

func New(srcPath string) Path {
	p := path{
		path: srcPath,
	}
	return p
}

// Return if the path seems to be a remote path
func (p path) IsRemote() bool {
	r, _ := regexp.Compile(".*@.*:.+")
	return r.Match([]byte(p.path))
}

// Return if the path seems to be a local path
func (p path) IsLocal() bool {
	return !p.IsRemote()
}

// If remote, will return the user
func (p path) GetUser() string {
	splitted := strings.Split(p.path, "@")
	if splitted[0] == p.path {
		return ""
	}
	return splitted[0]
}

// If remote, will return the server (IP or dns name)
func (p path) GetServer() string {
	splitted := strings.Split(p.path, "@")
	if splitted[0] == p.path {
		return ""
	}
	splitted = strings.Split(splitted[1], ":")
	return splitted[0]
}

// Return the path (for a local or remote path)
func (p path) GetFilePath() string {
	if p.IsRemote() {
		splitted := strings.Split(p.path, ":")
		if splitted[0] == p.path {
			return ""
		}
		return splitted[1]
	}
	return filepath.Clean(p.path)
}
