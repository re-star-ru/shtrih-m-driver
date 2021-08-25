package kkt

import (
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

const pingDeadline = time.Second

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

// NewKKT init new kkt device
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
				if err := kkt.Do(healhCheck); err != nil {
					log.Println(err)
				}
				time.Sleep(time.Second * 5)
			}
		}() // run healthcheck
	}

	return
}

// Do is function for starting request, create connection and close after exit
func (kkt *KKT) Do(cb func(kkt *KKT) (err error)) error {
	kkt.Lock()
	defer kkt.Unlock()

	t := time.Now()
	defer func(t time.Time) {
		log.Println("cmd time:", time.Since(t))
	}(t)

	if err := kkt.dial(); err != nil {
		err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.Addr, err)
		return err
	}
	defer func() {
		if e := kkt.conn.Close(); e != nil {
			log.Println("deferred closing error:", e)
		}
	}()

	//if err := kkt.prepareRequest(); err != nil {
	//	err = fmt.Errorf("kkt %s : prepare request error: %w", kkt.Addr, err)
	//	return err
	//}

	return cb(kkt)
}

// send printCmd[]
// if specified send dontPrintOneCheck
// send writeCashierInn
// send closeCheck

func PrintCheckHandler(check models.CheckPackage) func(kkt *KKT) error {
	return func(kkt *KKT) (err error) {
		log.Println("check:", check)

		// check state
		if !kkt.canPrintCheck() { // check State
			err = fmt.Errorf("cant print check, wrong kkt State %v", kkt.State.Current())
			return
		}

		// set state not print one check if specified
		if check.NotPrint {
			if err = notPrintOneCheck(kkt); err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("not print ok")

		// add operationV2 to check
		for _, v := range check.Operations {
			if err = sendOperationsV2(kkt, v); err != nil {
				log.Println(err)
				return err
			}
			log.Println("send operationv2 ok")
		}

		return sendCloseCheckV2(kkt, check)
		// close check V2
		//data, err := commands.CreateFNCloseCheck(check)
		//if err != nil {
		//	return err
		//}
		//
		//msg := createMessage(data)
		//log.Println("msg: ", msg)
		//
		//resp, err := kkt.sendMessage(msg)
		//if err != nil {
		//	err = fmt.Errorf("kkt %s: send operationV2 message error: %w", kkt.Addr, err)
		//	return err
		//}
		//
		//if err = kkt.parseCmd(resp); err != nil {
		//	err = fmt.Errorf("kkt %s: parse recieve closeCheckV2 message error: %w", kkt.Addr, err)
		//	return
		//}
		//
		//return err
	}
}

func sendCloseCheckV2(kkt *KKT, check models.CheckPackage) error {
	// close check V2
	data, err := commands.CreateFNCloseCheck(check)
	if err != nil {
		return err
	}

	msg := createMessage(data)
	log.Println("msg: ", msg)

	return sendENQ(kkt, msg)

	//resp, err := kkt.sendMessage(msg)
	//if err != nil {
	//	err = fmt.Errorf("kkt %s: send operationV2 message error: %w", kkt.Addr, err)
	//	return err
	//}
	//
	//if err = kkt.parseCmd(resp); err != nil {
	//	err = fmt.Errorf("kkt %s: parse recieve closeCheckV2 message error: %w", kkt.Addr, err)
	//	return
	//}
}

func sendOperationsV2(kkt *KKT, o models.Operation) error {
	data, err := commands.CreateFNOperationV2(o)
	if err != nil {
		return err
	}

	msg := createMessage(data)

	return sendENQ(kkt, msg)

	//resp, err := kkt.sendMessage(msg)
	//if err != nil {
	//	err = fmt.Errorf("kkt %s: send operationV2 message error: %w", kkt.Addr, err)
	//	return err
	//}
	//
	//if err = kkt.parseCmd(resp); err != nil {
	//	err = fmt.Errorf("kkt %s: parse recive operationV2 message error: %w", kkt.Addr, err)
	//	return err
	//}
	//
	//return nil
}

func healhCheck(kkt *KKT) (err error) {
	msg := createMessage(commands.CreateShortStatus())
	//resp, err := kkt.sendMessage(msg)

	return sendENQ(kkt, msg)
	//
	//if err != nil {
	//	err = fmt.Errorf("kkt %s : send message error: %w", kkt.Addr, err)
	//	return
	//}
	//
	//if err = kkt.parseCmd(resp); err != nil {
	//	log.Println("error while parsing response command:", err)
	//	return
	//}
	//
	//return nil
}

func notPrintOneCheck(kkt *KKT) (err error) {
	msg := createMessage(commands.CreateNotPrintOneCheck())

	return sendENQ(kkt, msg)

	//resp, err := kkt.sendMessage(msg)
	//if err != nil {
	//	err = fmt.Errorf("kkt %s : send message error: %w", kkt.Addr, err)
	//	return
	//}
	//
	//if err = kkt.parseCmd(resp); err != nil {
	//	log.Println("error while parsing response command:", err)
	//	return
	//}
	//
	//return nil
}
