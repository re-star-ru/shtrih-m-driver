package command

import "bytes"

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

	txData []byte
	rxData []byte
}

func (p *PrinterCommand) EncodeData() []byte {

	code := p.commandCode

	buf := bytes.NewBuffer([]byte{})

	if code > 255 {
		buf.WriteByte(byte(code >> 8 & 255))
		buf.WriteByte(byte(code & 255))
	}
	if code <= 255 {
		buf.WriteByte(byte(code))
	}

	return buf.Bytes()
}

func (p *PrinterCommand) GetCode() int {
	return p.commandCode
}

func (p *PrinterCommand) GetText() string {
	return p.text
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
	}
}
