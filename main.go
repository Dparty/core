package main

import (
	"fmt"

	"github.com/Dparty/core/controllers"
	"github.com/Dparty/core/services"
)

func main() {
	// var err error
	err := services.UpdatePasswordForce(1681137588590612482, "hello")
	fmt.Println(err)
	// viper.SetConfigName(".env") // name of config file (without extension)
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")   // optionally look for config in the working directory
	// err = viper.ReadInConfig() // Find and read the config file
	// if err != nil {            // Handle errors reading the config file
	// 	panic(fmt.Errorf("databases fatal error config file: %w", err))
	// }
	// port := ":" + viper.GetString("server.port")
	controllers.Init()
	controllers.Run(":8080")
}
