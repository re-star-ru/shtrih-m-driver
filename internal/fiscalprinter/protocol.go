package fiscalprinter

import (
	"log"
	"net"
	"shtrih-drv/internal/fiscalprinter/port"
	"shtrih-drv/internal/logger"
	"time"
)

func NewPrinterProtocol(logger logger.Logger) *PrinterProtocol {
	return &PrinterProtocol{
		logger: logger,
	}
}

type PrinterProtocol struct {
	byteTimeout        int
	maxEnqNumber       int
	maxNakAnswerNumber int
	maxAckNumber       int

	port   *port.SocketPort
	frame  Frame
	logger logger.Logger
}

func (p *PrinterProtocol) Connect() error {
	// константы из printer protocol

	p.byteTimeout = 100
	//STX := byte(2);
	//ENQ :=byte(5) ;
	//ACK :=byte(6) ;
	//NAK :=byte(21) ;

	p.maxEnqNumber = 3

	//private int maxNakCommandNumber = 3;
	p.maxNakAnswerNumber = 3
	p.maxAckNumber = 3
	//private int maxRepeatCount = 3;
	//private byte[] txData = new byte[0];
	//private byte[] rxData = new byte[0];

	//
	ackNumber := 0
	enqNumber := 0

	conn, err := net.Dial("tcp", "10.51.0.71:7778")
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Millisecond * 700))
	p.port = port.NewSocketPort(conn, p.logger)

	for {
		//port.setTimeout(byteTimeout)
		if err := p.port.Write([]byte{5}); err != nil {
			log.Fatal(err)
		}
		B := p.readControlByte()
		switch B {
		case 6:
			p.readAnswer(p.byteTimeout)
			ackNumber++
		case 21:
			return nil
		default:
			time.Sleep(time.Millisecond * 100)
			enqNumber++
		}
		//catch (IOException var4) {
		//		enqNumber++;
		//}

		if ackNumber >= p.maxAckNumber {
			//throw new DeviceException(2, Localizer.getString("NoConnection"));
			return port.NoConnectionError
		}

		if enqNumber < p.maxEnqNumber {
			break
		}
	}

	return port.NoConnectionError
}

func (p *PrinterProtocol) readControlByte() int {
	result := 0
	for {
		result = p.portReadByte()
		if result == 225 {
			break
		}
	}
	return result
}

func (p *PrinterProtocol) portReadByte() int {
	b, err := p.port.ReadByte()
	if err != nil {
		p.logger.Fatal(err)
	}
	return b
}

func (p *PrinterProtocol) portReadBytes(len int) []byte {
	b, err := p.port.ReadBytes(len)
	if err != nil {
		p.logger.Fatal(err)
	}
	return b
}

func (p *PrinterProtocol) readAnswer(timeout int) ([]byte, error) {
	enqNumber := 0
	nakCount := 0

label36: // TODO: убрать это нахер из моего кода
	for {
		//port.Set(timeout)

		//for (p.portReadByte() != 2) {}

		//this.port.setTimeout(this.byteTimeout);

		dataLength := p.portReadByte() + 1
		commandData := p.portReadBytes(dataLength)

		crc := commandData[len(commandData)-1]

		commandData = commandData[:len(commandData)-1]

		if p.frame.GetCrc(commandData) == crc {
			p.port.Write([]byte{6})
			return commandData, nil
		}

		if nakCount >= p.maxNakAnswerNumber {
			return nil, port.ReadAnswerError
		}

		nakCount++
		p.port.Write([]byte{21})

		for {
			p.port.Write([]byte{5})
			enqNumber++
			B := p.readControlByte()
			if B == 6 {
				continue label36
			}

			if B == 21 {
				return nil, port.ReadAnswerError
			}

			if enqNumber < p.maxEnqNumber {
				break
			}
		}

		return nil, port.ReadAnswerError
	}
}
