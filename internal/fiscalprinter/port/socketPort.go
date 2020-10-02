package port

import (
	"net"
	"shtrih-drv/internal/logger"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	NoConnectionError = Error("No connection error")
	ReadAnswerError   = Error("read Answer Error")
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

func NewSocketPort(conn net.Conn, logger logger.Logger) *SocketPort {
	return &SocketPort{
		conn:   conn,
		logger: logger,
	}
}

type SocketPort struct {
	conn   net.Conn
	logger logger.Logger
}

func (p SocketPort) ReadByte() (int, error) {

	//input stream
	buf := []byte{1}
	n, err := p.conn.Read(buf)
	if err != nil {
		p.logger.Fatal(err)
	}
	p.logger.Debug(n)

	//b := p.inputStream.read()
	if n == -1 {
		return 0, NoConnectionError
	}

	return n, nil
}

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
	i := 0

	for i < 2 {
		n, err := p.conn.Write(b)
		p.logger.Debug(n, b)

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
