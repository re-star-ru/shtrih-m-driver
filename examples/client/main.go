package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Control bytes
const (
	ENQ byte = 0x05
	STX byte = 0x02
	ACK byte = 0x06 // 6
	NAK byte = 0x15 // 21
)

const (
	SHORT_STATUS byte = 0x10
)

const pingDeadline = time.Millisecond * 500

var ErrChecksum error = errors.New("checksum does not match")

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("dial to kkt")

	kkt := newKKT("10.51.0.71:7778", time.Second*5)
	msg := createMessage(SHORT_STATUS, []byte{0x1E, 0x00, 0x00, 0x00})

	resp, err := kkt.SendRequest(msg)
	if err != nil {
		log.Println("error while send message:", err)
		return
	}

	log.Println("responded message:", hex.Dump(resp))
}

type KKT struct {
	addr string
	d    net.Dialer
	conn net.Conn
	sync.Mutex

	ctrlByte []byte
}

func newKKT(addr string, connTimeout time.Duration) (kkt *KKT) {
	kkt = &KKT{}
	kkt.addr = addr
	kkt.d.Timeout = connTimeout
	kkt.ctrlByte = make([]byte, 1)

	return
}

//func (kkt *KKT) connect() {
//	kkt.Lock()
//	defer kkt.Unlock()
//
//	var err error
//
//	msg := createMessage(SHORT_STATUS, []byte{0x1E, 0x00, 0x00, 0x00}) // short status message
//
//	if err := kkt.SendMessage(msg); err != nil {
//		log.Println("err while send message: ")
//	}
//
//	kkt.conn, err = kkt.d.Dial("tcp", kkt.addr)
//	if err != nil {
//		log.Println("err dial dial:", err)
//		return
//	}
//	defer func() {
//		if err := kkt.conn.Close(); err != nil {
//			log.Println("Error while closing conn in defer", err)
//		}
//	}()
//	log.Println("dial ok")
//	canSendCmd, err := ping(kkt.conn)
//	if err != nil {
//		log.Println("No connection:", err)
//		log.Println("Retry...")
//		return
//	}
//	log.Println("Can send cmd conn to kkt:", canSendCmd)
//
//	// if you get error continue, if ok break
//	defaultTries := 5
//	for i := 0; i < defaultTries; i++ {
//		if err := sendMessage(kkt.conn, msg); err != nil {
//			log.Println("err on try", i, "err:", err)
//			continue
//		}
//
//		if kkt.awaitACK() {
//			log.Println("ack recived, ok")
//		}
//
//		break
//	}
//
//	if err != nil {
//		log.Println("err while send cmd:", err)
//		return
//	}
//}

func (kkt *KKT) SendRequest(req []byte) (resp []byte, err error) {
	kkt.Lock()
	defer kkt.Unlock()

	// retry for message?

	// dial conn
	if err = kkt.dial(); err != nil {
		err = fmt.Errorf("kkt %s : dial: no connection: %w", kkt.addr, err)
		return
	}
	defer func() {
		if err := kkt.conn.Close(); err != nil {
			log.Println("deferred closing error:", err)
		}
	}()
	// end dial

	if err = kkt.prepareRequest(); err != nil {
		err = fmt.Errorf("kkt %s : prepare request error: %w", kkt.addr, err)
		return
	}

	if err = sendMessage(kkt.conn, req); err != nil {
		err = fmt.Errorf("kkt %s : send message error: %w", kkt.addr, err)
		return
	}

	resp, err = kkt.receiveMessage()
	if err != nil {
		err = fmt.Errorf("kkt %s : receive message error: %w", kkt.addr, err)
		return
	}

	return
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
		n, err := kkt.conn.Write([]byte{ENQ})
		if err != nil {
			continue
		}
		log.Println("bytes writed to conn:", n)

		// ////// read
		n, err = kkt.conn.Read(kkt.ctrlByte)
		if err != nil {
			continue
		}
		log.Println("bytes readed from conn:", n)

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
	err = kkt.readNAK()
	if err != nil {
		err = fmt.Errorf("err while read NAK: %w", err)
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

func (kkt *KKT) readNAK() error {
	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return err
	}
	if kkt.ctrlByte[0] != NAK {
		return fmt.Errorf("got wrong control byte: %v, expect: %v", kkt.ctrlByte[0], NAK)
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
		return nil, ErrChecksum
	}

	return msg, nil
}

func (kkt *KKT) readLRC() (byte, error) {
	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return 0, err
	}

	return kkt.ctrlByte[0], nil
}

// Служебное сообщение
func ping(conn net.Conn) (bool, error) {
	if err := conn.SetDeadline(time.Now().Add(pingDeadline)); err != nil {
		log.Fatal(err)
	}

	n, err := conn.Write([]byte{ENQ})
	// write enq for req message
	if err != nil {
		return false, nil
	}
	log.Println("bytes writed to conn:", n)

	b := make([]byte, 1)
	n, err = conn.Read(b)
	if err != nil {
		return false, err
	}

	log.Println("readed bytes from buf:", n)

	switch b[0] {
	case NAK:
		return true, nil
	case ACK:
		return false, fmt.Errorf("kkt busy")
	default:
		return false, fmt.Errorf("wrong byte for ENQ")
	}
}

func createMessage(cmdID byte, cmdData []byte) []byte {
	l := byte(len(cmdData)) // may be panic if overflow?
	m := []byte{STX, l, cmdID}

	m = append(m, cmdData...)
	m = append(m, computeLRC(m[1:]))

	log.Println(len(m))

	return m
}

func sendMessage(w io.Writer, message []byte) error {
	n, err := w.Write(message)
	if err != nil {
		return err
	}
	log.Println("bytes writed:", n)

	return nil
}

func computeLRC(data []byte) (lrc byte) {
	for _, v := range data {
		lrc ^= v
	}

	return
}
