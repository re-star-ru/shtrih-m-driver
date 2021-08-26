package transport

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	STX byte = 0x02 // 2
	ENQ byte = 0x05 // 5
	ACK byte = 0x06 // 6
	NAK byte = 0x15 // 21
)

type K struct {
	controlByte byte
	conn        net.Conn
	sendMsgBuf  bytes.Buffer
}

func New(conn net.Conn) *K {
	return &K{
		controlByte: 0,
		conn:        conn,
		sendMsgBuf:  bytes.Buffer{},
	}
}

func (k *K) SendMessage(msg []byte) ([]byte, error) {
	k.sendMsgBuf.Write(msg)
	resp, err := k.sendENQ()
	k.sendMsgBuf.Reset()

	return resp, err
}

func (k *K) sendENQ() ([]byte, error) {
	k.writeByte(ENQ)
	k.readControlByte()

	switch k.controlByte {
	case ACK:
		log.Printf("ENQ ACK")
		msg, err := k.reciveMsg()
		if err != nil {
			log.Println(err)
		}
		log.Printf("received message in send enq: %x\n", msg)

		return k.sendENQ()
	case NAK:
		if err := k.sendMsg(); err != nil {
			log.Println(err)
		}

		return k.reciveMsg()
	default:
		log.Printf("wrong control byte %X, retry after sleep\n", k.controlByte)
		time.Sleep(time.Millisecond * 600)

		return k.sendENQ()
	}
}

func (k *K) sendMsg() error {
	msgLen := k.sendMsgBuf.Len()
	if msgLen > 255 {
		return fmt.Errorf("owerflow msg length: %v", msgLen)
	}

	var resp []byte
	resp = append(resp, byte(msgLen))
	resp = append(resp, k.sendMsgBuf.Bytes()...)
	crc := getChecksum(resp)

	// write datagram
	buf := bytes.NewBuffer(make([]byte, 4))
	buf.WriteByte(STX)
	buf.WriteByte(byte(msgLen))
	buf.Write(k.sendMsgBuf.Bytes())
	buf.WriteByte(crc)

	for i := 0; ; i++ {
		if _, err := k.conn.Write(buf.Bytes()); err != nil {
			log.Println("err in send msg loop:", err.Error())
		}

		k.readControlByte()
		switch k.controlByte {
		case ACK:
			return nil
		default:
			if i < 10 { // 10
				log.Printf("continue %v, ctrlByte: 0x%X \n", i, k.controlByte)
				continue
			}
			err := fmt.Errorf("wrong contol byte send message end %x", k.controlByte)
			return err
		}
	}
}

// read stx
// read len
func (k *K) reciveMsg() ([]byte, error) {
	if b := k.readByte(); b != STX {
		return nil, fmt.Errorf("wrong stx: %x", b)
	}

	lenMsg := k.readByte()
	msg := make([]byte, lenMsg)
	if _, err := k.conn.Read(msg); err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	lrc := k.readByte()

	k.writeByte(ACK) // write ack after read last byte in msg

	var resp []byte
	resp = append(resp, lenMsg)
	resp = append(resp, msg...)

	if !checksum(resp, lrc) {
		return nil, fmt.Errorf("invalid checksum %v", lrc)
	}

	return msg, nil
}

func checksum(data []byte, lrc byte) bool {
	return getChecksum(data) == lrc
}

func getChecksum(data []byte) (lrc byte) {
	for _, v := range data {
		lrc ^= v
	}

	return
}

func (k *K) readControlByte() {
	k.controlByte = k.readByte()
}

func (k *K) writeByte(b byte) {
	if _, err := k.conn.Write([]byte{b}); err != nil {
		log.Println("error write single byte:", err.Error())
	}
}

func (k *K) readByte() byte {
	b := make([]byte, 1)
	if _, err := k.conn.Read(b); err != nil {
		log.Println("error read single byte:", err.Error())
	}
	return b[0]
}
