package main

import "log"

const (
	RECV = ACK
	SEND = NAK
)

const (
	STX byte = 0x02 // 2
	ENQ byte = 0x05 // 5
	ACK byte = 0x06 // 6
	NAK byte = 0x15 // 21
)

const MSG = 2

func processSend() {
	kt := &k{}

	kt.send(ENQ)

	switch kt.state {
	case ACK:
		kt.recv()
	case NAK:
		kt.send(MSG)
		kt.recvCtrl()
		switch kt.ctrl {
		case ACK:
			log.Println("OK")
			kt.recv()
			if kt.err != nil {
				log.Println("WRONG MSG", kt.err)
				return
			}

		case NAK:
			log.Println("ERR")
		}
	default:
		log.Println("RETRY IF HAS TRIES")
	}
}

type k struct {
	state byte
	ctrl  byte
	err   error
}

func (k *k) send(b byte) {

}
func (k *k) recv()     {}
func (k *k) recvCtrl() {}
