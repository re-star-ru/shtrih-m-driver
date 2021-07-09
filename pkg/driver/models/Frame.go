package models

import (
	"bytes"
	"errors"
)

type Frame struct {
	DATA []byte
	CMD  []byte
	STX  byte
	DLEN byte
	ERR  byte
	CRC  byte
}

func (f *Frame) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(f.STX)
	buf.WriteByte(f.DLEN)
	buf.Write(f.CMD)
	buf.WriteByte(f.ERR)
	buf.Write(f.DATA)
	buf.WriteByte(f.CRC)

	return buf.Bytes()
}

func (f *Frame) CheckCRC() error {
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
