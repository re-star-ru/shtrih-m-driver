package main

import (
	"bufio"

	"log"
	"net"
	"shtrih-drv/internal/fiscalprinter"
	"shtrih-drv/internal/fiscalprinter/command"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func main() {

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal()
	}
	slogger := logger.Sugar()

	slogger.Info("Shtih driver starting")
	//slogger.Debug("This is a DEBUG message")
	//slogger.Info("This is an INFO message")
	////slogger.Info("This is an INFO message with fields", "region", "us-west", "id", 2)
	//slogger.Warn("This is a WARN message")
	//slogger.Error("This is an ERROR message")

	printer := fiscalprinter.NewPrinterProtocol(slogger)
	err = printer.Connect()
	if err != nil {
		slogger.Error(err)
		return
	}

	//code := 300
	//log.Println(code >> 8 & 255)
	//log.Println(code & 255)

	c := command.NewReadDeviceMetrics()
	if err := printer.SendCommand(c); err != nil {
		slogger.Fatal(err)
	}

	t := time.Now()
	ls := command.NewReadLongStatus()
	if err := printer.SendCommand(ls); err != nil {
		slogger.Fatal(err)
	}
	slogger.Debug(time.Since(t))

	//ReadPrinterModelParameters next

	//conn, err := Connect("10.51.0.71:7778")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := conn.w.WriteByte(5); err != nil {
	//	log.Fatal(err)
	//}
	//conn.w.Flush()
	//
	//b, err := conn.r.ReadByte()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(b)
	//
	//time.Sleep(time.Millisecond * 100)
}

type TcpClient struct {
	r    *bufio.Reader
	w    *bufio.Writer
	buf  []byte
	conn net.Conn
}

func Connect(host string) (TcpClient, error) {
	log.Println("Connect")
	conn, err := net.Dial("tcp", host)
	conn.SetDeadline(time.Now().Add(time.Millisecond * 1000))

	return TcpClient{
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
		make([]byte, 1024),
		conn,
	}, err
}

//
//func (tcp TcpClient) Send(dataSend string) (string, error) {
//	var (
//		returnData string
//		returnErr  error
//	)
//	tcp.w.WriteString(dataSend + "\r\n")
//	tcp.w.Flush()
//
//ILOOP:
//	for {
//		n, err := tcp.r.Read(tcp.buf)
//		data := string(tcp.buf[:n])
//		switch err {
//		case io.EOF:
//			break ILOOP
//		case nil:
//			returnData = data
//			returnErr = nil
//			if isTransportOver(data) {
//				break ILOOP
//			}
//		default:
//			returnData = ""
//			returnErr = err
//		}
//	}
//	return returnData, returnErr
//}
//func (tcp TcpClient) Close() {
//	tcp.conn.Close()
//}
//func isTransportOver(data string) (over bool) {
//	over = strings.HasSuffix(data, "\r\n")
//	return
//}
//
////func read(cn net.Conn) {
////	log.Println("read...")
////	var buf []byte
////
////	for {
////		n, err := cn.Read(buf)
////		log.Println(buf, n, err)
//
//	}
//}
