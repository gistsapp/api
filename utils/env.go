package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Get(key string) string {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()

	if err != nil {
		return os.Getenv(key) //for production
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		msg := fmt.Sprintf("Env variable '%q' not found\n", key)
		panic(msg)
		return ""
	}

	return value
}
