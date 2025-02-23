package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	log "go.uber.org/zap"
)

const (
	APP_NAME                    = "APP_NAME"
	HTTP_PORT                   = "HTTP_PORT"
	HTTP_HOST                   = "HTTP_HOST"
	ENV                         = "ENV"
	DB_ADDRESS                  = "DB_ADDRESS"
	DB_NAME                     = "DB_NAME"
	DB_USER                     = "DB_USER"
	DB_PW                       = "DB_PW"
	OTEL_EXPORTER_OTLP_ENDPOINT = "OTEL_EXPORTER_OTLP_ENDPOINT"
)

func Get(key string) string {

	viper.AutomaticEnv()
	key = strings.ToUpper(key)
	return fmt.Sprintf("%v", viper.Get(key))
}
func init() {

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	viper.WatchConfig()

	if err != nil {
		log.S().Error("Fatal error config file ", err)
	}
	return
}
