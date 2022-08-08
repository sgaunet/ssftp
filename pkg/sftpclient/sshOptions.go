package sftpclient

import "strings"

// Structs to handle multiple SSHOptions arguments from the command line
type SshOptions []string

func (i *SshOptions) String() string {
	return strings.Join(*i, " ")
}

func (i *SshOptions) Set(value string) error {
	*i = append(*i, value)
	return nil
}
