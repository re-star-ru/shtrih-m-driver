package transport

import (
	"bytes"
	"errors"
	"fmt"
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

var (
	ErrInvalidChecksum   = errors.New("invalid checksum")
	ErrWrongSTX          = errors.New("wrong stx")
	ErrOverflowMsgLength = errors.New("overflow msg length")
	ErrWrongControlByte  = errors.New("wrong control byte")
)

type K struct {
	controlByte byte
	c           net.Conn
	sendMsgBuf  bytes.Buffer
}

func New(c net.Conn) *K {
	return &K{
		controlByte: 0,
		c:           c,
		sendMsgBuf:  bytes.Buffer{},
	}
}

// Close is function for close and delete conn
func (k *K) Close() (err error) {
	if k.c != nil {
		err = k.c.Close()
		k.c = nil
		if err != nil {
			return fmt.Errorf("transport closing conn error %w", err)
		}
	}

	return nil
}

func (k *K) SendMessage(msg []byte) ([]byte, error) {
	k.sendMsgBuf.Reset()
	defer k.sendMsgBuf.Reset()
	k.sendMsgBuf.Write(msg)
	return k.sendENQ()
}

func (k *K) sendENQ() ([]byte, error) {
	if k.c == nil {
		return nil, net.ErrClosed
	}

	if err := k.writeByte(ENQ); err != nil {
		err = fmt.Errorf("err write ENQ: %w", err)
		return nil, err
	}
	time.Sleep(time.Millisecond * 50)

	if err := k.readControlByte(); err != nil {
		return nil, err
	}

	switch k.controlByte {
	case ACK:
		log.Printf("ENQ ACK")

		msg, err := k.reciveMsg()
		if err != nil {
			log.Err(err).Msg("err in ack")
		}

		log.Debug().Msgf("received message in send enq: %x\n", msg)

		return k.sendENQ()
	case NAK:
		if err := k.sendMsg(); err != nil {
			log.Err(err).Msg("err in nak")
		}

		return k.reciveMsg()
	default:
		log.Debug().Msgf("wrong control byte %X, retry after sleep\n", k.controlByte)
		time.Sleep(time.Millisecond * 50)

		return k.sendENQ()
	}
}

func (k *K) sendMsg() error {
	msgLen := k.sendMsgBuf.Len()
	if msgLen > 255 {
		return fmt.Errorf("sendMsg: %w: %v", ErrOverflowMsgLength, msgLen)
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

	for attempt := 0; attempt < 10; attempt++ {
		if _, err := k.c.Write(buf.Bytes()); err != nil {
			err = fmt.Errorf("err in attempt loop: %w", err)
			log.Err(err).Send()
			continue
		}

		if err := k.readControlByte(); err != nil {
			err = fmt.Errorf("err in attempt loop: %w", err)
			log.Err(err).Send()
			continue
		}

		switch k.controlByte {
		case ACK:
			return nil
		default:
			if attempt < 10 { // 10
				log.Debug().Msgf("attempt %v, ctrlByte: 0x%X \n", attempt, k.controlByte)
				time.Sleep(time.Millisecond * 10)
				continue
			}

			err := fmt.Errorf("send message end: %w: %x", ErrWrongControlByte, k.controlByte)

			return err
		}
	}

	return fmt.Errorf("cannot send msg")
}

// read stx.
// read len.
// mb error while io.EOF
func (k *K) reciveMsg() ([]byte, error) {
	stx, err := k.readByte()
	if err != nil {
		err = fmt.Errorf("error read stx: %w", err)
		return nil, err
	}

	if stx != STX {
		return nil, fmt.Errorf("reciveMsg: %w: %x", ErrWrongSTX, stx)
	}

	lenMsg, err := k.readByte()
	if err != nil {
		err = fmt.Errorf("error read lenMsg: %w", err)
		return nil, err
	}

	msg := make([]byte, lenMsg)
	if _, err = k.c.Read(msg); err != nil {
		return nil, fmt.Errorf("error read message: %w", err)
	}

	lrc, err := k.readByte()
	if err != nil {
		err = fmt.Errorf("error read checksum: %w", err)
		return nil, err
	}

	err = k.writeByte(ACK) // write ack after read last byte in msg
	if err != nil {
		err = fmt.Errorf("error write ACK error: %w", err)
		return nil, err
	}

	var resp []byte
	resp = append(resp, lenMsg)
	resp = append(resp, msg...)

	if !checksum(resp, lrc) {
		return nil, fmt.Errorf("error checksum: %w: %v", ErrInvalidChecksum, lrc)
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
	if err != nil {
		err = fmt.Errorf("error read control byte: %w", err)
	}
	return err
}

// error write
func (k *K) writeByte(b byte) error {
	if _, err := k.c.Write([]byte{b}); err != nil {
		err = fmt.Errorf("error write single byte: %w", err)
		return err
	}
	return nil
}

// error read
func (k *K) readByte() (byte, error) {
	buf := make([]byte, 1)

	_, err := k.c.Read(buf)
	if err != nil {
		err = fmt.Errorf("error read single byte: %w", err)
		return buf[0], err
	}

	return buf[0], nil
}
