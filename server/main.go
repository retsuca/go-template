package server

import (
	"github.com/julienschmidt/httprouter"
	"go-template/internal/config"
	"go-template/server/controllers"
	log "go.uber.org/zap"
	"net/http"
)

func init() {
	router := httprouter.New()
	router.GET("/test", controllers.GetTest)

	log.S().Fatal(http.ListenAndServe(":"+config.Port, router))

}
