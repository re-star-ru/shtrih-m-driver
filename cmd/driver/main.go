package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("golang shtrih")

	//timeout := time.Second * 1

	conn, err := net.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()

	for {
		source := ""
		log.Println("Запрос серверу")
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			log.Println(err)
			return
		}
		// получаем ответ
		log.Println("Ответ::")
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				break
			}
			log.Println(string(buff[0:n]))
			conn.SetReadDeadline(time.Now().Add(time.Microsecond * 700))
		}
		log.Println()

	}

	//conn.SetDeadline(time.Now().Add(timeout))
	//
	//
	//go writeByte(conn)
	//
	//
	//conn.SetDeadline(time.Now().Add(time.Second * 5))
	//buf := make([]byte, 1)
	//b, err := conn.Read(buf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(b)

	//fp := fiscalprinter.NewFiscalPrinter()

	//log.Println(fp.GetSerialNumber())
	//time.Sleep(time.Second * 5)
}

func writeByte(conn net.Conn) {
	wb := make([]byte, 5)
	res, err := conn.Write(wb)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
