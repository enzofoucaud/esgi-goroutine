package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// init logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	switch os.Args[1] {
	case "resize":
		resizeImages()
	case "dining":
		diningPhilosophers()
	default:
		log.Error().Msg("invalid argument")
		os.Exit(1)
	}
}
