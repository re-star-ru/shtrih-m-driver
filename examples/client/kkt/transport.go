package kkt

import (
	"fmt"
	"log"
	"time"
)

// Control bytes
const (
	STX byte = 0x02 // 2
	ENQ byte = 0x05 // 5
	ACK byte = 0x06 // 6
	NAK byte = 0x15 // 21
)

func (kkt *KKT) dial() (err error) {
	const dialRetries = 3

	for i := 0; i < dialRetries; i++ {
		kkt.conn, err = kkt.d.Dial("tcp", kkt.Addr)
		if err == nil {
			return
		}

		time.Sleep(time.Second * 1) // default timeout for retry
	}

	return
}

func sendENQ(kkt *KKT, msg []byte) error {
	if err := kkt.conn.SetDeadline(time.Now().Add(pingDeadline)); err != nil {
		return err
	}

	if _, err := kkt.conn.Write([]byte{ENQ}); err != nil {
		return err
	}

	if _, err := kkt.conn.Read(kkt.ctrlByte); err != nil {
		return err
	}

	switch kkt.ctrlByte[0] {
	case ACK:
		// read message
		resp, err := kkt.receiveMessage()
		if err != nil {
			err = fmt.Errorf("kkt %s : err while receive message error: %w", kkt.Addr, err)
			return err
		}
		if err = kkt.parseCmd(resp); err != nil {
			err = fmt.Errorf("kkt %s : err while parse response: %w", kkt.Addr, err)
			return err
		}
	case NAK:
		// send message
		resp, err := kkt.sendMessage(msg)
		if err != nil {
			err = fmt.Errorf("kkt %s : send message error: %w", kkt.Addr, err)
			log.Println(err)
			return err
		}
		if err = kkt.parseCmd(resp); err != nil {
			log.Println("error while parsing response command:", err)
			return err
		}

	default:
		// wait
	}

	return nil
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

		//  write
		_, err := kkt.conn.Write([]byte{ENQ})
		if err != nil {
			continue
		}

		//  read
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

	err = kkt.sendACK()
	if err != nil {
		err = fmt.Errorf("err while send ack: %w", err)
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

func (kkt *KKT) sendMessage(message []byte) (resp []byte, err error) {
	_, err = kkt.conn.Write(message)
	if err != nil {
		return nil, err
	}

	return kkt.receiveMessage()
}

func computeLRC(data []byte) (lrc byte) {
	for _, v := range data {
		lrc ^= v
	}

	return
}
