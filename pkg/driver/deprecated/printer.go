package deprecated

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/pkg/logger"

	"golang.org/x/text/encoding/charmap"

	"golang.org/x/text/encoding"
)

type Printer struct {
	logger   logger.Logger
	client   *client
	password uint32
	decoder  *encoding.Decoder
	encoder  *encoding.Encoder
}

func (p *Printer) setFiscalReceiptType(fiscalReceiptType int) error {
	return nil
}

func NewPrinter(logger logger.Logger, host string, password uint32) *Printer {
	return &Printer{
		logger:   logger,
		client:   newClient(logger, host),
		password: password,
		decoder:  charmap.Windows1251.NewDecoder(),
		encoder:  charmap.Windows1251.NewEncoder(),
	}
}

func (p *Printer) FnReadStatus() {
	p.logger.Debug("Send command FnReadStatus")

	data, err := p.sendCommand(models.FnReadStatus)
	if err != nil {
		p.logger.Error(err)
		return
	}

	r := bufio.NewReader(bytes.NewReader(data))
	////
	fsStatusByte, _ := r.ReadByte()
	code := uint64(fsStatusByte)
	p.logger.Debug(fmt.Sprintf("status bits: %b", code))

	p.logger.Debug(fmt.Sprintf("\nпроведена настройка фн: %v \nфискальный режим открыт: %v \n"+
		"фискальный режим закрыт: %v \nзакончена передача фискальных данных в ОФД: %v \n",
		code&1 != 0, code&2 != 0, code&4 != 0, code&8 != 0))

	fsDocTypeByte, _ := r.ReadByte()
	fsIsDocRecivedByte, _ := r.ReadByte()
	fsIsDayOpenedByte, _ := r.ReadByte()
	fsFlagsByte, _ := r.ReadByte()

	fsDateYearByte, _ := r.ReadByte()
	fsDateMounthByte, _ := r.ReadByte()
	fsDateDayByte, _ := r.ReadByte()

	fsTimeHourByte, _ := r.ReadByte()
	fsTimeMinByte, _ := r.ReadByte()

	fsSerialBytes := make([]byte, 16)
	r.Read(fsSerialBytes)

	fsLastDocNumber := make([]byte, 4)
	r.Read(fsLastDocNumber)

	p.logger.Debug(fmt.Sprintf(
		"\nfs doc type: %v\n"+
			"fs is doc received: %v\n"+
			"fs is day opened: %v \n"+
			"fs flags: %v \n"+
			"fs date: %v.%v.%v\n"+
			"fs time: %v:%v\n"+
			"fs serial: %s\n"+
			"fs last doc number: %v",
		fsDocTypeByte, fsIsDocRecivedByte, fsIsDayOpenedByte, fsFlagsByte,
		fsDateYearByte, fsDateMounthByte, fsDateDayByte,
		fsTimeHourByte, fsTimeMinByte,
		fsSerialBytes, binary.LittleEndian.Uint32(fsLastDocNumber)),
	)
}

func (p *Printer) ReadShortStatus() byte {
	p.logger.Debug("Send command ReadShortStatus")

	data, cmdLen := p.createCommandData(models.ReadShortStatus)

	buf := bytes.NewBuffer(data)

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

	in := bufio.NewReader(bytes.NewReader(rFrame.DATA))

	operatorNumber, _ := in.ReadByte()

	flags := make([]byte, 2)
	in.Read(flags)

	mode, _ := in.ReadByte() // & 15; ? wft
	subMode, _ := in.ReadByte()
	receiptOperationsLo, _ := in.ReadByte()
	batteryState, _ := in.ReadByte()
	//double batteryVoltage = (double)batteryState / 255.0D * 100.0D * 5.0D / 100.0D;
	//double powerVoltage = (double)powerState * 24.0D / 216.0D * 100.0D / 100.0D;
	powerState, _ := in.ReadByte()
	receiptOperationsHi, _ := in.ReadByte()
	reserved1, _ := in.ReadByte()
	reserved2, _ := in.ReadByte()
	reserved3, _ := in.ReadByte()
	lastResult, _ := in.ReadByte()

	//int receiptOperations = receiptOperationsLo + (receiptOperationsHi << 8);
	str := fmt.Sprintf("\noperator number: %v\nflags: %v\nmode: %v\nsubMode: %v\n"+
		"receiptOperationsLo: %v\nbatteryState: %v\npowerState: %v\n"+
		"receiptOperationsHi: %v\nreserved1: %v\nreserved2: %v\nreserved3: %v\nlastResult: %v", operatorNumber, flags, mode, subMode,
		receiptOperationsLo, batteryState, powerState, receiptOperationsHi, reserved1, reserved2, reserved3, lastResult)

	p.logger.Debug(str)

	return mode
}

func (p *Printer) CheckStatus() bool {
	if p.ReadShortStatus() == 2 {
		return true
	}
	return false
}

func (p *Printer) PrintReportWithoutClearing() {
	p.logger.Debug("Send command PrintReportWithoutClearing")

	_, err := p.sendCommand(models.PrintReportWithoutClearing)
	if err != nil {
		p.logger.Error(err)
		return
	}
}

func (p *Printer) ReadFieldInfo(table, field byte) {
	cmdBinary, cmdLen := p.createCommandData(models.ReadFieldInfo)
	buf := bytes.NewBuffer(cmdBinary)
	buf.WriteByte(table)
	buf.WriteByte(field)
	cmdBinary = buf.Bytes()

	frame := p.client.createFrame(cmdBinary)

	con, _ := net.Dial("tcp", p.client.host)
	defer con.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))

	if err := p.client.sendFrame(frame, con, rw); err != nil {
		p.logger.Fatal(err)
	}

	rFrame, err := p.client.receiveFrame(con, byte(cmdLen), rw)
	if err != nil {
		p.logger.Fatal(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("Field info\n", hex.Dump(rFrame.DATA))
}
