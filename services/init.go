package services

import (
	"github.com/Dparty/model"
	"github.com/Dparty/model/core"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// var err error
	// viper.SetConfigName(".env.yaml")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")   // optionally look for config in the working directory
	// err = viper.ReadInConfig() // Find and read the config file
	// if err != nil {            // Handle errors reading the config file
	// 	panic(fmt.Errorf("databases fatal error config file: %w", err))
	// }
	user := "warmsilver"
	password := "warmsilver"
	host := "localhost"
	port := "3306"
	database := "warmsilver"
	DB, _ = model.NewConnection(user, password, host, port, database)
	// if err != nil {
	// 	panic(err)
	// }
	DB.AutoMigrate(&core.Account{})
}
