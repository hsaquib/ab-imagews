package api

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/api/private"
	"github.com/hsaquib/ab-imagews/service"
	rLog "github.com/hsaquib/rest-log"
)

func V1Router(provider *service.Provider, logger rLog.Logger) *chi.Mux {
	r := chi.NewRouter()
	privateRouter := private.NewPrivateRouter(provider, logger)
	r.Mount("/private", privateRouter.Router())
	return r
}
