package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "go-template/docs"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"

	httpclient "go-template/internal/clients/httpClient"
	"go-template/internal/config"
	logger "go-template/pkg/logger"
	"go-template/server/http/handler"
)

func CreateHTPPServer(ctx context.Context, host, port string) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(""))

	appName := config.Get(config.APP_NAME)

	e.Use(echoprometheus.NewMiddleware(strings.ReplaceAll(appName, "-", "_")))

	e.GET("/metrics", echoprometheus.NewHandler())

	client := httpclient.NewClient(httpclient.ClientOptions{
		BaseURL: &url.URL{Scheme: "https", Host: "hacker-news.firebaseio.com", Path: "v0/"},
		// Sock5Proxy:         config.Get(config.SOCKS5_PROXY),
		InsecureSkipVerify: false,
	})

	h := handler.NewHandler(client)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", h.Hello)
	e.GET("/withparam", h.HelloWithParam)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(host + ":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Fatal error http server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	//nolint:mnd //because
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
