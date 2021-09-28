package main

import (
	"time"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
)

type confKKT map[string]ck

type ck struct {
	addr string // ip addres for kkt
	inn  string // default inn for fiscal operations
}

const defaultTimeout = time.Second * 120

func initKkts(confs confKKT) (kks map[string]*kkt.KKT, err error) {
	kks = make(map[string]*kkt.KKT)

	for key, c := range confs {
		kk, err := kkt.NewKKT(key, c.addr, c.inn, defaultTimeout, true)
		if err != nil {
			return nil, err
		}
		kks[key] = kk
	}

	return
}
