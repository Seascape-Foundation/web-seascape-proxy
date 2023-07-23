package main

import (
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/proxy"
)

func main() {
	logger, err := log.New("main", false)
	if err != nil {
		log.Fatal("failed to create a log instance", "error", err)
	}

	appConfig, err := configuration.New(logger)
	if err != nil {
		logger.Fatal("configuration.NewAppConfig", "error", err)
	}

	////////////////////////////////////////////////////////////////////////
	//
	// Initialize the proxy
	//
	////////////////////////////////////////////////////////////////////////

	// the proxy creation will validate the config
	web, err := NewWebController(logger)
	if err != nil {
		logger.Fatal("failed to create a web controller", "error", err)
	}

	service := proxy.New(appConfig, logger)
	service.SetCustomSource(web)

	if err != nil {
		logger.Fatal("failed to add source controller to the proxy", "error", err)
	}
	service.Controller.RequireDestination(configuration.ReplierType)

	err = service.Prepare()
	if err != nil {
		logger.Fatal("failed to prepare the service", "error", err)
	}

	service.Run()
}
