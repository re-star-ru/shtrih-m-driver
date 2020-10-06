package shtrih

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"shtrih-drv/internal/logger"
)

type client struct {
	logger logger.Logger
}

func newClient(logger logger.Logger) *client {
	return &client{logger: logger}
}

func (c *client) ping() {
	con, err := net.Dial("tcp", "10.51.0.71:7778")
	defer con.Close()
	if err != nil {
		c.logger.Fatal(err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	rw.WriteByte(0x05)
	rw.Flush()
	b, _ := rw.ReadByte()
	switch b {
	case 0x6:
		c.logger.Debug("OK, need recive, recive now")
		rw.WriteByte(0x06)
		rw.Flush()
	case 0x15:
		c.logger.Debug("OK, nothing to recive")
	default:
		rw.WriteByte(0x06)
		rw.Flush()
		c.logger.Fatal("ERR, ping byte:", b)
	}
}

func (c *client) createFrame(data []byte) []byte {
	// frame buffer
	frameBuf := bytes.NewBuffer([]byte{})
	frameBuf.WriteByte(0x02) // write start

	dl := len(data)
	frameBuf.WriteByte(byte(dl)) // write data len
	frameBuf.Write(data)         // write data

	crc := byte(dl)
	for i := 0; i < dl; i++ {
		crc ^= data[i]
	}
	frameBuf.WriteByte(crc) // write control sum

	return frameBuf.Bytes()
}

func (c *client) sendFrame(frame []byte) error {
	con, err := net.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		c.logger.Fatal(err)
	}
	defer con.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	rw.Write(frame)
	rw.Flush()

	b, err := rw.ReadByte()
	if err != nil {
		c.logger.Fatal(err)
	}

	switch b {
	case 0x06:
		return nil
	case 0x15:
		return errors.New("21, nothig to recive")
	default:
		return errors.New(fmt.Sprint("control byte is:", b, ": wft?"))
	}
}

func checkCRC(data []byte, rcrc byte) error {
	dlen := len(data)

	crc := byte(dlen)
	for i := 0; i < dlen; i++ {
		crc ^= data[i]
	}

	if rcrc != crc {
		return errors.New("crc does not match")
	}

	return nil
}

func (c *client) receiveDataFromFrame() ([]byte, error) {
	con, err := net.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		c.logger.Fatal(err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	defer func() {
		rw.WriteByte(byte(0x06))
		rw.Flush()

		con.Close()
	}()

	var frame bytes.Buffer

	stx, _ := rw.ReadByte() // read byte STX (0x02) err need
	println("stx byte:", stx)

	dlen, _ := rw.ReadByte() // read byte dataLen (0x02)
	println("data len:", dlen)

	data := make([]byte, dlen)
	n, _ := rw.Read(data)
	println("read data bytes len:", n)

	crc, _ := rw.ReadByte() // read crc byte
	println("crc:", crc)

	frame.WriteByte(stx)
	frame.WriteByte(dlen)
	frame.Write(data)
	frame.WriteByte(crc)
	c.logger.Debug("<- recive frame")
	print(hex.Dump(frame.Bytes()))

	if err := checkCRC(data, crc); err != nil {
		return nil, err
	}
	return data, nil
}
