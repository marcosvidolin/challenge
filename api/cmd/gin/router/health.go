package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type cacheAdapter interface {
	Ping(ctx context.Context) error
}

type dbAdapter interface {
	Ping() error
}

type healthRouter struct {
	db    dbAdapter
	cache cacheAdapter
}

func NewHealthRouter(db dbAdapter, cache cacheAdapter) *healthRouter {
	return &healthRouter{
		db:    db,
		cache: cache,
	}
}

func (r *healthRouter) HealthRouter(api *gin.RouterGroup) {
	api.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()

		dbStatus := "connected"
		if err := r.db.Ping(); err != nil {
			dbStatus = "not connected"
		}

		cacheStatus := "connected"
		if err := r.cache.Ping(ctx); err != nil {
			cacheStatus = "not connected"
		}

		c.JSON(http.StatusOK, gin.H{
			"database": dbStatus,
			"cache":    cacheStatus,
		})
	})
}
