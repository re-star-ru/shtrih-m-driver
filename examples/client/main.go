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

	// client code
	kk := kkt.NewKKT("10.51.0.73:7778", time.Second*5, true)
	kk2 := kkt.NewKKT("10.51.0.74:7778", time.Second*5, true)

	service := rest.New(kk, kk2)
	service.Run()
}
