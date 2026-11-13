// Package httpapi provides the local HTTP transport. Health reports only that
// this HTTP application responds; it does not check workflow or provider readiness.
package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter constructs routes without starting a listener or accessing the engine.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return router
}
