package plugins

import (
	"github.com/coreos/go-semver/semver"
	"github.com/gin-gonic/gin"
)

type GoWikiPlugin interface {
	Version() *semver.Version
	RunEngine(*gin.Engine)
}
