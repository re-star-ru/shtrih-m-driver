package fiscalprinter

import (
	"bytes"
	"shtrih-drv/internal/fiscalprinter/port"
)

func NewFrame() *Frame {
	return &Frame{
		2,
	}
}

type Frame struct {
	STX byte
}

func (f *Frame) GetCrc(data []byte) byte {
	crc := byte(len(data))

	for i := 0; i < len(data); i++ {
		crc ^= data[i]
	}

	return crc
}

func (f *Frame) encode(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	if len(data) > 255 {
		return nil, port.DataLenghtExeeds
	} else {
		buf.WriteByte(2)
		buf.WriteByte(byte(len(data)))
		buf.Write(data)
		buf.WriteByte(f.GetCrc(data))

		return buf.Bytes(), nil
	}
}
