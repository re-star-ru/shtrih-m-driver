package shtrih

import (
	"bytes"
	"encoding/binary"
	"net"
)

func (p *Printer) createCommand(command uint16) []byte {
	dataBuffer := bytes.NewBuffer([]byte{})

	cb := make([]byte, 2)
	binary.BigEndian.PutUint16(cb, command)
	cb = bytes.TrimPrefix(cb, []byte{0})

	dataBuffer.Write(cb)

	passwordBinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(passwordBinary, p.password)
	dataBuffer.Write(passwordBinary) // write password

	return dataBuffer.Bytes()
}

func (p *Printer) sendCommand(command uint16) (net.Conn, error) {
	cmdBinary := p.createCommand(command)
	frame := p.client.createFrame(cmdBinary)

	con, _ := net.Dial("tcp", p.client.host)
	if err := p.client.sendFrame(frame, con); err != nil {
		return nil, err
	}

	return con, nil
}
