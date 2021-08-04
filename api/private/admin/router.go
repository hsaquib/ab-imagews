package admin

import (
	"github.com/go-chi/chi"
	rLog "github.com/hsaquib/rest-log"
)

type adminRouter struct {
	Log rLog.Logger
}

func NewAdminRouter(rLogger rLog.Logger) *adminRouter {
	return &adminRouter{
		Log: rLogger,
	}
}

func (ar *adminRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/upload", UploadImage)

	//r.With(middleware.AuthenticatedAdminOnly).Post("/signup", cr.signup)
	//r.With(middleware.AuthenticatedAdminOnly).Get("/profile", cr.getAdminProfile)
	//r.With(middleware.AuthenticatedAdminOnly).Patch("/profile", cr.updateAdminProfile)
	//r.With(middleware.AuthenticatedAdminOnly).Get("/verify-token", cr.verifyAccessToken)
	//r.With(middleware.JWTTokenOnly).Get("/refresh-token", cr.refreshToken)
	//r.With(middleware.AuthenticatedAdminOnly).Put("/password", cr.updatePassword)

	return r
}
