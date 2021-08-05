package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/fess932/shtrih-m-driver/examples/client/commands"
	"github.com/fess932/shtrih-m-driver/pkg/consts"

	"github.com/go-chi/chi/v5"

	"github.com/looplab/fsm"
)

// Control bytes
const (
	ENQ byte = 0x05
	STX byte = 0x02
	ACK byte = 0x06 // 6
	NAK byte = 0x15 // 21
)

const pingDeadline = time.Millisecond * 500

var errChecksum = errors.New("checksum does not match")

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("dial to kkt")

	// client code
	kkt := newKKT("10.51.0.73:7778", time.Second*5, true)

	if err := kkt.Do(printCheckHandler); err != nil {
		log.Println(err)
	}

	{ // http handler
		r := chi.NewRouter()

		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			if _, err := fmt.Fprintf(w,
				"state kkt: %v, substate: %v\ncan print: %v\n",
				kkt.state.Current(), kkt.substate.Current(), kkt.canPrintCheck()); err != nil {
				log.Println(err)
			}
		})

		r.Get("/print", func(w http.ResponseWriter, r *http.Request) {
			if err := kkt.Do(printCheckHandler); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}
		})
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}

type KKT struct {
	addr string
	d    net.Dialer
	conn net.Conn
	sync.Mutex

	ctrlByte []byte
	state    *fsm.FSM
	substate *fsm.FSM
}

func printCheckHandler(kkt *KKT) (err error) {
	if !kkt.canPrintCheck() { // check state
		err = fmt.Errorf("cant print check, wrong kkt state %v", kkt.state.Current())
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
		err = fmt.Errorf("kkt %s: send operation message error: %w", kkt.addr, err)
		return
	}

	resp, err := kkt.receiveMessage()
	if err != nil {
		err = fmt.Errorf("kkt %s: revice operation message error: %w", kkt.addr, err)
		return
	}

	if err = kkt.parseCmd(resp); err != nil {
		err = fmt.Errorf("kkt %s: parce recive operation message error: %w", kkt.addr, err)
		return
	}

	return nil
}

func healhCheck(kkt *KKT) (err error) {
	msg := createMessage(commands.CreateShortStatus())

	if err = sendMessage(kkt.conn, msg); err != nil {
		err = fmt.Errorf("kkt %s : send message error: %w", kkt.addr, err)
		return
	}

	resp, err := kkt.receiveMessage()
	if err != nil {
		err = fmt.Errorf("kkt %s : receive message error: %w", kkt.addr, err)
		return
	}

	if err = kkt.parseCmd(resp); err != nil {
		log.Println("error while parsing response command:", err)
		return
	}

	return nil
}

func (kkt *KKT) canPrintCheck() bool {
	return kkt.state.Can(printCheck) && kkt.substate.Can(printCheck)
}

func newKKT(addr string, connTimeout time.Duration, healthCheck bool) (kkt *KKT) {
	kkt = &KKT{}
	kkt.addr = addr
	kkt.d.Timeout = connTimeout
	kkt.ctrlByte = make([]byte, 1)
	kkt.state = newState()
	kkt.substate = newSubstate()

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

// Events
const (
	printCheck  = "printCheck"
	shiftOpen   = "shiftOpen"
	shiftClose  = "shiftClose"
	shiftReopen = "shiftReopen"
)

func newState() *fsm.FSM {
	return fsm.NewFSM(
		"",
		fsm.Events{
			{Name: printCheck, Src: []string{"shiftOpen"}, Dst: "shiftOpen"},
			{Name: shiftOpen, Src: []string{"shiftClosed"}, Dst: "shiftOpen"},
			{Name: shiftClose, Src: []string{"shiftOpen"}, Dst: "shiftClosed"},
			{Name: shiftReopen, Src: []string{"shiftExpired"}, Dst: "shiftClosed"},
		},
		fsm.Callbacks{},
	)
}

func newSubstate() *fsm.FSM {
	return fsm.NewFSM(
		"",
		fsm.Events{{
			Name: printCheck, Src: []string{"paperLoaded"}, Dst: "paperLoaded",
		}},
		fsm.Callbacks{},
	)
}

func (kkt *KKT) Do(cb func(kkt *KKT) (err error)) (err error) {
	kkt.Lock()
	defer kkt.Unlock()

	if err = kkt.dial(); err != nil {
		err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.addr, err)
		return
	}
	defer func() {
		if err := kkt.conn.Close(); err != nil {
			log.Println("deferred closing error:", err)
		}
	}()

	if err = kkt.prepareRequest(); err != nil {
		err = fmt.Errorf("kkt %s : prepare request error: %w", kkt.addr, err)
		return
	}

	return cb(kkt)
}

