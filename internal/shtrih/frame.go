package shtrih

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
)

const errLen = 1

type frame struct {
	STX  byte
	DLEN byte
	CMD  []byte
	ERR  byte
	DATA []byte
	CRC  byte
}

func (f *frame) bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(f.STX)
	buf.WriteByte(f.DLEN)
	buf.Write(f.CMD)
	buf.WriteByte(f.ERR)
	buf.Write(f.DATA)
	buf.WriteByte(f.CRC)

	return buf.Bytes()
}

func (f *frame) checkCRC() error {

	buf := bytes.NewBuffer([]byte{})
	buf.Write(f.CMD)
	buf.WriteByte(f.ERR)
	buf.Write(f.DATA)
	dataForCheck := buf.Bytes()

	dlen := len(dataForCheck)

	crc := byte(dlen)
	for i := 0; i < dlen; i++ {
		crc ^= dataForCheck[i]
	}

	if f.CRC != crc {
		return errors.New("crc does not match")
	}

	return nil
}

func (c *client) receiveFrame(con net.Conn, cmdLen byte) (*frame, error) {
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	defer func() {
		rw.WriteByte(byte(NAK))
		rw.Flush()
		con.Close()
	}()

	var FRM frame
	var err error

	FRM.STX, err = rw.ReadByte() // read byte STX (0x02) err need
	if err != nil {
		log.Fatal(err)
	}

	FRM.DLEN, err = rw.ReadByte() // read byte dataLen

	FRM.CMD = make([]byte, cmdLen)
	rw.Read(FRM.CMD) // read cmd bytes

	FRM.ERR, _ = rw.ReadByte() // read err byte

	FRM.DATA = make([]byte, FRM.DLEN-cmdLen-errLen)
	rw.Read(FRM.DATA) // read data bytes

	FRM.CRC, _ = rw.ReadByte() // read crc byte

	c.logger.Debug("<- recive frame: \n",
		fmt.Sprintf("stx: %v, dlen: %v, crc: %v  \n", FRM.STX, FRM.DLEN, FRM.CRC),
		hex.Dump(FRM.bytes()))

	dataCheck := bytes.NewBuffer([]byte{})
	dataCheck.Write(FRM.CMD)
	dataCheck.WriteByte(FRM.ERR)
	dataCheck.Write(FRM.DATA)
	if err := FRM.checkCRC(); err != nil {
		return nil, err
	}

	return &FRM, nil
}
