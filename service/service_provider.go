package service

import (
	"github.com/hsaquib/ab-imagews/config"
	rLog "github.com/hsaquib/rest-log"
)

type Provider struct {
	FileProcessor *FileProcessor
}

func InitProvider(cfg *config.AppConfig, logger rLog.Logger) *Provider {

	return &Provider{
		FileProcessor: NewFileProcessor(cfg, logger),
	}
}
