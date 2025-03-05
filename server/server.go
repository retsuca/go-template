package server

import (
	"context"
	"errors"
	"go-template/internal/config"
	logger "go-template/pkg/logger"
	"go-template/server/handler"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	httpclient "go-template/internal/clients/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateHTPPServer(host, port string) {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(""))

	appName := config.Get(config.APP_NAME)

	e.Use(echoprometheus.NewMiddleware(strings.Replace(appName, "-", "_", -1)))

	e.GET("/metrics", echoprometheus.NewHandler())

	client := httpclient.NewClient(&url.URL{Scheme: "https", Host: "hacker-news.firebaseio.com", Path: "v0/"})

	h := handler.NewHandler(client)

	e.GET("/", h.Hello)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {

		if err := e.Start(host + ":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FatalErr("Fatal error http server ", err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
