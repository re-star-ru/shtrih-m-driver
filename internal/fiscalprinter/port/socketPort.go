package port

import (
	"bufio"
	"net"
	"time"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	NoConnectionError = Error("No connection error")
	ReadAnswerError   = Error("read Answer Error")
	DataLenghtExeeds  = Error("Data length exeeds 256 bytes")
)

//package com.shtrih.fiscalprinter.port; todo

//_, err := net.Dial("tcp", "10.51.0.71:7778")
//if err != nil {
//	log.Println(err.Error())
//}
//type SocketPort struct {
//	PortNameport
//	PortNameport
//}

type TcpClient struct {
	R    *bufio.Reader
	W    *bufio.Writer
	buf  []byte
	Conn net.Conn
}

func Connect(host string) (TcpClient, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return TcpClient{}, err
	}

	conn.SetDeadline(time.Now().Add(time.Millisecond * 2000))

	return TcpClient{
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
		make([]byte, 1024),
		conn,
	}, err
}

func NewSocketPort(conn net.Conn) *SocketPort {
	return &SocketPort{
		conn: conn,
	}
}

type SocketPort struct {
	conn net.Conn
}

//func (p SocketPort) ReadByte() (int, error) {
//
//	buf := make([]byte, 0, 4096) // big buffer
//	tmp := make([]byte, 256)
//
//	for {
//		n, err := p.conn.Read(tmp)
//		if err != nil {
//			if err != io.EOF {
//				log.Println("read error:", err)
//			}
//			break
//		}
//		//fmt.Println("got", n, "bytes.")
//		buf = append(buf, tmp[:n]...)
//	}
//	log.Println("total size:", len(buf))
//
//	////input stream
//	//buf := []byte{0}
//	//n, err := p.conn.Read(buf)
//	//if err != nil {
//	//	return 0, err
//	//}
//	//
//	////b := p.inputStream.read()
//	//if n == -1 {
//	//	return 0, NoConnectionError
//	//}
//	//
//	return len(buf), nil
//}

func (p SocketPort) ReadBytes(len int) ([]byte, error) {
	data := make([]byte, len)

	count := 0

	for offset := 0; len > 0; offset += count {
		count, err := p.conn.Read(data)

		if err != nil {
			return nil, NoConnectionError
		}
		if count == -1 {
			return nil, NoConnectionError
		}

		len -= count
	}

	return data, nil
}

func (p SocketPort) Write(b []byte) error {

	//s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//f := os.NewFile(uintptr(s))

	//p.conn.

	i := 0
	for i < 2 {
		_, err := p.conn.Write(b)
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

//func (p SocketPort) open(timeout int) {
//	conn, err := net.Dial("tcp", "10.51.0.71:7778")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//if (!this.isConnected()) {
//	//	this.socket = new Socket();
//	//	this.socket.setReuseAddress(true);
//	//	this.socket.setSoTimeout(this.readTimeout);
//	//	this.socket.setTcpNoDelay(true);
//	//	StringTokenizer tokenizer = new StringTokenizer(this.portName, ":");
//	//	String host = tokenizer.nextToken();
//	//	int port = Integer.parseInt(tokenizer.nextToken());
//	//	this.socket.connect(new InetSocketAddress(host, port), this.openTimeout);
//	//	this.inputStream = this.socket.getInputStream();
//	//	this.outputStream = this.socket.getOutputStream();
//	//}
//}
