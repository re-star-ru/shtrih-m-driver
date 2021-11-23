package kkt

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt/transport"
	"github.com/re-star-ru/shtrih-m-driver/app/models/kkterrors"
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

// NewKKT init new kkt device.
func NewKKT(key, addr, inn string, connTimeout time.Duration, healthCheck bool) (*KKT, error) {
	s := strings.Split(key, "-")
	if len(s) != 2 {
		return nil, fmt.Errorf("неправильный ключ для ккт: %v", key)
	}

	kkt := &KKT{
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
				if err := kkt.Do(doHealhCheck); err != nil {
					log.Err(err).Send()
				}

				time.Sleep(time.Second * 30)
			}
		}() // run healthcheck
	}

	return kkt, nil
}

func (kkt *KKT) connect(ctx context.Context) (err error) {
	kkt.c, err = kkt.d.DialContext(ctx, "tcp", kkt.Addr)
	if err != nil {
		return err
	}
	kkt.m = transport.New(kkt.c)
	return nil
}

// Do is function for starting request, create connection and close after exit
// Handle context right here.
func (kkt *KKT) Do(cb func(kkt *KKT) (err error)) error {
	kkt.Lock()
	defer func() {
		if e := kkt.c.Close(); e != nil {
			log.Err(e).Msg("deferred closing error:")
		}
		kkt.Unlock()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), kkt.d.Timeout)

	t := time.Now()

	defer cancel()
	defer func(t time.Time) {
		log.Printf("kkt: %v, cmd time: %v", kkt.Addr, time.Since(t))
	}(t)

	ch := make(chan error)

	go func(ctx context.Context, ch chan error, cb func(kkt *KKT) (err error)) {
		if err := kkt.connect(ctx); err != nil {
			err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.Addr, err)
			ch <- err
			return
		}

		ch <- cb(kkt)

		for {
			select {
			case <-ctx.Done():
				log.Print("TIMEOUT WITH CONTEXT!")
				ch <- kkterrors.ErrTimeout
				return
			}
		}

	}(ctx, ch, cb)

	return <-ch
}

//
//func (kkt *KKT) goDo(ctx context.Context, ch chan error, cb func(kkt *KKT) (err error)) {
//
//	if err := kkt.connect(); err != nil {
//		err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.Addr, err)
//		ch <- err
//	}
//
//	ch <- cb(kkt)
//}
