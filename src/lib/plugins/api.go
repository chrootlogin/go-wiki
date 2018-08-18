package plugins

import "github.com/gin-gonic/gin"

// Exported Interface
type GoWikiPluginAPIInterface interface {
	Engine() *gin.Engine
}

type GoWikiPluginAPI struct {
	engine *gin.Engine
}

func (api GoWikiPluginAPI) Engine() *gin.Engine {
	return api.engine
}

func (api GoWikiPluginAPI) SetEngine(e *gin.Engine) {
	api.engine = e
}