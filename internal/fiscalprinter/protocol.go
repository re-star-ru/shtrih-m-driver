package fiscalprinter

import (
	"encoding/hex"
	"shtrih-drv/internal/fiscalprinter/command"
	"shtrih-drv/internal/fiscalprinter/port"
	"shtrih-drv/internal/logger"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func NewPrinterProtocol(logger logger.Logger) *PrinterProtocol {
	// константы из driver protocol

	//p.byteTimeout = 100
	//STX := byte(2);
	//ENQ :=byte(5) ;
	//ACK :=byte(6) ;
	//NAK :=byte(21) ;

	//p.maxEnqNumber = 3

	//private int maxNakCommandNumber = 3;
	//p.maxNakAnswerNumber = 3
	//p.maxAckNumber = 3
	//private int maxRepeatCount = 3;
	//private byte[] txData = new byte[0];
	//private byte[] rxData = new byte[0];

	return &PrinterProtocol{
		logger:             logger,
		byteTimeout:        1000,
		maxEnqNumber:       3,
		maxNakAnswerNumber: 3,
		maxAckNumber:       3,
	}
}

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

	if err := p.portWrite(ENQ); err != nil {
		p.logger.Fatal(err)
	}

	B, err := p.portReadByte()
	if err != nil {
		return err
	}

	switch B {
	case ACK:
		p.readAnswer(p.byteTimeout)
		ackNumber++
	case NAK:
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

func (p *PrinterProtocol) readControlByte() (int, error) {
	var result int
	var err error

	p.logger.Debug("read control byte")
	for {
		result, err = p.portReadByte()
		if err != nil {
			return 0, err
		}
		if result != 255 {
			break
		}
	}
	return result, nil
}

func (p *PrinterProtocol) portWrite(b int) error {
	data := []byte{byte(b)}
	return p.portWriteData(data)
}

func (p *PrinterProtocol) portWriteData(b []byte) error {
	//Logger2.logTx(logger, data)

	p.logger.Debug("-> ", "\n", hex.Dump(b))
	_, err := p.client.W.Write(b)
	if err != nil {
		return err
	}
	return p.client.W.Flush()
}

func (p *PrinterProtocol) portReadByte() (int, error) {
	b, err := p.client.R.ReadByte()
	if err != nil {
		return 0, err
	}
	p.logger.Debug("<- ", "\n", hex.Dump([]byte{byte(b)}))

	return int(b), nil
}

func (p *PrinterProtocol) portReadBytes(l int) []byte {
	p.logger.Debug("port read bytes start")

	buf := make([]byte, l)

	b, err := p.client.R.Read(buf)
	p.logger.Debug("port reads bytes: ", b)

	if err != nil {
		p.logger.Fatal(err)
	}
	p.logger.Debug("<- ", "\n", hex.Dump(buf))
	return buf
}

func (p *PrinterProtocol) readAnswer(timeout int) ([]byte, error) {
	enqNumber := 0
	nakCount := 0
	p.logger.Debug("read answer")

label36: // TODO: убрать это нахер из моего кода
	for {
		//port.Set(timeout)

		//for (p.portReadByte() != 2) {}
		for {
			b, err := p.portReadByte()
			if err != nil {
				return nil, err
			}
			if b == STX {
				break
			}
		}

		//this.port.setTimeout(this.byteTimeout);

		dataLength, err := p.portReadByte()
		if err != nil {
			return nil, err
		}
		dataLength++
		commandData := p.portReadBytes(dataLength)
		//p.logger.Debug(commandData)

		crc := commandData[len(commandData)-1]

		commandData = commandData[:len(commandData)-1]

		if p.frame.GetCrc(commandData) == crc {
			if err := p.portWrite(ACK); err != nil {
				return nil, err
			}
			return commandData, nil
		}

		if nakCount >= p.maxNakAnswerNumber {
			return nil, port.ReadAnswerError
		}

		nakCount++
		if err := p.portWrite(NAK); err != nil {
			return nil, err
		}

		for {
			if err := p.portWrite(ENQ); err != nil {
				return nil, err
			}
			enqNumber++
			B, err := p.readControlByte()
			if err != nil {
				return nil, err
			}
			if B == ACK {
				continue label36
			}

			if B == NAK {
				return nil, port.ReadAnswerError
			}

			if enqNumber < p.maxEnqNumber {
				break
			}
		}

		return nil, port.ReadAnswerError
	}
}

func (p *PrinterProtocol) SendCommand(command command.PrinterCommander) error {
	p.logger.Debug("send command: ", command.GetText())
	tx, _ := command.EncodeData()

	rx, err := p.sendEncodedCommand(tx, p.byteTimeout)
	if err != nil {
		return err
	}

	p.logger.Debug("rx: ", rx)
	rxe, err := charmap.Windows1251.NewDecoder().Bytes(rx)
	if err != nil {
		p.logger.Fatal(err)
	}
	p.logger.Debug("txe: ", string(rxe))

	return nil
}

func (p *PrinterProtocol) sendEncodedCommand(data []byte, timeout int) ([]byte, error) {
	var err error
	p.txData, err = p.frame.encode(data)
	if err != nil {
		return nil, err
	}
	//p.logger.Debug("send encoded command after encode: ", hex.Dump(p.txData))

	rx, err := p.send(p.txData)
	if err != nil {
		return nil, err
	}
	//p.logger.Debug("recive encoded command after encode: ", hex.Dump(rx))

	p.rxData, err = p.frame.encode(rx)
	//p.logger.Debug("recive encoded command after encode: ", hex.Dump(p.rxData))

	//private byte[] sendCommand(byte[] data, int timeout) throws Exception {
	//	this.txData = this.frame.encode(data);
	//	byte[] rx = this.send(this.txData, timeout);
	//	this.rxData = this.frame.encode(rx);
	//	return rx;
	//}
	return rx, err
}

func (p *PrinterProtocol) send(data []byte) ([]byte, error) {
	var ackNumber, enqNumber, B int
	var err error

	if err := p.portWrite(ENQ); err != nil {
		return nil, err
	}
	enqNumber++

	B, err = p.readControlByte()
	if err != nil {
		return nil, err
	}

	switch B {

	case ACK:
		if _, err := p.readAnswer(p.byteTimeout); err != nil {
			return nil, err
		}
		ackNumber++
		break
	case NAK:
		if err := p.writeCommand(data); err != nil {
			return nil, err
		}
		return p.readAnswer(p.byteTimeout)
	default:
		time.Sleep(time.Millisecond * 100) // TODO: wtf
	}

	return nil, nil
}

func (p *PrinterProtocol) writeCommand(data []byte) error {

	if err := p.portWriteData(data); err != nil {
		return err
	}

	B, err := p.readControlByte()
	if err != nil {
		return err
	}

	switch B {
	case ACK:
		return nil
	case NAK:
		break
	default:
		return nil
	}

	return nil
}
