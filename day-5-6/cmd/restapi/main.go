package main

import "alterra-agmc-day-5-6/internal/app"

func main() {
	if err := app.NewRestApiApp().Run(); err != nil {
		panic("error running rest api application")
	}
}
