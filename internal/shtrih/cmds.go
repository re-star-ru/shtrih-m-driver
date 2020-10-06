package shtrih

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

func (p *Printer) createCommand(password uint32, command uint16) []byte {
	dataBuffer := bytes.NewBuffer([]byte{})

	cb := make([]byte, 2)
	binary.BigEndian.PutUint16(cb, command)
	cb = bytes.TrimPrefix(cb, []byte{0})

	dataBuffer.Write(cb)

	passwordBinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(passwordBinary, password)
	dataBuffer.Write(passwordBinary) // write password

	return dataBuffer.Bytes()
}

func (p *Printer) sendCommand(password uint32, command uint16) error {
	cmdBinary := p.createCommand(password, command)
	frame := p.client.createFrame(cmdBinary)
	println(hex.Dump(frame))

	if err := p.client.sendFrame(frame); err != nil {
		return err
	}

	return nil
}

//func ReadSerial() {
//	ping()
//
//	if err := sendCommand(30, FnReadStatus); err != nil {
//		log.Fatal(err)
//	}
//	data, err := receiveDataFromFrame()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	r := bufio.NewReader(bytes.NewReader(data))
//	cmdBin, err := r.ReadBytes(0)
//	println("cmdbin: ")
//	println(hex.Dump(cmdBin))
//
//	////
//	fsStatusByte, _ := r.ReadByte()
//	println("fs status byte:", fsStatusByte)
//	code := uint64(fsStatusByte)
//
//	fmt.Printf("status bits: %064b \n", code)
//	println("проведена настройка фн:", code&1 != 0)
//	println("isFiscalModeOpened:", code&2 != 0)
//	println("isFiscalModeClosed:", code&4 != 0)
//	println("Закончена передача фискальных данных в ОФД ", code&8 != 0)
//	println()
//
//	fsDocTypeByte, _ := r.ReadByte()
//	println("fs doc type:", fsDocTypeByte)
//
//	fsIsDocRecivedByte, _ := r.ReadByte()
//	println("fs is doc received:", fsIsDocRecivedByte)
//
//	fsIsDayOpenedByte, _ := r.ReadByte()
//	println("fs is day opened:", fsIsDayOpenedByte)
//
//	fsFlagsByte, _ := r.ReadByte()
//	println("fs flags:", fsFlagsByte)
//
//	fsDateYearByte, _ := r.ReadByte()
//	fsDateMounthByte, _ := r.ReadByte()
//	fsDateDayByte, _ := r.ReadByte()
//	println("year, mounth, day:", fsDateYearByte, fsDateMounthByte, fsDateDayByte)
//
//	fsTimeHourByte, _ := r.ReadByte()
//	fsTimeMinByte, _ := r.ReadByte()
//	println("hour, minute:", fsTimeHourByte, fsTimeMinByte)
//
//	fsSerialBytes := make([]byte, 16)
//	r.Read(fsSerialBytes)
//	println("fs serial: ", string(fsSerialBytes))
//	print(hex.Dump(fsSerialBytes))
//
//	fsLastDocNumber := make([]byte, 4)
//	r.Read(fsLastDocNumber)
//
//	println("fs last doc number: ", binary.LittleEndian.Uint32(fsLastDocNumber))
//	print(hex.Dump(fsLastDocNumber))
//}
