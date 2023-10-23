package main

import (
	"github.com/ahmetson/log-lib"
	webProxy "github.com/ahmetson/web-proxy"
)

func main() {
	logger, err := log.New("web-proxy", false)
	if err != nil {
		log.Fatal("failed to create a log instance", "error", err)
	}

	proxy, err := webProxy.New()
	if err != nil {
		logger.Fatal("service.NewProxy", "error", err)
	}

	wg, err := proxy.Start()
	if err != nil {
		logger.Fatal("proxy.Start", "error", err)
	}

	wg.Wait()
}
