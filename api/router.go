package api

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/api/private"
	rLog "github.com/hsaquib/rest-log"
)

func V1Router(logger rLog.Logger) *chi.Mux {
	r := chi.NewRouter()
	privateRouter := private.NewPrivateRouter(logger)
	r.Mount("/private", privateRouter.Router())
	return r
}
