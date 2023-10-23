package webproxy

import (
	"fmt"
	"github.com/ahmetson/handler-lib/base"
	"github.com/ahmetson/handler-lib/config"
	"github.com/ahmetson/service-lib"
	webHandler "github.com/ahmetson/web-lib"
)

// The New method returns a service with the web handlers.
func New() (*service.Proxy, error) {
	webDefiner := func() base.Interface {
		return webHandler.New()
	}

	proxy, err := service.NewProxy()
	if err != nil {
		return nil, fmt.Errorf("service.NewProxy: %w", err)
	}
	proxy.SetHandlerDefiner(config.ReplierType, webDefiner)
	proxy.SetHandlerDefiner(config.SyncReplierType, webDefiner)

	return proxy, nil
}
