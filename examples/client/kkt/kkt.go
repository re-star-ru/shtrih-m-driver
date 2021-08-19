package kkt

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/fess932/shtrih-m-driver/examples/client/commands"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/looplab/fsm"
)

const pingDeadline = time.Millisecond * 500

type KKT struct {
	Organization string
	Place        string
	Addr         string
	d            net.Dialer
	conn         net.Conn
	sync.Mutex

	ctrlByte []byte
	State    *fsm.FSM
	Substate *fsm.FSM
}

// send printCmd[]
// if specified send dontPrintOneCheck
// send writeCashierInn
// send closeCheck

func PrintCheckHandler(check models.CheckPackage) func(kkt *KKT) error {
	return func(kkt *KKT) (err error) {
		log.Println("check:", check)

		if !kkt.canPrintCheck() { // check State
			err = fmt.Errorf("cant print check, wrong kkt State %v", kkt.State.Current())
			return
		}

		for _, v := range check.Operations {
			data, err := commands.CreateFNOperationV2(v)
			if err != nil {
				return err
			}

			log.Println("Cmd len ", len(data))
			log.Println("Data cmd create fn \n", hex.Dump(data))

			msg := createMessage(data)
			log.Println("msg: ", msg)
		}

		//if err = sendMessage(kkt.conn, msg); err != nil {
		//	err = fmt.Errorf("kkt %s: send operation message error: %w", kkt.Addr, err)
		//	return
		//}
		//
		//resp, err := kkt.receiveMessage()
		//if err != nil {
		//	err = fmt.Errorf("kkt %s: revice operation message error: %w", kkt.Addr, err)
		//	return
		//}
		//
		//if err = kkt.parseCmd(resp); err != nil {
		//	err = fmt.Errorf("kkt %s: parce recive operation message error: %w", kkt.Addr, err)
		//	return
		//}

		return nil
	}
}

func healhCheck(kkt *KKT) (err error) {
	msg := createMessage(commands.CreateShortStatus())

	if err = sendMessage(kkt.conn, msg); err != nil {
		err = fmt.Errorf("kkt %s : send message error: %w", kkt.Addr, err)
		return
	}

	resp, err := kkt.receiveMessage()
	if err != nil {
		err = fmt.Errorf("kkt %s : receive message error: %w", kkt.Addr, err)
		return
	}

	if err = kkt.parseCmd(resp); err != nil {
		log.Println("error while parsing response command:", err)
		return
	}

	return nil
}

func (kkt *KKT) canPrintCheck() bool {
	return kkt.State.Can(printCheck) && kkt.Substate.Can(printCheck)
}

func NewKKT(key, addr string, connTimeout time.Duration, healthCheck bool) (kkt *KKT, err error) {
	s := strings.Split(key, "-")
	if len(s) != 2 {
		return nil, fmt.Errorf("неправильный ключ для ккт: %v", key)
	}

	kkt = &KKT{
		Organization: s[0],
		Place:        s[1],
		Addr:         addr,
	}

	kkt.d.Timeout = connTimeout
	kkt.ctrlByte = make([]byte, 1)
	kkt.State = newState()
	kkt.Substate = newSubstate()

	if healthCheck { // run healthcheck
		go func() {
			for {
				t := time.Now()
				if err := kkt.Do(healhCheck); err != nil {
					log.Println(err)
				}

				log.Println("cmd time:", time.Since(t))
				time.Sleep(time.Second * 5)
			}
		}() // run healthcheck
	}

	return
}

func (kkt *KKT) Do(cb func(kkt *KKT) (err error)) (err error) {
	kkt.Lock()
	defer kkt.Unlock()

	if err = kkt.dial(); err != nil {
		err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.Addr, err)
		return
	}
	defer func() {
		if err := kkt.conn.Close(); err != nil {
			log.Println("deferred closing error:", err)
		}
	}()

	if err = kkt.prepareRequest(); err != nil {
		err = fmt.Errorf("kkt %s : prepare request error: %w", kkt.Addr, err)
		return
	}

	return cb(kkt)
}
