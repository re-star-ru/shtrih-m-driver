package printerUsecase

import (
	"bytes"
	"encoding/binary"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/pkg/driver/client"
	"github.com/fess932/shtrih-m-driver/pkg/driver/printer"
)

type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

type printerUsecase struct {
	logger   Logger
	client   client.Usecase
	password uint32
}

func (p *printerUsecase) send(cmd []byte, cmdLen int) (*models.Frame, error) {
	frameToSend := p.createFrame(cmd)
	return p.client.Send(frameToSend, cmdLen)
}

func (p *printerUsecase) createFrame(data []byte) []byte {
	frameBuf := bytes.NewBuffer([]byte{})

	frameBuf.WriteByte(0x02) // write start

	dl := len(data)

	frameBuf.WriteByte(byte(dl)) // write data len
	frameBuf.Write(data)         // write data

	crc := byte(dl)
	for i := 0; i < dl; i++ {
		crc ^= data[i]
	}

	frameBuf.WriteByte(crc) // write control sum

	return frameBuf.Bytes()
}

func (p *printerUsecase) createCommandBuffer(command uint16, password uint32) (data *bytes.Buffer, cmdLen int) {
	dataBuffer := bytes.NewBuffer([]byte{})
	cb := make([]byte, 2)

	binary.BigEndian.PutUint16(cb, command)

	cb = bytes.TrimPrefix(cb, []byte{0})

	dataBuffer.Write(cb) // write command

	passwordBinary := make([]byte, 4)

	binary.LittleEndian.PutUint32(passwordBinary, password)

	dataBuffer.Write(passwordBinary) // write password

	return dataBuffer, len(cb)
}

func (p *printerUsecase) intToBytesWithLen(val int64, bytesLen int64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}

	return buf.Bytes()[:bytesLen], nil
}

func NewPrinterUsecase(logger Logger, usecase client.Usecase, password uint32) printer.Usecase {
	return &printerUsecase{logger: logger, client: usecase, password: password}
}
