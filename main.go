package main

import (
	"fmt"

	"github.com/Dparty/core/controllers"
	"github.com/Dparty/model"
	"github.com/spf13/viper"
)

func main() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	db, err := model.NewConnection(user, password, host, port, database)
	if err != nil {
		panic(err)
	}
	controllers.Init(db)
	controllers.Run(":8080")
}
