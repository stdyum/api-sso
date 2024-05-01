package http

import (
	netHttp "net/http"

	"github.com/stdyum/api-common/hc"
	"github.com/stdyum/api-sso/internal/use/controllers/controller"
	"github.com/stdyum/api-sso/internal/use/controllers/controller/dto"
)

type Handler interface {
	Run() error

	Login(ctx *hc.Context)
	UpdateToken(ctx *hc.Context)
}

type handler struct {
	port       string
	controller controller.Controller
}

func New(port string, controller controller.Controller) Handler {
	return &handler{
		port:       port,
		controller: controller,
	}
}

func (h *handler) Login(ctx *hc.Context) {
	var requestDTO dto.LoginRequest
	if err := ctx.BindJSON(&requestDTO); err != nil {
		_ = ctx.Error(err)
		return
	}

	tokens, err := h.controller.Login(ctx, requestDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	maxAge := 30 * 24 * 60 * 60 * 1000
	domain := "." + ctx.Request.Host
	ctx.SetCookie("access", tokens.Access, maxAge, "/", domain, true, true)
	ctx.SetCookie("refresh", tokens.Refresh, maxAge, "/", domain, true, true)

	ctx.Status(netHttp.StatusNoContent)
}

func (h *handler) UpdateToken(ctx *hc.Context) {
	refresh, err := ctx.Cookie("refresh")
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	requestDTO := dto.UpdateRequest{
		RefreshToken: refresh,
	}

	tokens, err := h.controller.Update(ctx, requestDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	maxAge := 30 * 24 * 60 * 60 * 1000
	ctx.SetCookie("access", tokens.Access, maxAge, "/", "", true, true)
	ctx.SetCookie("refresh", tokens.Refresh, maxAge, "/", "", true, true)

	ctx.Status(netHttp.StatusNoContent)
}
