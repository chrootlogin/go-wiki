package plugins

import "github.com/coreos/go-semver/semver"

type GoWikiPlugin interface {
	Version() semver.Version
}
