package transport

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog/log"
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
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	k.sendMsgBuf.Reset()

	return resp, err
}

func (k *K) sendENQ() ([]byte, error) {
	k.writeByte(ENQ)
	if err := k.readControlByte(); errors.Is(err, io.EOF) {
		log.Err(err).Send()
		return nil, err
	}

	switch k.controlByte {
	case ACK:
		log.Printf("ENQ ACK")
		msg, err := k.reciveMsg()
		if err != nil {
			log.Err(err).Send()
		}
		log.Printf("received message in send enq: %x\n", msg)

		return k.sendENQ()
	case NAK:
		if err := k.sendMsg(); err != nil {
			log.Print(err)
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
			log.Err(err).Msg("err in send msg loop")
		}

		if err := k.readControlByte(); errors.Is(err, io.EOF) {
			log.Err(err).Msg("err in read control byte in msg loop")
		}

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
	if b, _ := k.readByte(); b != STX {
		return nil, fmt.Errorf("wrong stx: %x", b)
	}

	lenMsg, _ := k.readByte()
	msg := make([]byte, lenMsg)
	if _, err := k.conn.Read(msg); err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	lrc, _ := k.readByte()

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

func (k *K) readControlByte() (err error) {
	k.controlByte, err = k.readByte()
	return
}

func (k *K) writeByte(b byte) {
	if _, err := k.conn.Write([]byte{b}); err != nil {
		log.Print("error write single byte:", err.Error())
	}
}

func (k *K) readByte() (byte, error) {
	b := make([]byte, 1)

	_, err := k.conn.Read(b)

	if err != nil {
		log.Err(err).Msg("error read single byte")
	}

	if errors.Is(err, io.EOF) {
		return 0, io.EOF
	}

	return b[0], nil
}
