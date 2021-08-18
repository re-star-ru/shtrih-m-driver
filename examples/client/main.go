package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fess932/shtrih-m-driver/examples/client/rest"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("dial to kkt")

	kks, err := initKkts([]confKKT{
		{"EV", "N", "10.51.0.73:7778"},
		{"SM", "N", "10.51.0.74:7778"},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	service := rest.New(kks...)
	service.Run()
}

type confKKT struct {
	Org   string
	Place string
	Addr  string
}

func initKkts(confs []confKKT) (kks []*kkt.KKT, err error) {
	//1

	for _, conf := range confs {
		if kktExist(kks, conf.Addr) {
			return nil, fmt.Errorf("kkt already exist: %v", conf)
		}
		kk := kkt.NewKKT(conf.Org, conf.Place, conf.Addr, time.Second*5, true)
		kks = append(kks, kk)
	}

	return
}

func kktExist(kks []*kkt.KKT, addr string) bool {
	for _, kk := range kks {
		if kk.Addr == addr {
			return true
		}
	}

	return false
}
