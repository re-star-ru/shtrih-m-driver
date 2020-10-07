package shtrih

import (
	"bytes"
	"encoding/binary"
	"net"
)

func (p *Printer) createCommandData(command uint16) ([]byte, int) {
	dataBuffer := bytes.NewBuffer([]byte{})

	cb := make([]byte, 2)
	binary.BigEndian.PutUint16(cb, command)
	cb = bytes.TrimPrefix(cb, []byte{0})

	dataBuffer.Write(cb)

	passwordBinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(passwordBinary, p.password)
	dataBuffer.Write(passwordBinary) // write password

	return dataBuffer.Bytes(), len(cb)
}

func (p *Printer) sendCommand(command uint16) ([]byte, error) {
	cmdBinary, cmdLen := p.createCommandData(command)
	frame := p.client.createFrame(cmdBinary)

	con, _ := net.Dial("tcp", p.client.host)
	if err := p.client.sendFrame(frame, con); err != nil {
		return nil, err
	}

	rFrame, err := p.client.receiveFrame(con, byte(cmdLen))
	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		return nil, err
	}

	return rFrame.DATA, nil
}
