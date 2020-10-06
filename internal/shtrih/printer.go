package shtrih

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"shtrih-drv/internal/logger"
)

type Printer struct {
	logger logger.Logger
	client *client
}

func NewPrinter(logger logger.Logger) *Printer {
	return &Printer{
		logger: logger,
		client: newClient(logger),
	}
}

func (p *Printer) ReadSerial() {
	p.client.ping()

	if err := p.sendCommand(30, FnReadStatus); err != nil {
		p.logger.Fatal(err)
	}
	data, err := p.client.receiveDataFromFrame()
	if err != nil {
		p.logger.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewReader(data))
	cmdBin, err := r.ReadBytes(0)
	println("cmdbin: ")
	println(hex.Dump(cmdBin))

	////
	fsStatusByte, _ := r.ReadByte()
	println("fs status byte:", fsStatusByte)
	code := uint64(fsStatusByte)

	fmt.Printf("status bits: %064b \n", code)
	println("проведена настройка фн:", code&1 != 0)
	println("isFiscalModeOpened:", code&2 != 0)
	println("isFiscalModeClosed:", code&4 != 0)
	println("Закончена передача фискальных данных в ОФД ", code&8 != 0)
	println()

	fsDocTypeByte, _ := r.ReadByte()
	println("fs doc type:", fsDocTypeByte)

	fsIsDocRecivedByte, _ := r.ReadByte()
	println("fs is doc received:", fsIsDocRecivedByte)

	fsIsDayOpenedByte, _ := r.ReadByte()
	println("fs is day opened:", fsIsDayOpenedByte)

	fsFlagsByte, _ := r.ReadByte()
	println("fs flags:", fsFlagsByte)

	fsDateYearByte, _ := r.ReadByte()
	fsDateMounthByte, _ := r.ReadByte()
	fsDateDayByte, _ := r.ReadByte()
	println("year, mounth, day:", fsDateYearByte, fsDateMounthByte, fsDateDayByte)

	fsTimeHourByte, _ := r.ReadByte()
	fsTimeMinByte, _ := r.ReadByte()
	println("hour, minute:", fsTimeHourByte, fsTimeMinByte)

	fsSerialBytes := make([]byte, 16)
	r.Read(fsSerialBytes)
	println("fs serial: ", string(fsSerialBytes))
	print(hex.Dump(fsSerialBytes))

	fsLastDocNumber := make([]byte, 4)
	r.Read(fsLastDocNumber)

	println("fs last doc number: ", binary.LittleEndian.Uint32(fsLastDocNumber))
	print(hex.Dump(fsLastDocNumber))
}
