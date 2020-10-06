package shtrih

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"shtrih-drv/internal/logger"
)

type Printer struct {
	logger   logger.Logger
	client   *client
	password uint32
}

func NewPrinter(logger logger.Logger, host string, password uint32) *Printer {
	return &Printer{
		logger:   logger,
		client:   newClient(logger, host),
		password: password,
	}
}

func (p *Printer) Ping() {
	p.client.ping()
}

func (p *Printer) FnReadStatus() {
	p.client.ping()

	con, err := p.sendCommand(FnReadStatus)
	if err != nil {
		p.logger.Fatal(err)
	}
	defer con.Close()

	data, err := p.client.receiveDataFromFrame(con)
	if err != nil {
		p.logger.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewReader(data))
	r.ReadBytes(0) // чтение байтов команды и null

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
	p.client.ping()

	con, err := p.sendCommand(ReadShortStatus)
	if err != nil {
		p.logger.Fatal(err)
	}
	defer con.Close()

	_, err = p.client.receiveDataFromFrame(con)
	if err != nil {
		p.logger.Fatal(err)
	}

}
