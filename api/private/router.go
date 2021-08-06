package private

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/api/private/admin"
	"github.com/hsaquib/ab-imagews/api/private/merchant"
	"github.com/hsaquib/ab-imagews/service"
	rLog "github.com/hsaquib/rest-log"
)

type privateRouter struct {
	serviceProvider *service.Provider
	Log             rLog.Logger
}

func NewPrivateRouter(provider *service.Provider, rLogger rLog.Logger) *privateRouter {
	return &privateRouter{
		serviceProvider: provider,
		Log:             rLogger,
	}
}

func (pr *privateRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/admin", admin.NewAdminRouter(pr.serviceProvider, pr.Log).Router())
	r.Mount("/merchant", merchant.NewMerchantRouter(pr.serviceProvider, pr.Log).Router())
	return r
}
