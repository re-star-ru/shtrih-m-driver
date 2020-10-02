package fiscalprinter

import (
	"encoding/hex"
	"shtrih-drv/internal/fiscalprinter/command"
	"shtrih-drv/internal/fiscalprinter/port"
	"shtrih-drv/internal/logger"
	"time"
)

func NewPrinterProtocol(logger logger.Logger) *PrinterProtocol {
	return &PrinterProtocol{
		logger: logger,
	}
}

const (
	ACKNOWLEDGE          = 0x6
	NEGATIVE_ACKNOWLEDGE = 0x15
)

type PrinterProtocol struct {
	byteTimeout        int
	maxEnqNumber       int
	maxNakAnswerNumber int
	maxAckNumber       int

	client port.TcpClient
	frame  Frame
	logger logger.Logger

	txData []byte
	rxData []byte
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

	p.logger.Debug("Connect")
	con, err := port.Connect("10.51.0.71:7778")
	if err != nil {
		return err
	}
	p.client = con
	//defer conn.Close()

	//port.setTimeout(byteTimeout)
	p.portWrite(5)
	if err := con.W.WriteByte(5); err != nil {
		p.logger.Fatal(err)
	}
	con.W.Flush()

	B, err := con.R.ReadByte()
	if err != nil {
		return err
	}
	p.logger.Debug(B)

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

	if ackNumber >= p.maxAckNumber {
		return port.NoConnectionError
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

func (p *PrinterProtocol) portWrite(b int) error {
	//Logger2.logTx(logger, data)

	data := []byte{byte(b)}
	p.logger.Debug("-> ", hex.Dump(data))
	return p.port.Write(data)
}

func (p *PrinterProtocol) portReadByte() int {
	b, err := p.port.ReadByte()
	if err != nil {
		p.logger.Fatal(err)
	}
	p.logger.Debug("<- ", hex.Dump([]byte{byte(b)}))

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

func (p *PrinterProtocol) sendCommand(command command.PrinterCommand) {
	p.logger.Debug("send command: ", command.GetText())

	tx := command.EncodeData()

	rx := p.sendEncodedCommand(tx, p.byteTimeout)
	p.logger.Debug(rx)
}

func (p *PrinterProtocol) sendEncodedCommand(data []byte, timeout int) []byte {
	p.logger.Debug("send encoded command: ", hex.Dump(data))

	var err error
	p.txData, err = p.frame.encode(data)
	if err != nil {
		p.logger.Fatal(err)
	}
	//private byte[] sendCommand(byte[] data, int timeout) throws Exception {
	//	this.txData = this.frame.encode(data);
	//	byte[] rx = this.send(this.txData, timeout);
	//	this.rxData = this.frame.encode(rx);
	//	return rx;
	//}
	return nil
}
