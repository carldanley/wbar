package main

import (
	"errors"
	_ "image/jpeg"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/carldanley/wbar/pkg"
	"github.com/joho/godotenv"
)

var ErrorNoBridgeHostSpecified = errors.New("please specify a wyze bridge host via the environment variable")

func main() {
	// create our channel for signal interrupts
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	// figure out where to load the .env file from
	envFileLocation := ".env"
	if location, exists := os.LookupEnv("ENV_FILE"); exists {
		envFileLocation = location
	}

	// attempt to load our .env file
	log.Printf("Loading environment file from location: %s\n", envFileLocation)

	err := godotenv.Load(envFileLocation)
	if err != nil {
		log.Fatalf("Could not load environment file from location: %v", err)
	}

	// make sure the host is specified
	if pkg.GetWyzeBridgeHost() == "" {
		panic(ErrorNoBridgeHostSpecified)
	}

	// start scanning for issues with time syncs
	go pkg.StartScanning()

	// wait for an interrupt signal
	<-signalChannel
}
