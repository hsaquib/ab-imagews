package admin

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/service"
	rLog "github.com/hsaquib/rest-log"
)

type adminRouter struct {
	ServiceProvider *service.Provider
	Log             rLog.Logger
}

func NewAdminRouter(provider *service.Provider, rLogger rLog.Logger) *adminRouter {
	return &adminRouter{
		ServiceProvider: provider,
		Log:             rLogger,
	}
}

func (ar *adminRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/image", ar.imageRouter())

	//r.With(middleware.AuthenticatedAdminOnly).Post("/signup", cr.signup)
	//r.With(middleware.AuthenticatedAdminOnly).Get("/profile", cr.getAdminProfile)
	//r.With(middleware.AuthenticatedAdminOnly).Patch("/profile", cr.updateAdminProfile)
	//r.With(middleware.AuthenticatedAdminOnly).Get("/verify-token", cr.verifyAccessToken)
	//r.With(middleware.JWTTokenOnly).Get("/refresh-token", cr.refreshToken)
	//r.With(middleware.AuthenticatedAdminOnly).Put("/password", cr.updatePassword)

	return r
}
