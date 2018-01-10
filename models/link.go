package models

import (
	"regexp"
)

var re = regexp.MustCompile(`^((\.\./|[a-zA-Z0-9_/\-\\])*)$`)

type Link string

func NewLink(path string) *Link {
	var link = Link(path)
	return &link
}

func (link *Link) Raw() string {
	return string(*link)
}

func (link *Link) Parsed() bool {
	return re.MatchString(string(*link))
}
