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
)

type MessageSendCloser interface {
	SendMessage(msg []byte) (resp []byte, err error)
	Close() error
}

type KKT struct {
	Organization string
	Place        string
	Addr         string
	CashierInn   string

	// t is timeout for cmd
	t time.Duration

	sync.Mutex
	m MessageSendCloser
	d net.Dialer

	State    *fsm.FSM
	Substate *fsm.FSM
}

// NewKKT init new kkt device.
func NewKKT(key, addr, inn string, timeout time.Duration) (*KKT, error) {
	s := strings.Split(key, "-")
	if len(s) != 2 {
		return nil, fmt.Errorf("неправильный ключ для ккт: %v", key)
	}

	kkt := &KKT{
		Organization: s[0],
		Place:        s[1],
		Addr:         addr,
		CashierInn:   inn,
		t:            timeout,
	}

	kkt.State = newState()
	kkt.Substate = newSubstate()

	//if healthCheck { // run healthcheck
	//	go func() {
	//		for {
	//			if err := kkt.Do(doHealhCheck); err != nil {
	//				log.Err(err).Msg("healthcheck error")
	//			}
	//
	//			time.Sleep(time.Second * 30)
	//		}
	//	}() // run healthcheck
	//}

	return kkt, nil
}

func (kkt *KKT) connect(ctx context.Context) error {
	kkt.Lock()
	c, err := kkt.d.DialContext(ctx, "tcp", kkt.Addr)
	if err != nil {
		return err
	}
	kkt.m = transport.New(c)
	return nil
}

func (kkt *KKT) close() (err error) {
	defer kkt.Unlock()
	if kkt.m != nil {
		err = kkt.m.Close()
	}

	return err
}

// Do is function for starting request, create connection and close after exit
// Handle context right here.
func (kkt *KKT) Do(cb func(kkt *KKT) (err error)) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), kkt.t)
	t := time.Now()

	defer func(t time.Time) {
		if e := kkt.close(); e != nil {
			log.Err(e).Msg("kkt close error")
		}
		cancel()
		log.Printf("kkt: %v, cmd time: %v", kkt.Addr, time.Since(t))
	}(t)

	if err = kkt.connect(ctx); err != nil {
		err = fmt.Errorf("kkt %s: %w", kkt.Addr, err)
		return err
	}

	ch := make(chan error)
	go func(chan<- error) {
		ch <- cb(kkt)
	}(ch)

	select {
	case <-ctx.Done():
		log.Debug().Msg("Context deadline, close conn")
	case err = <-ch:
		if err != nil {
			log.Err(err).Msg("cmd done with error")
		}
	}

	return err
}
