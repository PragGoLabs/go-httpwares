// Copyright 2017 Mark Nevill. All Rights Reserved.
// See LICENSE for licensing terms.

package http_prometheus

import (
	"net/http"
	"regexp"

	"github.com/Solution/go-httpwares/tags"
)

type meta struct {
	name, handler, method, host, path string
}

var compiledPaths = map[string]*regexp.Regexp{}

func reqMeta(req *http.Request, opts *options, inbound bool) *meta {
	m := &meta{name: opts.name, method: req.Method}

	var tags map[string]interface{}
	if inbound {
		tags = http_ctxtags.ExtractInbound(req).Values()
	} else {
		tags = http_ctxtags.ExtractOutbound(req).Values()
	}
	var v interface{}
	if m.name == "" {
		v, _ = tags[http_ctxtags.TagForCallService]
		m.name, _ = v.(string)
	}
	v, _ = tags[http_ctxtags.TagForHandlerName]
	hname, _ := v.(string)
	if hname != "" {
		v, _ = tags[http_ctxtags.TagForHandlerGroup]
		hgroup, _ := v.(string)
		if hgroup == "" {
			hgroup = "unknown"
		}
		m.handler = hgroup + "." + hname
	}

	if opts.hosts {
		m.host = req.URL.Host
		if m.host == "" {
			m.host = req.Host
		}
	}

	if opts.paths {
		m.path = resolvePath(opts, req.URL.Path)
	}

	return m
}


func resolvePath(opts *options, path string) string {
	if !opts.registerUniquePaths {
		return path
	}

	for _, uniquePath := range opts.uniquePaths {
		if compiledPaths[uniquePath] == nil {
			compiledPaths[uniquePath] = regexp.MustCompile(uniquePath)
		}

		if compiledPaths[uniquePath].Match([]byte(path)) {
			return uniquePath
		}
	}

	return path
}
