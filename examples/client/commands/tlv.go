package commands

import (
	"encoding/binary"
	"errors"
	"io"
)

func writeTlv(w io.Writer, tag, Len uint16, value []byte) error {
	buf := make([]byte, 2)

	binary.LittleEndian.PutUint16(buf, tag) // код тега
	w.Write(buf)

	binary.LittleEndian.PutUint16(buf, Len) // длинна тега 	// значение тега
	w.Write(buf)

	if len(value) != int(Len) {
		return errors.New("длинна не совпадает со значением")
	}

	w.Write(value)

	return nil
}
