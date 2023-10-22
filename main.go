package main

import (
	"github.com/ahmetson/handler-lib/base"
	"github.com/ahmetson/handler-lib/config"
	"github.com/ahmetson/log-lib"
	"github.com/ahmetson/service-lib"
	webHandler "github.com/ahmetson/web-lib"
)

func main() {
	logger, err := log.New("web-proxy", false)
	if err != nil {
		log.Fatal("failed to create a log instance", "error", err)
	}

	proxy, err := service.NewProxy()
	if err != nil {
		logger.Fatal("service.NewProxy", "error", err)
	}
	webDefiner := func() base.Interface {
		return webHandler.New()
	}
	proxy.SetHandlerDefiner(config.ReplierType, webDefiner)
	proxy.SetHandlerDefiner(config.SyncReplierType, webDefiner)

	wg, err := proxy.Start()
	if err != nil {
		logger.Fatal("proxy.Start", "error", err)
	}

	wg.Wait()
}
