package admin

import (
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/api/private/handler"
	"github.com/hsaquib/ab-imagews/middleware"
	"net/http"
)

func (ar *adminRouter) imageRouter() http.Handler {
	r := chi.NewRouter()
	reqHandler := handler.NewApiHandler(ar.ServiceProvider, ar.Log)

	r.With(middleware.AuthenticatedAdminOnly).Post("/upload_with_variants", reqHandler.UploadImageWithVariantsByAdmin)
	r.Post("/upload_mult", reqHandler.UploadImageListWithVariantsByMerchant)
	return r
}
