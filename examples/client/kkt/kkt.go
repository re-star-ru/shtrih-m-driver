package kkt

import (
	"fmt"
	"github.com/fess932/shtrih-m-driver/examples/client/kkt/transport"
	"github.com/looplab/fsm"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Messager interface {
	SendMessage(msg []byte) (resp []byte, err error)
}

type KKT struct {
	Organization string
	Place        string
	Addr         string
	CashierInn   string

	sync.Mutex
	m Messager
	d net.Dialer
	c net.Conn

	State    *fsm.FSM
	Substate *fsm.FSM
}

// NewKKT init new kkt device
func NewKKT(key, addr, inn string, connTimeout time.Duration, healthCheck bool) (kkt *KKT, err error) {
	s := strings.Split(key, "-")
	if len(s) != 2 {
		return nil, fmt.Errorf("неправильный ключ для ккт: %v", key)
	}

	kkt = &KKT{
		Organization: s[0],
		Place:        s[1],
		Addr:         addr,
		CashierInn:   inn,
	}

	kkt.d.Timeout = connTimeout
	kkt.State = newState()
	kkt.Substate = newSubstate()

	if healthCheck { // run healthcheck
		go func() {
			for {
				if err := kkt.Do(healhCheck); err != nil {
					log.Println(err)
				}
				time.Sleep(time.Second * 30)
			}
		}() // run healthcheck
	}

	return
}

func (kkt *KKT) connect() (err error) {
	const dialRetries = 3

	for i := 0; i < dialRetries; i++ {
		kkt.c, err = kkt.d.Dial("tcp", kkt.Addr)
		if err != nil {
			time.Sleep(time.Second * 1) // default timeout for retry
			continue
		}

		kkt.m = transport.New(kkt.c)
		return nil
	}

	return
}

// Do is function for starting request, create connection and close after exit
func (kkt *KKT) Do(cb func(kkt *KKT) (err error)) error {
	kkt.Lock()
	defer kkt.Unlock()

	t := time.Now()
	defer func(t time.Time) {
		log.Printf("kkt: %v, cmd time: %v", kkt.Addr, time.Since(t))
	}(t)

	if err := kkt.connect(); err != nil {
		err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.Addr, err)
		return err
	}
	defer func() {
		if e := kkt.c.Close(); e != nil {
			log.Println("deferred closing error:", e)
		}
	}()

	return cb(kkt)
}
