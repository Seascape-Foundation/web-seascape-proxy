package main

import (
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/proxy"
	"github.com/ahmetson/service-lib/remote"
)

func main() {
	logger, err := log.New("main", false)
	if err != nil {
		log.Fatal("failed to create a log instance", "error", err)
	}

	appConfig, err := configuration.NewAppConfig(logger)
	if err != nil {
		log.Fatal("configuration.NewAppConfig", "error", err)
	}
	if len(appConfig.Services) == 0 {
		log.Fatal("seascape.yml doesn't have services section")
	}

	////////////////////////////////////////////////////////////////////////
	//
	// Initialize the proxy
	//
	////////////////////////////////////////////////////////////////////////

	// We won't handle anything
	handler := func(messages []string,
		controllerLogger log.Logger,
		_ []*proxy.DestinationClient,
		clients remote.Clients) ([]string, error) {
		controllerLogger.Info("request", "messages", messages)
		return messages, nil
	}

	// the proxy creation will validate the config
	sourceConfig, err := appConfig.Services[0].GetController(proxy.SourceName)
	if err != nil {
		log.Fatal("failed to get source controller's configuration from seascape.yml", "error", err)
	}
	web, err := NewWebController(logger)
	web.AddConfig(sourceConfig)

	service, err := proxy.New(appConfig.Services[0], logger)
	if err != nil {
		log.Fatal("proxy.New", "error", err)
	}

	service.SetRequestHandler(handler)
	err = service.AddSourceController(proxy.SourceName, web)
	if err != nil {
		log.Fatal("failed to add source controller to the proxy", "error", err)
	}

	service.Run()
}
