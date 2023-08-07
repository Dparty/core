package services

import (
	"fmt"

	"github.com/Dparty/common/cloud"
	"github.com/Dparty/model"
	"github.com/Dparty/model/core"
	"github.com/Dparty/model/restaurant"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

var DB *gorm.DB

var CosClient cloud.CosClient
var Bucket string

func init() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")   // optionally look for config in the working directory
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	Bucket = viper.GetString("cos.Bucket")
	CosClient.Region = viper.GetString("cos.Region")
	CosClient.SecretID = viper.GetString("cos.SecretID")
	CosClient.SecretKey = viper.GetString("cos.SecretKey")
}

func init() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")   // optionally look for config in the working directory
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	DB, err = model.NewConnection(user, password, host, port, database)
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&restaurant.Item{})
	DB.AutoMigrate(&restaurant.Restaurant{})
	DB.AutoMigrate(&core.Account{})
}
