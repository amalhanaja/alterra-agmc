package main

import (
	"alterra-agmc-day-7/internal/app"
	"log"
)

func main() {
	if err := app.NewRestApiApp().Run(); err != nil {
		log.Fatalf("Error Running Application: %v", err)
	}
}
