package common

import (
	"github.com/spf13/viper"
)

func ConfigInit() {
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic("Failed while reding configuration")
	}
}

func GetServerPort() string {
	return viper.GetString("server.port")
}

func GetDatabaseHost() string {
	return viper.GetString("database.host")
}

func GetBasicUsername() string {
	return viper.GetString("basic.username")
}

func GetBasicPassword() string {
	return viper.GetString("basic.password")
}

func GetBearerKey() string {
	return viper.GetString("bearer.key")
}

func GetRefreshKey() string {
	return viper.GetString("bearer.refreshKey")
}
