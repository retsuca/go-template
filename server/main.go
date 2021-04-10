package server

import (
	"github.com/gorilla/mux"
	"go-template/internal/config"
	"go-template/server/controllers"
	log "go.uber.org/zap"
	"net/http"
)

func init() {
	r := mux.NewRouter()


	r.HandleFunc("/test", controllers.GetTest)

	log.S().Fatal(http.ListenAndServe(":"+config.Get("port"), r))

}
