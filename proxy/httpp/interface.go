package httpp

import (
	"net/http"
	"regexp"
)


type HostPath struct {
	Host, Path string
}

type Regex struct {
	Match, Sub string

	match *regexp.Regexp
}

func (t *Regex) Compile() error {
	var err error
	t.match, err = regexp.Compile(t.Match)
	return err
}

func (t *Regex) Apply(req *http.Request) {
	t.match.ReplaceAllString(req.URL.Path, t.Sub)
	req.URL.RawPath = req.URL.EscapedPath()
}

type Transform struct {
	// Apply a regular expression to adapt the resulting URL.
	Regex *Regex
	// Maintain the original path of the proxy. Normally, it is stripped.
	// For example: if you map "proxy.address/path/p1/" to "backend.address/path2/", a request
	// for "proxy.address/path/p1/test" will by default land to "backend.address/path2/test".
	// If you set Maintain to true, it will instead land on "backend.address/path2/path/p1/test".
	Maintain bool
}

func (t *Transform) Apply(req *http.Request) bool {
	if t.Regex != nil {
		t.Regex.Apply(req)
	}
	return t.Maintain
}

func (t *Transform) Compile(fromurl, tourl string) error {
	if t.Regex != nil {
		return t.Regex.Compile()
	}
	return nil
}

type MappingAuth string

const (
	MappingAuthenticated MappingAuth = ""
	MappingPublic        MappingAuth = "public"
)

type Mapping struct {
	From HostPath
	To   string

	Transform []*Transform
	Auth      MappingAuth
}

type Config struct {
	// Which URLs to map to which other URLs.
	Mapping []Mapping

	// Extra domains for which to obtain a certificate.
	Domains []string
}


