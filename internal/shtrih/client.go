package shtrih

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"shtrih-drv/internal/logger"
)

type client struct {
	logger logger.Logger
	host   string
}

func newClient(logger logger.Logger, host string) *client {
	return &client{logger: logger, host: host}
}

func (c *client) ping() {
	con, err := net.Dial("tcp", c.host)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer con.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	rw.WriteByte(ENQ)
	rw.Flush()
	c.logger.Debug("-> send ENQ")

	b, _ := rw.ReadByte()
	c.logger.Debug("<- recive control byte")
	switch b {
	case ACK:
		c.logger.Debug("OK, ACK, recive now")
		rw.WriteByte(ACK)
		rw.Flush()
	case NAK:
		c.logger.Debug("OK, NAK, nothing to recive")
	default:
		rw.WriteByte(ACK)
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

func (c *client) sendFrame(frame []byte, con net.Conn) error {
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	rw.Write(frame)
	rw.Flush()
	c.logger.Debug("-> send frame: \n", hex.Dump(frame))

	b, err := rw.ReadByte()
	if err != nil {
		con.Close()
		c.logger.Fatal(err)
	}
	c.logger.Debug("<- recive control byte: ", b)

	switch b {
	case ACK:
		return nil
	case NAK:
		return errors.New("NAK, nothig to recive")
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

func (c *client) receiveDataFromFrame(con net.Conn) ([]byte, error) {
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	defer func() {
		rw.WriteByte(byte(NAK))
		rw.Flush()
		con.Close()
	}()

	var frame bytes.Buffer

	stx, err := rw.ReadByte() // read byte STX (0x02) err need
	if err != nil {
		log.Fatal(err)
	}
	dlen, _ := rw.ReadByte() // read byte dataLen
	data := make([]byte, dlen)
	n, _ := rw.Read(data)   // read data bytes
	crc, _ := rw.ReadByte() // read crc byte

	frame.WriteByte(stx)
	frame.WriteByte(dlen)
	frame.Write(data)
	frame.WriteByte(crc)
	c.logger.Debug("<- recive frame: \n",
		fmt.Sprintf("stx: %v, data len: %v, crc: %v \n", stx, dlen, n),
		hex.Dump(frame.Bytes()))

	if err := checkCRC(data, crc); err != nil {
		return nil, err
	}
	return data, nil
}
