package shtrih

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

func (p *Printer) createCommandData(command uint16) (data []byte, cmdLen int) {
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
	defer con.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))

	if err := p.client.sendFrame(frame, con, rw); err != nil {
		return nil, err
	}

	rFrame, err := p.client.receiveFrame(con, byte(cmdLen), rw)
	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		return nil, err
	}

	var cmdCodeRecived uint16
	if len(rFrame.CMD) == 1 {
		cmdCodeRecived = uint16(rFrame.CMD[0])
	}
	if len(rFrame.CMD) == 2 {
		cmdCodeRecived = binary.BigEndian.Uint16(rFrame.CMD)
	}
	if cmdCodeRecived != command {
		return nil, errors.New("отправленная и полученная команды не совпадаютс")
	}

	return rFrame.DATA, nil
	//return nil, nil
}

func (p *Printer) WriteTable(tableNumber byte, rowNumber uint16, fieldNumber byte, fieldValue string) {
	data, cmdLen := p.createCommandData(WriteTable)

	buf := bytes.NewBuffer(data)

	buf.WriteByte(tableNumber)

	cb := make([]byte, 2)
	binary.LittleEndian.PutUint16(cb, rowNumber)
	buf.Write(cb)

	buf.WriteByte(fieldNumber)

	fvb, _ := p.encoder.Bytes([]byte(fieldValue)) // конвентируем строку в win1251
	buf.Write(fvb)
	frame := p.client.createFrame(buf.Bytes())

	con, _ := net.Dial("tcp", p.client.host)
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))

	defer con.Close()
	if err := p.client.sendFrame(frame, con, rw); err != nil {
		p.logger.Fatal(err)
	}

	rFrame, err := p.client.receiveFrame(con, byte(cmdLen), rw)
	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}
}
