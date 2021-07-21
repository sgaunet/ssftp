package pathh

import (
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

type path struct {
	path string
}

func New(srcPath string) Path {
	p := path{
		path: srcPath,
	}
	return p
}

func (p path) IsRemote() bool {
	// r, _ := regexp.Compile("([a-z]+)@([a-z]+)")
	// r, _ := regexp.Compile("([a-z]+)@([a-z]+)")
	// return r.Match([]byte(p.path))
	ipv6_regex := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	ipv4_regex := `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
	domain_regex := `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
	localserver_regex := `^[a-zA-Z0-9]+$`

	match, _ := regexp.MatchString(ipv4_regex+`|`+ipv6_regex+`|`+domain_regex+`|`+localserver_regex, p.GetServer())
	return match
}

func (p path) IsLocal() bool {
	return !p.IsRemote()
}

func (p path) GetUser() string {
	splitted := strings.Split(p.path, "@")
	if splitted[0] == p.path {
		return ""
	}
	return splitted[0]
}

func (p path) GetServer() string {
	splitted := strings.Split(p.path, "@")
	if splitted[0] == p.path {
		return ""
	}

	splitted = strings.Split(splitted[1], ":")
	return splitted[0]
}

func (p path) GetFilePath() string {

	if p.IsRemote() {
		splitted := strings.Split(p.path, ":")
		if splitted[0] == p.path {
			return ""
		}

		return splitted[1]
	}
	return p.path
}
