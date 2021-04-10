package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

var Port = getConf("port", "8061")
var Host = getConf("host", "localhost")

var Env = getConf("env", "local")
var DBUser = getConf("DbUser", "postgres")
var DBPW = getConf("DBPW", "dbpw")

func getConf(key, defaultvalue string) string {

	viper.AutomaticEnv()
	key = strings.ToUpper(key)
	viper.SetDefault(key, defaultvalue)
	return fmt.Sprintf("%v", viper.Get(key))

}
func init() {

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Error("Fatal error config file ", err)
	}

}
