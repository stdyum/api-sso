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
	Logout(ctx *hc.Context)
	UpdateToken(ctx *hc.Context)

	Authorize(ctx *hc.Context)

	SetTokens(ctx *hc.Context)

	SetDefaultEnrollmentId(ctx *hc.Context)
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

	h.setTokenCookies(ctx, tokens)

	ctx.Status(netHttp.StatusNoContent)
}

func (h *handler) Logout(ctx *hc.Context) {
	h.setDefaultSubdomainCookieWithMaxAge(ctx, "access", "", -1)
	h.setDefaultSubdomainCookieWithMaxAge(ctx, "refresh", "", -1)

	ctx.Status(netHttp.StatusNoContent)
}

func (h *handler) UpdateToken(ctx *hc.Context) {
	refresh, err := ctx.Cookie("refresh")
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	enrollment, _ := ctx.Cookie("enrollment")

	request := dto.UpdateRequest{
		RefreshToken: refresh,
	}

	response, err := h.controller.Update(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	h.setTokenCookies(ctx, response.Tokens)

	response.Enrollment.Id = enrollment

	ctx.JSON(netHttp.StatusOK, response)
}

func (h *handler) Authorize(ctx *hc.Context) {
	access, err := ctx.Cookie("access")
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	refresh, _ := ctx.Cookie("refresh")
	enrollment, _ := ctx.Cookie("enrollment")

	request := dto.AuthorizeRequest{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	response, err := h.controller.Authorize(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	h.setTokenCookies(ctx, response.Tokens)

	response.Enrollment.Id = enrollment

	ctx.JSON(netHttp.StatusOK, response)
}

func (h *handler) SetTokens(ctx *hc.Context) {
	var request dto.SetTokensRequest
	if err := ctx.BindJSON(&request); err != nil {
		_ = ctx.Error(err)
		return
	}

	token := dto.TokenPairResponse{
		Access:  request.AccessToken,
		Refresh: request.RefreshToken,
	}

	h.setTokenCookies(ctx, token)

	ctx.Status(netHttp.StatusOK)
}

func (h *handler) SetDefaultEnrollmentId(ctx *hc.Context) {
	var request dto.SetDefaultEnrollmentIdRequest
	if err := ctx.BindJSON(&request); err != nil {
		_ = ctx.Error(err)
		return
	}

	h.setDefaultSubdomainCookie(ctx, "enrollment", request.Id)

	ctx.Status(netHttp.StatusNoContent)
}

func (h *handler) setTokenCookies(ctx *hc.Context, tokens dto.TokenPairResponse) {
	h.setDefaultSubdomainCookie(ctx, "access", tokens.Access)
	h.setDefaultSubdomainCookie(ctx, "refresh", tokens.Refresh)
}

func (h *handler) setDefaultSubdomainCookie(ctx *hc.Context, name, value string) {
	maxAge := 30 * 24 * 60 * 60 * 1000
	h.setDefaultSubdomainCookieWithMaxAge(ctx, name, value, maxAge)
}

func (h *handler) setDefaultSubdomainCookieWithMaxAge(ctx *hc.Context, name, value string, maxAge int) {
	domain := h.generateCookieHost(ctx.Request.Host)
	ctx.SetCookie(name, value, maxAge, "/", domain, true, true)
}

func (h *handler) generateCookieHost(root string) string {
	dotAmount := 0
	host := ""
	for i := len(root) - 1; i >= 0; i-- {
		if root[i] == '/' && i != 0 && host[i-1] != '/' {
			dotAmount = 0
			host = ""
			continue
		}

		if dotAmount == 2 {
			continue
		}

		if root[i] == '.' {
			dotAmount++
			if dotAmount == 2 {
				continue
			}
		}

		host += string(root[i])
	}

	return "." + reverse(host)
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
