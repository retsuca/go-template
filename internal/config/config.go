package config

import (
	"fmt"
	"github.com/spf13/viper"
	log "go.uber.org/zap"
	"strings"
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
