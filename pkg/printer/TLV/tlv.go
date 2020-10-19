package TLV

import (
	"bytes"
	"encoding/binary"
)

//type Writer struct {
//	buf bufio.Writer
//}

// TEXT в TLV кодировка : CP866

// tlv types:
// string cp866
// uint16

const (
	StringTLV = iota
	Uint16TLV
	Uint32TLV
)

type TLVData []byte

type TLVStruct struct {
	Tag  string
	Len  string
	Data struct {
		Type  int
		Value []byte
	}
}

// тип   длинна  значение (длинна сейчас 16)1
//11 04 | 10 00 | 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32

func New(tag uint16, dataLen uint16) TLVData {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, tag)
	binary.Write(buf, binary.LittleEndian, dataLen)

	return buf.Bytes()
}

func WriteUint16() {

}

func WriteString() {

}
