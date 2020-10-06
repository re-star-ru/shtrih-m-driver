package main

//
//import (
//	"bufio"
//	"bytes"
//	"encoding/binary"
//	"encoding/hex"
//	"errors"
//	"fmt"
//	"log"
//	"net"
//)
//
////buf.WriteByte(2) // start
////buf.WriteByte(byte(len(data))) // len data
////buf.Write(data) // data
////buf.WriteByte(f.GetCrc(data)) // control sum (crc)
//
//// win drv
////02 05 10 1E 00 00 00 0B
//// go
////02 05 10 1e 00 00 00 0b
//
////
////02 05 10 1E 00 00 00 0B
////02 05 0a 1e 00 00 00 11
////02 06 FF 02 02 00 00 00 F9
////02 06 FF 01 02 00 00 00 FA
////02 05 10    02 00 00 00 17
////start data len  data				   crc
////(02) (  06   )  (FF 01 1E 00 00 00)  E6
//
//// data
//// command  password
//// (FF 01) (1E 00 00 00)
//// 1 байт -начало команды
//// 2 байт, 3 байт - команда либо 2,3,4 байты команда фискального принтера
//// 4-7 байты - пароль либо 5-8 байты для фискального принтера
//// 8 байт - контрольная сумма либо 9 байт для фискального принтера
//
////02 начало команды 1 байт - начало команды
////FF 01  - признак команда Фискального Накопителя (ФН) или 05 - признак команда Принтера (П)
////(00 - FF) команда
////(00 - FF) пароль
//func main() {
//	//println(hex.Dump(createCommand(30, ReadShortStatus)))
//	//println(hex.Dump(createCommand(30, FnReadStatus)))
//	//println(hex.Dump(createCommand(30, FnReadSerial)))
//
//	readSerial()
//}
//
//func createFrame(data []byte) []byte {
//	// frame buffer
//	frameBuf := bytes.NewBuffer([]byte{})
//	frameBuf.WriteByte(0x02) // write start
//
//	dl := len(data)
//	frameBuf.WriteByte(byte(dl)) // write data len
//	frameBuf.Write(data)         // write data
//
//	crc := byte(dl)
//	for i := 0; i < dl; i++ {
//		crc ^= data[i]
//	}
//	frameBuf.WriteByte(crc) // write control sum
//
//	return frameBuf.Bytes()
//}
//
//const (
//	// команды принтера
//	ReadShortStatus = 16
//
//	// команды фискального накопителя
//	FnReadStatus = 65281
//	FnReadSerial = 65282
//)
//
//func createCommand(password uint32, command uint16) []byte {
//	dataBuffer := bytes.NewBuffer([]byte{})
//
//	cb := make([]byte, 2)
//	binary.BigEndian.PutUint16(cb, command)
//	cb = bytes.TrimPrefix(cb, []byte{0})
//
//	dataBuffer.Write(cb)
//
//	passwordBinary := make([]byte, 4)
//	binary.LittleEndian.PutUint32(passwordBinary, password)
//	dataBuffer.Write(passwordBinary) // write password
//
//	return dataBuffer.Bytes()
//}
//
//func ping() {
//	con, err := net.Dial("tcp", "10.51.0.71:7778")
//	defer con.Close()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
//	rw.WriteByte(0x05)
//	rw.Flush()
//	b, _ := rw.ReadByte()
//	switch b {
//	case 0x6:
//		log.Println("OK, need recive, recive now")
//		rw.WriteByte(0x06)
//		rw.Flush()
//	case 0x15:
//		log.Println("OK, nothing to recive")
//	default:
//		rw.WriteByte(0x06)
//		rw.Flush()
//		log.Fatal("ERR, ping byte:", b)
//	}
//}
//
//func readSerial() {
//	ping()
//
//	if err := sendCommand(30, FnReadStatus); err != nil {
//		log.Fatal(err)
//	}
//	data, err := reciveDataFromFrame()
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
//
//func sendCommand(password uint32, command uint16) error {
//	cmdBinary := createCommand(password, command)
//	frame := createFrame(cmdBinary)
//	println(hex.Dump(frame))
//
//	if err := sendFrame(frame); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func sendFrame(frame []byte) error {
//	con, err := net.Dial("tcp", "10.51.0.71:7778")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer con.Close()
//
//	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
//	rw.Write(frame)
//	rw.Flush()
//
//	b, err := rw.ReadByte()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	switch b {
//	case 0x06:
//		return nil
//	case 0x15:
//		return errors.New("21, nothig to recive")
//	default:
//		return errors.New(fmt.Sprint("control byte is:", b, ": wft?"))
//	}
//}
//
//func checkCRC(data []byte, rcrc byte) error {
//	dlen := len(data)
//
//	crc := byte(dlen)
//	for i := 0; i < dlen; i++ {
//		crc ^= data[i]
//	}
//
//	if rcrc != crc {
//		return errors.New("crc does not match")
//	}
//
//	return nil
//}
//
//func reciveDataFromFrame() ([]byte, error) {
//	con, err := net.Dial("tcp", "10.51.0.71:7778")
//	if err != nil {
//		log.Fatal(err)
//	}
//	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
//	defer func() {
//		rw.WriteByte(byte(0x06))
//		rw.Flush()
//
//		con.Close()
//	}()
//
//	var frame bytes.Buffer
//
//	stx, _ := rw.ReadByte() // read byte STX (0x02) err need
//	println("stx byte:", stx)
//
//	dlen, _ := rw.ReadByte() // read byte dataLen (0x02)
//	println("data len:", dlen)
//
//	data := make([]byte, dlen)
//	n, _ := rw.Read(data)
//	println("read data bytes len:", n)
//
//	crc, _ := rw.ReadByte() // read crc byte
//	println("crc:", crc)
//
//	frame.WriteByte(stx)
//	frame.WriteByte(dlen)
//	frame.Write(data)
//	frame.WriteByte(crc)
//	log.Println("<- recive frame")
//	print(hex.Dump(frame.Bytes()))
//
//	if err := checkCRC(data, crc); err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//// fn read status
//// 02 06 FF 01 1E 00 00 00 E6
//// 02 | 21 | FF | 01 | 00 | 03 00 00 01 00 14 0A 06 0B 1A 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32 3C 2E 00 00 | C0
//
//// fn read status go
//// 02 06 ff 01 1e 00 00 00 e6
//// stx| len|   cmd   | nil|   data 																					  |	crc
//// 02 | 21 | ff | 01 | 00 | 03 00 00 01 00 14 0a 06 0b 1a 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32 3c 2e 00 00 | c0
//
////----------------------------------------
////Состояние фазы жизни    : 03h, 3
////----------------------------------------
////Проведена настройка ФН: [да]
////Открыт фискальный режим: [да]
////Закрыт фискальный режим: [нет]
////Закончена передача фискальных данных в ОФД: [нет]
////----------------------------------------
////Текущий документ        : 00h, Нет открытого документа
////Данные документа        : 00h, Нет данных документа
////Состояние смены         : 01h, Смена открыта
////Флаги предупреждения    : 00h, 0
////Дата и время ФН         : 06.10.2020 11:26:00
////Номер ФН                : 9281000100007442
////							9281000100007442
////Номер последнего ФД     : 11836
//
////02 | 10 | 10 | 00 | 1E 92 02 02 00 00 9F E5 18 01 00 89 08 00 | 6E
////02 | 10 | 10 | 00 | 1e 92 02 02 00 00 9f e5 18 01 00 89 08 00 | 6e
////02 | 10 | 10 | 00 | 1E 92 02 02 00 00 9F E5 18 01 00 89 08 00 | 6E
////					   (delim)
////stx  19-(1 comnd 2) 3 - null   4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 | crc
////02 | 13 | FF 02    |   00   | 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32 | E8
////02 | 13 | ff 02    |   00   | 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32 | e8
////----------------------------------------
////Краткий запрос состояния:
////----------------------------------------
////Режим:
////2, Открытая смена; 24 часа не кончились
////----------------------------------------
////Подрежим                  : 0, Бумага присутствует
////Статус режима             : 0
////Количество операций в чеке: 0
////Напряжение батареи, В     : 3,12
////Напряжение источника, В   : 25,44
////----------------------------------------
////ФлагиKKT                  : 0292h, 658
////----------------------------------------
////Увеличенная точность количества  : [нет]
////Бумага на выходе из накопителя   : [нет]
////Бумага на входе в накопитель     : [нет]
////Денежный ящик открыт             : [нет]
////Крышка корпуса поднята           : [нет]
////Рычаг термоголовки чека опущен   : [да]
////Рычаг термоголовки журнала опущен: [да]
////Оптический датчик чека           : [да]
////Оптический датчик журнала        : [да]
////2 знака после запятой в цене     : [да]
////Нижний датчик ПД                 : [да]
////Верхний датчик ПД                : [да]
////Рулон чековой ленты              : [да]
////Рулон контрольной ленты          : [да]
////----------------------------------------
////02 | 10 | 10 | 00 | 1E 92 02 02 00 00 9F E5 18 01 00 89 08 00 | 6E
//
////int operatorNumber = in.readByte();
////int flags = in.readShort();
////int mode = in.readByte() & 15;
////int subMode = in.readByte();
////int receiptOperationsLo = in.readByte();
////int batteryState = in.readByte();
////int powerState = in.readByte();
////int FMResultCode = in.readByte();
////int EJResultCode = in.readByte();
////int receiptOperationsHi = in.readByte();
////double batteryVoltage = (double)batteryState / 255.0D * 100.0D * 5.0D / 100.0D;
////double powerVoltage = (double)powerState * 24.0D / 216.0D * 100.0D / 100.0D;
////int receiptOperations = receiptOperationsLo + (receiptOperationsHi << 8);
