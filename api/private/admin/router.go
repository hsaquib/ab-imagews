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
	return r
}
