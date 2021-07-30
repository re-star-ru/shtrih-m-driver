package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	ENQ byte = 0x05
	STX byte = 0x02
	ACK byte = 0x06 // 6
	NAK byte = 0x15 // 21
)

const SHORT_STATUS byte = 0x10

func main() {
	log.Println("dial to kkt")

	var d net.Dialer

	d.Timeout = time.Second * 5 // todo retry dial

	conn, err := d.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		log.Fatal("err dial dial:", err)
	}
	log.Println("dial 1 ok")

	conn2, err := d.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		log.Fatal("err dial dial:", err)
	}
	log.Println("dial 2 ok")

	go func() {
		canSendCmd, err := ping(conn)
		if err != nil {
			log.Fatal("No connection:", err)
		}

		log.Println("Can send cmd conn to kkt:", canSendCmd)
	}()

	go func() {
		canSendCmd, err := ping(conn2)
		if err != nil {
			log.Fatal("No connection:", err)
		}

		log.Println("Can send cmd to conn2 kkt:", canSendCmd)
	}()

	time.Sleep(time.Second * 5)
	log.Println(conn, conn2)

	//buf := bufio.NewReader(conn)

	//for i := 0; i < 5; i++ {
	//	go func() {
	//		for {
	//			time.Sleep(time.Second * 5)
	//			canSendCmd, err := ping(conn)
	//			if err != nil {
	//				log.Fatal("No connection:", err)
	//			}
	//
	//			log.Println("Can send cmd to kkt:", canSendCmd)
	//		}
	//	}()
	//}

	//{
	//	time.Sleep(time.Second)
	//	canSendCmd, err := ping(conn, buf)
	//	if err != nil {
	//		log.Fatal("No connection:", err)
	//	}
	//	log.Println("Can send cmd to kkt:", canSendCmd)
	//}

	//if canSendCmd {
	//	if err := sendCmdShortStatus(conn); err != nil {
	//		log.Fatal("err while sendCmd", err)
	//	}
	//}
	//
	//log.Println("read bytes from buffered reader")
	//for {
	//	b, err := buf.ReadByte()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	log.Println("byte:", b)
	//}
}

func sendCmdShortStatus(w io.Writer) error {
	return sendMessage(w, SHORT_STATUS, []byte{0x1E, 0x00, 0x00, 0x00})
}

const pingDeadline = time.Millisecond * 1000

// Служебное сообщение
func ping(conn net.Conn) (bool, error) {

	if err := conn.SetDeadline(time.Now().Add(pingDeadline)); err != nil {
		log.Fatal(err)
	}

	n, err := conn.Write([]byte{ENQ})
	// write enq for req message
	if err != nil {
		return false, nil
	}
	log.Println("bytes writed to conn:", n)

	b := make([]byte, 1)
	n, err = conn.Read(b)
	if err != nil {
		return false, err
	}

	log.Println("readed bytes from buf:", n)

	switch b[0] {
	case NAK:
		return true, nil
	case ACK:
		return false, fmt.Errorf("kkt busy")
	default:
		return false, fmt.Errorf("wrong byte for ENQ")
	}
}

func sendMessage(w io.Writer, cmdID byte, cmdData []byte) error {

	N := byte(len(cmdData)) // may be panic if overflow?
	m := []byte{STX, N, cmdID}

	m = append(m, cmdData...)

	m = append(m, getLRC(m[1:]))

	log.Println(len(m))

	n, err := w.Write(m)
	if err != nil {
		return err
	}
	log.Println("bytes writed:", n)

	return nil
}

func getLRC(data []byte) (LRC byte) {
	for _, v := range data {
		LRC ^= v
	}

	return
}

//1 struct with fields with render end

//2 buffer with writed bytes
