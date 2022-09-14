package main

import (
	"alterra-agmc-day-2/config"
	"alterra-agmc-day-2/routes"
)

func main() {
	config.InitDB()

	e := routes.New()

	e.Logger.Fatal(e.Start(":8080"))

}
