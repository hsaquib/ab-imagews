package private

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/api/private/admin"
	rLog "github.com/hsaquib/rest-log"
)

type privateRouter struct {
	Log rLog.Logger
}

func NewPrivateRouter(rLogger rLog.Logger) *privateRouter {
	return &privateRouter{
		Log: rLogger,
	}
}

func (pr *privateRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/admin", admin.NewAdminRouter(pr.Log).Router())
	return r
}
