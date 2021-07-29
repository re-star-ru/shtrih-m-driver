package main

import (
	"encoding/hex"
	"log"
)

func main() {
	data := []byte{0x1E, 0x00, 0x00, 0x00}

	// var crc byte

	// for _, v := range data {
	// 	crc ^= v
	// }

	// log.Println("crc:", crc)

	sendMessage(0x10, data)
}

const STX byte = 0x02

func sendMessage(cmdID byte, cmdData []byte) {

	N := byte(len(cmdData) + 1) // may be panic if overflow?
	m := []byte{STX, N, cmdID}

	m = append(m, cmdData...)

	m = append(m, getLRC(m[1:]))

	log.Println(len(m))

	log.Println(hex.Dump(m))

}

func getLRC(data []byte) (LRC byte) {
	for _, v := range data {
		LRC ^= v
	}

	return
}
