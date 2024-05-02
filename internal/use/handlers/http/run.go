package http

import (
	"github.com/stdyum/api-common/hc"
	"github.com/stdyum/api-common/http/middlewares"
)

func (h *handler) Run() error {
	engine := hc.New()
	engine.Use(hc.Recovery())

	group := engine.Group("api/v1", hc.Logger(), middlewares.ErrorMiddleware())
	{
		group.POST("login", h.Login)
		group.POST("update", h.UpdateToken)
		group.POST("authorize", h.Authorize)
	}

	return engine.Run(":" + h.port)
}
