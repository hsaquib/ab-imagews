package merchant

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/service"
	rLog "github.com/hsaquib/rest-log"
)

type merchantRouter struct {
	ServiceProvider *service.Provider
	Log             rLog.Logger
}

func NewMerchantRouter(provider *service.Provider, logger rLog.Logger) *merchantRouter {
	return &merchantRouter{
		ServiceProvider: provider,
		Log:             logger,
	}
}

func (mr *merchantRouter) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/image", mr.imageRouter())
	return r
}