func (kkt *KKT) dial() (err error) {
	const dialRetries = 3

	for i := 0; i < dialRetries; i++ {
		kkt.conn, err = kkt.d.Dial("tcp", kkt.addr)
		if err == nil {
			return
		}

		time.Sleep(time.Second * 1) // default timeout for retry
	}

	return
}

func (kkt *KKT) prepareRequest() (err error) {
	const pingRetries = 3

	for i := 0; i < pingRetries; i++ {
		if i != 0 {
			time.Sleep(time.Second * 1)
		}

		if err = kkt.conn.SetDeadline(time.Now().Add(pingDeadline)); err != nil {
			continue
		}

		// ////// write
		_, err := kkt.conn.Write([]byte{ENQ})
		if err != nil {
			continue
		}

		// ////// read
		_, err = kkt.conn.Read(kkt.ctrlByte)
		if err != nil {
			continue
		}

		switch kkt.ctrlByte[0] {
		case ACK:
			err = kkt.sendACK()
			if err != nil {
				continue
			}
		case NAK:
			return nil
		default:
			continue
		}
	}

	return err
}

func (kkt *KKT) sendACK() error {
	_, err := kkt.conn.Write([]byte{ACK})
	return err
}

func (kkt *KKT) receiveMessage() (message []byte, err error) {
	err = kkt.readACK()
	if err != nil {
		err = fmt.Errorf("err while read ACK: %w", err)
		return
	}

	err = kkt.readSTX()
	if err != nil {
		err = fmt.Errorf("err while read STX: %w", err)
		return
	}

	var l byte
	l, err = kkt.readLen()
	if err != nil {
		err = fmt.Errorf("err while read len: %w", err)
		return
	}

	msg, err := kkt.readMessage(l)
	if err != nil {
		err = fmt.Errorf("err while read message: %w", err)
		return
	}

	return msg, nil
}

func (kkt *KKT) readACK() error {
	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return err
	}
	if kkt.ctrlByte[0] != ACK {
		return fmt.Errorf("got wrong control byte: %v, expect: %v", kkt.ctrlByte[0], ACK)
	}

	return nil
}

func (kkt *KKT) readSTX() error {
	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return err
	}
	if kkt.ctrlByte[0] != STX {
		return fmt.Errorf("got wrong control byte: %v, expect: %v", kkt.ctrlByte[0], STX)
	}

	return nil
}

func (kkt *KKT) readLen() (byte, error) {
	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return 0, err
	}

	return kkt.ctrlByte[0], nil
}

func (kkt *KKT) readMessage(messageLen byte) ([]byte, error) {
	msg := make([]byte, messageLen)
	if _, err := kkt.conn.Read(msg); err != nil {
		return nil, err
	}

	rLrc, err := kkt.readLRC()
	if err != nil {
		return nil, err
	}

	var resp []byte
	resp = append(resp, messageLen)
	resp = append(resp, msg...)

	lrc := computeLRC(resp)

	if lrc != rLrc {
		return nil, errChecksum
	}

	return msg, nil
}

func (kkt *KKT) readLRC() (byte, error) {
	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return 0, err
	}

	return kkt.ctrlByte[0], nil
}

func createMessage(cmdData []byte) []byte {
	l := byte(len(cmdData)) // may be panic if overflowed? cannot be more than 255
	m := []byte{STX, l}
	m = append(m, cmdData...)
	m = append(m, computeLRC(m[1:]))

	return m
}

func sendMessage(w io.Writer, message []byte) error {
	_, err := w.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func computeLRC(data []byte) (lrc byte) {
	for _, v := range data {
		lrc ^= v
	}

	return
}
