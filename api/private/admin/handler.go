package admin

import (
	"github.com/hsaquib/ab-imagews/service"
	rLog "github.com/hsaquib/rest-log"
)

type apiHandler struct {
	ServiceProvider *service.Provider
	Log             rLog.Logger
}

func NewApiHandler(srvProvider *service.Provider, log rLog.Logger) *apiHandler {
	return &apiHandler{
		ServiceProvider: srvProvider,
		Log:             log,
	}
}
