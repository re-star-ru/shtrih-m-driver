package commands

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrWongTlvLength = errors.New("len not equal to value(len)")

func writeTlv(w io.Writer, tag, tagLen uint16, value []byte) (err error) {
	buf := make([]byte, 2)

	binary.LittleEndian.PutUint16(buf, tag) // код тега

	_, err = w.Write(buf)
	if err != nil {
		return
	}

	binary.LittleEndian.PutUint16(buf, tagLen) // длинна тега

	_, err = w.Write(buf)
	if err != nil {
		return
	}

	if len(value) != int(tagLen) {
		return ErrWongTlvLength
	}

	_, err = w.Write(value) // значение тега

	return
}
