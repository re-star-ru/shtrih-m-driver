package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/rest"
)

var rev = 1

func main() {
	log.Logger = log.Logger.With().
		Caller().
		Logger()
	zerolog.TimeFieldFormat = time.StampMilli

	kks, err := initKkts(confKKT{
		"EV-S": ck{"10.51.0.71:7778", "263209745357"},
		"SM-S": ck{"10.51.0.72:7778", "262804786800"},

		"EV-N": ck{"10.51.0.73:7778", "263209745357"},
		"SM-N": ck{"10.51.0.74:7778", "262804786800"},
	})

	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	// addr
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080"
	}
	//
	service := rest.New(kks, addr)
	service.Run()
}
