package command

import (
	"bytes"
)

type PrinterCommand struct {
	DefaultTimeout     int
	timeout            int
	resultCode         int
	charsetName        string
	repeatEnabled      bool
	repeadNeeded       bool
	errorReportEnabled bool

	commandCode int

	text string

	txData   []byte
	rxData   []byte
	password int
}

func (p *PrinterCommand) EncodeData() ([]byte, error) {

	code := p.commandCode

	buf := bytes.NewBuffer([]byte{})

	if code > 255 {
		buf.WriteByte(byte(code >> 8 & 255))
		buf.WriteByte(byte(code & 255))
	}
	if code <= 255 {
		buf.WriteByte(byte(code))
	}

	buf2 := bytes.NewBuffer([]byte{})
	if err := p.Encode(buf2); err != nil {
		return nil, err
	}
	p.txData = buf2.Bytes()
	buf.Write(p.txData)

	return buf.Bytes(), nil
}

func (p *PrinterCommand) GetCode() int {
	return p.commandCode
}

func (p *PrinterCommand) GetText() string {
	return p.text
}

func (c *PrinterCommand) Encode(buf *bytes.Buffer) error {
	s := 30

	if err := buf.WriteByte(byte(s >> 0 & 255)); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(s >> 8 & 255)); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(s >> 16 & 255)); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(s >> 24 & 255)); err != nil {
		return err
	}

	return nil
}

func NewPrinterCommand() *PrinterCommand {
	return &PrinterCommand{
		10000,
		0,
		0,
		"Cp1251", // подумать везде сделать utf-8
		false,
		false,
		true,
		0,
		"default",
		[]byte{},
		[]byte{},
		30,
	}
}
