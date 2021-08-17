package kkt

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fess932/shtrih-m-driver/examples/client/commands"
	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/looplab/fsm"
)

const pingDeadline = time.Millisecond * 500

type KKT struct {
	Addr string
	d    net.Dialer
	conn net.Conn
	sync.Mutex

	ctrlByte []byte
	State    *fsm.FSM
	Substate *fsm.FSM
}

// send printCmd[]
// if specified send dontPrintOneCheck
// send writeCashierInn
// send closeCheck

func printCheckHandler(check *models.CheckPackage) func(kkt *KKT) error {
	return func(kkt *KKT) (err error) {
		log.Println("check:", check)

		if !kkt.canPrintCheck() { // check State
			err = fmt.Errorf("cant print check, wrong kkt State %v", kkt.State.Current())
			return
		}

		// validate input data
		o := commands.Operation{
			Type:    consts.Income,
			Subject: 0,
			Amount:  0,
			Price:   0,
			Sum:     0,
			Name:    "",
		}

		if err := o.Validate(); err != nil {
			return err
		}

		data, err := commands.CreateFNOperationV2(o)
		if err != nil {
			return err
		}

		log.Println("Cmd len ", len(data))
		log.Println("Data cmd create fn \n", hex.Dump(data))

		msg := createMessage(data)

		if err = sendMessage(kkt.conn, msg); err != nil {
			err = fmt.Errorf("kkt %s: send operation message error: %w", kkt.Addr, err)
			return
		}

		resp, err := kkt.receiveMessage()
		if err != nil {
			err = fmt.Errorf("kkt %s: revice operation message error: %w", kkt.Addr, err)
			return
		}

		if err = kkt.parseCmd(resp); err != nil {
			err = fmt.Errorf("kkt %s: parce recive operation message error: %w", kkt.Addr, err)
			return
		}

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

func NewKKT(addr string, connTimeout time.Duration, healthCheck bool) (kkt *KKT) {
	kkt = &KKT{}
	kkt.Addr = addr
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
