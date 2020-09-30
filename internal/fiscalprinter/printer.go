package fiscalprinter

func NewFiscalPrinter() *FiscalPrinter {
	return &FiscalPrinter{}
}

type FiscalPrinter struct {
	physicalDeviceName string
}

func (p FiscalPrinter) getPhysicalDeviceName() string {
	return p.physicalDeviceName
}

func (p FiscalPrinter) GetSerialNumber() string {
	// команда чтения серийного номера из бд устройства
	SMFPTR_DIO_READ_SERIAL := 22

	serial := [2]string{"", ""}

	// direct handler 2
	directIO(SMFPTR_DIO_READ_SERIAL, nil, serial)

	return serial[0]
}

func directIO(command int, data []int, object interface{}) {
	switch command {
	case 22:
		execute(data, object)
	}
}

func execute(data []int, object interface{}) {
	//DIOUtils.checkObjectMinLength((String[])((String[])object), 1);
	//String[] serial = (String[])((String[])object);
	//serial[0] = "";
	//LongPrinterStatus status = this.service.readLongStatus();
	//if (status.getRegistrationNumber() > 0) {
	//	ReadEJSerialNumber command = new ReadEJSerialNumber();
	//	command.setPassword(this.service.fiscalprinter.getSysPassword());
	//	this.service.fiscalprinter.execute(command);
	//	serial[0] = String.valueOf(command.getSerial());
	//}
}

//
//func decode(in CommandInputStream)  {
//	int operatorNumber = in.readByte();
//	int flags = in.readShort();
//	int mode = in.readByte() & 15;
//	int subMode = in.readByte();
//	int receiptOperationsLo = in.readByte();
//	int batteryState = in.readByte();
//	int powerState = in.readByte();
//	int FMResultCode = in.readByte();
//	int EJResultCode = in.readByte();
//	int receiptOperationsHi = in.readByte();
//	double batteryVoltage = (double)batteryState / 255.0D * 100.0D * 5.0D / 100.0D;
//	double powerVoltage = (double)powerState * 24.0D / 216.0D * 100.0D / 100.0D;
//	int receiptOperations = receiptOperationsLo + (receiptOperationsHi << 8);
//	this.status = new ShortPrinterStatus(mode, flags, subMode, FMResultCode, EJResultCode, receiptOperations, batteryVoltage, powerVoltage, operatorNumber);
//}
//
//type CommandInputStream struct {
//	stream io.ByteReader
//}
//
//func (in *CommandInputStream) readByte()  {
//	// package com.shtrih.fiscalprinter.command;
//	stream
//}
