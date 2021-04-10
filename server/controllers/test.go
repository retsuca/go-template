package controllers

import (
	"fmt"
	"go-template/internal/config"
	"net/http"
)

func GetTest(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	println(config.Get("port"))
	fmt.Fprint(w, "Test!\n")
}
