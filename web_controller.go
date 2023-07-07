package main

import (
	"fmt"
	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-service-lib/communication/command"
	"github.com/Seascape-Foundation/sds-service-lib/configuration"
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
	"github.com/valyala/fasthttp"
)

type WebController struct {
	Config             *configuration.Controller
	logger             log.Logger
	requiredExtensions []string
	extensionConfigs   key_value.KeyValue
	extensions         remote.Clients
}

func NewWebController(parent log.Logger) (*WebController, error) {
	logger, err := parent.Child("web-controller", true)
	if err != nil {
		return nil, fmt.Errorf("failed to create a log: %w", err)
	}

	webController := WebController{
		logger:             logger,
		requiredExtensions: make([]string, 0),
		extensionConfigs:   key_value.Empty(),
		extensions:         make(remote.Clients, 0),
	}

	return &webController, nil
}

func (web *WebController) AddConfig(config *configuration.Controller) {
	web.Config = config
}

// AddExtensionConfig adds the configuration of the extension that the controller depends on
func (web *WebController) AddExtensionConfig(extension *configuration.Extension) {
	web.extensionConfigs.Set(extension.Name, extension)
}

// RequireExtension marks the extensions that this controller depends on.
// Before running, the required extension should be added from the configuration.
// Otherwise, controller won't run.
func (web *WebController) RequireExtension(name string) {
	web.requiredExtensions = append(web.requiredExtensions, name)
}

// RequiredExtensions returns the list of extension names required by this controller
func (web *WebController) RequiredExtensions() []string {
	return web.requiredExtensions
}

// RegisterCommand adds a command along with its handler to this controller
func (web *WebController) RegisterCommand(_ command.Name, _ command.HandleFunc) {
	web.logger.Fatal("not implemented")
}

func (web *WebController) initExtensionClients() error {
	for _, extensionInterface := range web.extensionConfigs {
		extensionConfig := extensionInterface.(*configuration.Extension)
		extension, err := remote.NewReq(extensionConfig.Name, extensionConfig.Port, &web.logger)
		if err != nil {
			return fmt.Errorf("failed to create a request client: %w", err)
		}
		web.extensions.Set(extensionConfig.Name, extension)
	}

	return nil
}

func (web *WebController) Run() error {
	if len(web.Config.Instances) == 0 {
		return fmt.Errorf("no instance of the config")
	}

	// todo
	// init extension clients

	instanceConfig := web.Config.Instances[0]
	if instanceConfig.Port == 0 {
		web.logger.Fatal("instance port is invalid",
			"controller", instanceConfig.Name,
			"instance", instanceConfig.Instance,
			"port", instanceConfig.Port,
		)
	}

	addr := fmt.Sprintf(":%d", instanceConfig.Port)

	if err := fasthttp.ListenAndServe(addr, requestHandler); err != nil {
		return fmt.Errorf("error in ListenAndServe: %w at port %d", err, instanceConfig.Port)
	}

	return fmt.Errorf("http server was down")
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintf(ctx, "Hello, world!\n\n")

	_, _ = fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	_, _ = fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	_, _ = fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	_, _ = fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	_, _ = fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	_, _ = fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	_, _ = fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	_, _ = fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	_, _ = fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	_, _ = fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	_, _ = fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
	_, _ = fmt.Println(ctx, ctx.PostBody())

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)

	_, _ = ctx.WriteString("web proxy finished")
}
