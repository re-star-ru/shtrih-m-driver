package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	ENQ byte = 0x05
	STX byte = 0x02
	ACK byte = 0x06
	NAK byte = 0x15
)

const SHORT_STATUS byte = 0x10

func main() {
	log.Println("dial to kkt")
	conn, err := net.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		log.Fatal("err dial dial:", err)
	}

	buf := bufio.NewReader(conn)

	canSendCmd, err := ping(conn, buf)
	if err != nil {
		log.Fatal("No connection:", err)
	}
	log.Println("Can send cmd to kkt:", canSendCmd)

	if canSendCmd {
		if err := sendCmd(conn); err != nil {
			log.Fatal("err while sendCmd", err)
		}
	}

	log.Println("read bytes from buffered reader")
	for {
		b, err := buf.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("byte:", b)
	}
}

func sendCmd(w io.Writer) error {
	n, err := w.Write([]byte{SHORT_STATUS})
	if err != nil {
		return err
	}
	log.Println("bytes writed:", n)

	return nil
}

// Служебное сообщение
func ping(w io.Writer, buf *bufio.Reader) (bool, error) {
	// write enq for req message
	n, err := w.Write([]byte{ENQ}) // write timeout
	if err != nil {
		return false, nil
	}
	log.Println("bytes writed to conn:", n)

	b, err := buf.ReadByte()
	if err != nil {
		return false, err
	}

	log.Println("readed byte from buf:", b)

	switch b {
	case NAK:
		return true, nil
	case ACK:
		return false, fmt.Errorf("kkt busy")
	default:
		return false, fmt.Errorf("wrong byte for ENQ")
	}
}

type Message struct {
}

func sendMessage(cmdID byte, cmdData []byte) {

	N := byte(len(cmdData)) // may be panic if overflow?
	m := []byte{STX, N, cmdID}

	m = append(m, cmdData...)

	m = append(m, getLRC(m[1:]))

	log.Println(len(m))

}

func getLRC(data []byte) (LRC byte) {
	for _, v := range data {
		LRC ^= v
	}

	return
}

//1 struct with fields with render end

//2 buffer with writed bytes
