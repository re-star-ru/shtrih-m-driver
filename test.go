package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"
	"time"
)

//buf.WriteByte(2) // start
//buf.WriteByte(byte(len(data))) // len data
//buf.Write(data) // data
//buf.WriteByte(f.GetCrc(data)) // control sum (crc)

// win drv
//02 05 10 1E 00 00 00 0B
// go
//02 05 10 1e 00 00 00 0b

//
//02 05 10 1E 00 00 00 0B
//02 05 0a 1e 00 00 00 11
//02 06 FF 02 02 00 00 00 F9
//02 06 FF 01 02 00 00 00 FA
//02 05 10    02 00 00 00 17
//start data len  data				   crc
//(02) (  06   )  (FF 01 1E 00 00 00)  E6

// data
// command  password
// (FF 01) (1E 00 00 00)
// 1 байт -начало команды
// 2 байт, 3 байт - команда либо 2,3,4 байты команда фискального принтера
// 4-7 байты - пароль либо 5-8 байты для фискального принтера
// 8 байт - контрольная сумма либо 9 байт для фискального принтера

//02 начало команды 1 байт - начало команды
//06 FF 01  - признак команда Фискального Накопителя (ФН) или 05 - признак команда Принтера (П)
//(00 - FF) команда
//(00 - FF) пароль
func main() {
	status()
}

func createFrame(data []byte) []byte {
	// frame buffer
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

func readShortStatusCommand(password uint32) []byte {
	dataBuffer := bytes.NewBuffer([]byte{})
	dataBuffer.Write([]byte{0x10}) // write command

	passwordBinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(passwordBinary, password)
	dataBuffer.Write(passwordBinary) // write password

	return dataBuffer.Bytes()
}

func status() {
	rss := readShortStatusCommand(30)
	frame := createFrame(rss)
	println(hex.Dump(frame))
	sendFrame(frame)
}

func sendFrame(frame []byte) {
	t := time.Now()
	con, err := net.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		log.Fatal(err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	rw.Write(frame)
	rw.Flush()

	b, err := rw.ReadByte()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Control byte:", b)
	log.Println("hex:")
	println(hex.Dump([]byte{b}))

	res2 := make([]byte, 32)
	n2, err := rw.Read(res2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("num:", n2)
	log.Println("buf:", res2)
	log.Println("hex:")
	println(hex.Dump(res2))

	rw.WriteByte(byte(0x06))
	rw.Flush()
	con.Close()
	log.Println(time.Since(t))
}

//02 | 10 | 10 | 00 | 1E 92 02 02 00 00 9F E5 18 01 00 89 08 00 | 6E
//02 | 10 | 10 | 00 | 1e 92 02 02 00 00 9f e5 18 01 00 89 08 00 | 6e
//02 | 10 | 10 | 00 | 1E 92 02 02 00 00 9F E5 18 01 00 89 08 00 | 6E

//----------------------------------------
//Краткий запрос состояния:
//----------------------------------------
//Режим:
//2, Открытая смена; 24 часа не кончились
//----------------------------------------
//Подрежим                  : 0, Бумага присутствует
//Статус режима             : 0
//Количество операций в чеке: 0
//Напряжение батареи, В     : 3,12
//Напряжение источника, В   : 25,44
//----------------------------------------
//ФлагиKKT                  : 0292h, 658
//----------------------------------------
//Увеличенная точность количества  : [нет]
//Бумага на выходе из накопителя   : [нет]
//Бумага на входе в накопитель     : [нет]
//Денежный ящик открыт             : [нет]
//Крышка корпуса поднята           : [нет]
//Рычаг термоголовки чека опущен   : [да]
//Рычаг термоголовки журнала опущен: [да]
//Оптический датчик чека           : [да]
//Оптический датчик журнала        : [да]
//2 знака после запятой в цене     : [да]
//Нижний датчик ПД                 : [да]
//Верхний датчик ПД                : [да]
//Рулон чековой ленты              : [да]
//Рулон контрольной ленты          : [да]
//----------------------------------------
//02 | 10 | 10 | 00 | 1E 92 02 02 00 00 9F E5 18 01 00 89 08 00 | 6E

//int operatorNumber = in.readByte();
//int flags = in.readShort();
//int mode = in.readByte() & 15;
//int subMode = in.readByte();
//int receiptOperationsLo = in.readByte();
//int batteryState = in.readByte();
//int powerState = in.readByte();
//int FMResultCode = in.readByte();
//int EJResultCode = in.readByte();
//int receiptOperationsHi = in.readByte();
//double batteryVoltage = (double)batteryState / 255.0D * 100.0D * 5.0D / 100.0D;
//double powerVoltage = (double)powerState * 24.0D / 216.0D * 100.0D / 100.0D;
//int receiptOperations = receiptOperationsLo + (receiptOperationsHi << 8);
