package config

import (
	"fmt"
	"github.com/spf13/viper"
	log "go.uber.org/zap"
	"strings"
)

var Port = getConf("port", "8061")
var Host = getConf("host", "localhost")

var Env = getConf("env", "local")
var DBUser = getConf("DbUser", "postgres")
var DBPW = getConf("DBPW", "dbpw")

func getConf(key, defaultValue string) string {

	viper.AutomaticEnv()
	key = strings.ToUpper(key)
	viper.SetDefault(key, defaultValue)
	return fmt.Sprintf("%v", viper.Get(key))
}
func Init() {

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.S().Error("Fatal error config file ", err)
	}
}
