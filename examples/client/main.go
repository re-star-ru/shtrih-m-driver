package main

import (
	"log"
	"time"

	"github.com/fess932/shtrih-m-driver/examples/client/rest"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("dial to kkt")

	kks, err := initKkts(confKKT{
		"EV-N": "10.51.0.73:7778",
		"SM-N": "10.51.0.74:7778",
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	service := rest.New(kks)
	service.Run()
}

type confKKT map[string]string

func initKkts(confs confKKT) (kks map[string]*kkt.KKT, err error) {
	kks = make(map[string]*kkt.KKT)

	for key, addr := range confs {
		kk, err := kkt.NewKKT(key, addr, time.Second*5, true)
		if err != nil {
			return nil, err
		}
		kks[key] = kk
	}

	return
}
