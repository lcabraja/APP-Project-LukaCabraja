package main

import (
	"log"
	"os"
	"strings"

	"github.com/lcabraja/APP-Project-LukaCabraja/configuration"
	"github.com/lcabraja/APP-Project-LukaCabraja/instance"

	"github.com/joho/godotenv"
)

func main() {
	if loadenv, _ := os.LookupEnv("LOAD_ENV"); strings.ToLower(loadenv) != "false" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if clear, _ := os.LookupEnv("CLEAR_CONSOLE"); strings.ToLower(clear) != "false" {
		configuration.ClearConsole()
	}

	c := (&configuration.Configuration{}).ApplyDefaultsOnEmpty()

	app := instance.NewApp(c, true)

	app.Start()
}
