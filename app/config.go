package main

import (
	"fmt"
	"time"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
)

type confKKT map[string]ck

type ck struct {
	addr string // ip address for kkt
	inn  string // default inn for fiscal operations
}

const defaultTimeout = time.Second * 120

func initKkts(confs confKKT) (kks map[string]*kkt.KKT, err error) {
	kks = make(map[string]*kkt.KKT)

	for key, c := range confs {
		kk, err := kkt.NewKKT(key, c.addr, c.inn, defaultTimeout, true)
		if err != nil {
			return nil, fmt.Errorf("init kkt: %w", err)
		}

		kks[key] = kk
	}

	return
}
