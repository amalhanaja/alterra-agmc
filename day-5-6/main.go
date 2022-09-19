package main

import (
	"alterra-agmc-day-5-6/config"
	"alterra-agmc-day-5-6/lib/validator"
	"alterra-agmc-day-5-6/routes"
)

func main() {
	config.InitDB()

	e := routes.New()
	e.Validator = validator.NewCustomValidator()

	e.Logger.Fatal(e.Start(":8080"))

}
