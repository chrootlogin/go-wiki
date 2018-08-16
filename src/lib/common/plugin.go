package common

import "github.com/gin-gonic/gin"

type GoWikiPlugin interface {
	Run(e *gin.Engine)
}
