package main

import (
	"github.com/Seascape-Foundation/sds-service-lib/configuration"
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/sds-service-lib/proxy"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
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
		_ log.Logger,
		_ []*proxy.DestinationClient,
		_ remote.Clients) ([]string, string, error) {
		return messages, "destination", nil
	}

	service, err := proxy.New(appConfig.Services[0], logger)
	if err != nil {
		log.Fatal("proxy.New", "error", err)
	}

	service.SetRequestHandler(handler)
	// the proxy needs two parts:
	// 1. source controllers
	// 2. handler
}
