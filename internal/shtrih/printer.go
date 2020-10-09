package shtrih

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"shtrih-drv/internal/logger"

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

	data, err := p.sendCommand(FnReadStatus)
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

func (p *Printer) ReadShortStatus() {
	p.logger.Debug("Send command ReadShortStatus")

	data, err := p.sendCommand(ReadShortStatus)
	if err != nil {
		p.logger.Error(err)
		return
	}

	in := bufio.NewReader(bytes.NewReader(data))

	cmdBin, _ := in.ReadBytes(0) // чтение байтов команды и null
	println(hex.Dump(cmdBin))

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
	FMResultCode, _ := in.ReadByte()
	EJResultCode, _ := in.ReadByte()
	receiptOperationsHi, _ := in.ReadByte()

	//int receiptOperations = receiptOperationsLo + (receiptOperationsHi << 8);
	str := fmt.Sprintf("\noperator number: %v\nflags: %v\nmode: %v\nsubMode: %v\n"+
		"receiptOperationsLo: %v\nbatteryState: %v\npowerState: %v\n"+
		"FMResultCode: %v\nEJResultCode: %v\nreceiptOperationsHi: %v", operatorNumber, flags, mode, subMode,
		receiptOperationsLo, batteryState, powerState, FMResultCode, EJResultCode, receiptOperationsHi)

	p.logger.Debug(str)
}

func (p *Printer) PrintReportWithoutClearing() {
	p.logger.Debug("Send command PrintReportWithoutClearing")

	_, err := p.sendCommand(PrintReportWithoutClearing)
	if err != nil {
		p.logger.Error(err)
		return
	}
}

func (p *Printer) writeCashierName(cashierName string) {
	lines := make([]string, 1)
	lines[0] = cashierName
	//directIO(SMFPTR_DIO_WRITE_CASHIER_NAME, null, lines)
}

/**
 * Write cashier name *
 */
//public static final int SMFPTR_DIO_WRITE_CASHIER_NAME = 0x2C;

//DIO
const (
	SMFPTR_DIO_WRITE_CASHIER_NAME = 0x2C
)

// ////////////////////////////////////////////////////////////////////////
// table numbers
// ////////////////////////////////////////////////////////////////////////
// ECR type and mode
const (
	SmfpTableCashier byte = 2
)

func (p *Printer) writeCasierName(name string) error {
	//operator := p.readPrinterStatus().getOperator()
	//p.writeTable(SMFP_TABLE_CASHIER, operator, 2, name)
	return nil
}

//func writeAdminName(String name) throws Exception {
//	ReadShortStatus command = new ReadShortStatus(sysPassword);
//	execute(command);
//	int operator = command.getStatus().getOperatorNumber();
//	writeTable(SMFP_TABLE_CASHIER, operator, 2, name);
//}
//
//String[] lines = (String[]) object;
//DIOUtils.checkObjectMinLength(lines, 1);
//String cashierName = lines[0];
//getPrinter().writeAdminName(cashierName);
//getPrinter().writeCasierName(cashierName);
