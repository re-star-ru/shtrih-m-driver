package shtrih

import (
	"bytes"
	"errors"
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
