package main

import (
	//"log"
	// _ "net/http/pprof"

	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
	"github.com/re-star-ru/shtrih-m-driver/app/rest"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	//go func() {
	//	log.Fatal(http.ListenAndServe(":8090", nil))
	//}()

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

type confKKT map[string]ck

type ck struct {
	addr string
	inn  string
}

func initKkts(confs confKKT) (kks map[string]*kkt.KKT, err error) {
	kks = make(map[string]*kkt.KKT)

	for key, c := range confs {
		kk, err := kkt.NewKKT(key, c.addr, c.inn, time.Second*10, true)
		if err != nil {
			return nil, err
		}
		kks[key] = kk
	}

	return
}
