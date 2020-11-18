package deprecated

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/pkg/logger"
)

type client struct {
	logger logger.Logger
	host   string
}

func newClient(logger logger.Logger, host string) *client {
	return &client{logger: logger, host: host}
}

func (c *client) ping(rw *bufio.ReadWriter, con net.Conn) {
	//rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	c.logger.Debug("-> send ENQ")
	rw.WriteByte(models.ENQ)
	rw.Flush()

	b, _ := rw.ReadByte()
	c.logger.Debug("<- recive control byte:", b)

	switch b {
	case models.ACK:
		c.logger.Debug("OK, ACK, wait for recive now")
		rw.WriteByte(models.ACK)
		rw.Flush()
	case models.NAK:
		c.logger.Debug("OK, NAK, wait for cmd now")
		rw.WriteByte(models.ACK)
		rw.Flush()
	default:
		rw.WriteByte(models.ACK)
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

func (c *client) sendFrame(frame []byte, con net.Conn, rw *bufio.ReadWriter) error {
	//rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	c.ping(rw, con)

	rw.Write(frame)
	rw.Flush()
	c.logger.Debug("-> send frame: \n", hex.Dump(frame))

	b, err := rw.ReadByte()
	c.logger.Debug("<- recive control byte:", b)

	if err != nil {
		return err
	}
	switch b {
	case models.ACK:
		return nil
	case models.NAK:
		return errors.New("ошибка интерфейса либо неверная контрольная сумма")
	default:
		return errors.New("сообщение не принято либо не верные данные")
	}
}

func (c *client) receiveFrame(con net.Conn, cmdLen byte, rw *bufio.ReadWriter) (*models.Frame, error) {
	c.logger.Debug("<- Receive frame")

	//rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	defer func() {
		rw.WriteByte(models.ACK)
		rw.Flush()
	}()

	var FRM models.Frame
	var err error

	FRM.STX, err = rw.ReadByte() // read byte STX (0x02) err need
	if err != nil {
		c.logger.Fatal(err)
	}

	FRM.DLEN, err = rw.ReadByte() // read byte dataLen

	FRM.CMD = make([]byte, cmdLen)
	rw.Read(FRM.CMD) // read cmd bytes

	FRM.ERR, _ = rw.ReadByte() // read err byte

	FRM.DATA = make([]byte, FRM.DLEN-cmdLen-models.ErrLen)
	rw.Read(FRM.DATA) // read data bytes

	FRM.CRC, _ = rw.ReadByte() // read crc byte

	c.logger.Debug("<- recive frame: \n",
		fmt.Sprintf("stx: %v, dlen: %v, crc: %v  \n", FRM.STX, FRM.DLEN, FRM.CRC),
		hex.Dump(FRM.Bytes()))

	dataCheck := bytes.NewBuffer([]byte{})
	dataCheck.Write(FRM.CMD)
	dataCheck.WriteByte(FRM.ERR)
	dataCheck.Write(FRM.DATA)
	if err := FRM.CheckCRC(); err != nil {
		return nil, err
	}

	return &FRM, nil
}
