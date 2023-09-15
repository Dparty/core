package main

import (
	"github.com/Dparty/core/controllers"
)

func main() {
	controllers.Init()
	controllers.Run(":8080")
}
