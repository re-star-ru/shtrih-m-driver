package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

//buf.WriteByte(2) // start
//buf.WriteByte(byte(len(data))) // len data
//buf.Write(data) // data
//buf.WriteByte(f.GetCrc(data)) // control sum (crc)

//02 05 10 1E 00 00 00 0B
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
	//status()
}

func status() {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(0x02)           // write start
	buf.Write([]byte{0x05, 0x10}) // write command

	password := make([]byte, 4)
	binary.LittleEndian.PutUint32(password, 30)

	buf.Write(password) // write password

	buf.Write([]byte{0}) // write control sum
	println(hex.Dump(buf.Bytes()))

	//t := time.Now()
	//con, err := net.Dial("tcp", "10.51.0.71:7778")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	//rw.Write(buf.Bytes())
	//rw.Flush()
	//
	//b, err := rw.ReadByte()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Control byte:", b)
	//log.Println("hex:")
	//println(hex.Dump([]byte{b}))
	//
	//res2 := make([]byte, 32)
	//n2, err := rw.Read(res2)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("num:", n2)
	//log.Println("buf:", res2)
	//log.Println("hex:")
	//println(hex.Dump(res2))
	//
	//rw.Write([]byte{06})
	//rw.Flush()
	//
	//con.Close()
	//log.Println(time.Since(t))
}
