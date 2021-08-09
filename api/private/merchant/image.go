package merchant

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/api/private/handler"
	"github.com/hsaquib/ab-imagews/middleware"
	"net/http"
)

func (mr *merchantRouter) imageRouter() http.Handler {
	r := chi.NewRouter()
	reqHandler := handler.NewApiHandler(mr.ServiceProvider, mr.Log)

	r.With(middleware.AuthenticatedMerchantOnly).Post("/upload", reqHandler.UploadImageWithVariantsByMerchant)
	return r
}
