package swagger

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/config"
	"github.com/hsaquib/ab-imagews/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(cfg *config.AppConfig) *chi.Mux {
	r := chi.NewRouter()

	if cfg.Env == utils.DEV_ENV {
		r.Get("/*", httpSwagger.Handler())
	}

	return r
}
