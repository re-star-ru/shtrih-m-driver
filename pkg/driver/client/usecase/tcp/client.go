package tcp

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/fess932/shtrih-m-driver/pkg/driver/client"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

type Usecase struct {
	host    string
	timeout time.Duration
	logger  Logger
}

type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

func NewClientUsecase(host string, timeout time.Duration, logger Logger) client.Usecase {
	return &Usecase{host: host, timeout: timeout, logger: logger}
}

func (u *Usecase) Send(frameToSend []byte, cmdLen int) (*models.Frame, error) {
	con, err := net.Dial("tcp", u.host)

	if err != nil {
		return nil, err
	}

	if err := con.SetDeadline(time.Now().Add(u.timeout)); err != nil {
		return nil, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))

	defer func() {
		if err := con.Close(); err != nil {
			u.logger.Error(err)
			return
		}
	}()

	// проверяем статус онлайн кассы, если нет ошибки отправляем фрейм с командой
	if err := u.checkPortStatus(rw); err != nil {
		u.logger.Error(err)
		return nil, err
	}

	// отправляем фрейм с командой
	if err := u.sendFrame(frameToSend, rw); err != nil {
		return nil, err
	}

	return u.receiveFrame(byte(cmdLen), rw)
}

func (u *Usecase) sendFrame(frame []byte, rw *bufio.ReadWriter) error {
	u.logger.Debug("-> send frame: \n", hex.Dump(frame))

	if _, err := rw.Write(frame); err != nil {
		u.logger.Error(err)
	}
	if err := rw.Flush(); err != nil {
		u.logger.Error(err)
	}

	b, err := rw.ReadByte() // read control byte
	u.logger.Debug("<- recive control byte:", b)

	if err != nil {
		return err
	}

	switch b {
	case models.ACK:
		return nil
	case models.NAK:
		return errors.New("ошибка интерфейса либо неверная контрольная сумма")
	default:
		return errors.New("сообщение не принято либо не верные данные")
	}
}

func (u *Usecase) checkPortStatus(rw *bufio.ReadWriter) error {

	u.logger.Debug("-> send ENQ, check port status")
	currentAttempt := 0

	for currentAttempt < models.MaxENQAttempts {
		currentAttempt++

		if err := rw.WriteByte(models.ENQ); err != nil {
			return err
		}

		if err := rw.Flush(); err != nil {
			return err
		}

		b, err := rw.ReadByte()
		if err != nil {
			return err
		}

		u.logger.Debug("<- recive control byte:", b)

		switch b {
		case models.ACK:
			u.logger.Debug("ACK, wait and retry ENQ")
			time.Sleep(models.DefaultAttemptTimeout)
			continue

		case models.NAK:
			u.logger.Debug("NAK, send cmd")
			return nil

		default:
			u.logger.Error("ERR, ping byte:", b)
			return errors.New(fmt.Sprintln("ERR, ping byte:", b))
		}
	}

	return errors.New("no answer")
}

func (u *Usecase) ping(rw *bufio.ReadWriter) {
	u.logger.Debug("-> send ENQ")

	if err := rw.WriteByte(models.ENQ); err != nil {
		u.logger.Error(err)
	}

	if err := rw.Flush(); err != nil {
		u.logger.Error(err)
	}

	b, _ := rw.ReadByte()
	u.logger.Debug("<- recive control byte:", b)

	switch b {
	case models.ACK:
		u.logger.Debug("OK, ACK, wait for recive now")
		rw.WriteByte(models.ACK)
		rw.Flush()
	case models.NAK:
		u.logger.Debug("OK, NAK, wait for cmd now")
		rw.WriteByte(models.ACK)
		rw.Flush()
	default:
		rw.WriteByte(models.ACK)
		rw.Flush()
		u.logger.Fatal("ERR, ping byte:", b)
	}
}

func (u *Usecase) receiveFrame(cmdLen byte, rw *bufio.ReadWriter) (*models.Frame, error) {
	u.logger.Debug("<- start receive frame")

	defer func() {
		rw.WriteByte(models.ACK)
		rw.Flush()
	}()

	var FRM models.Frame
	var err error

	FRM.STX, err = rw.ReadByte() // read byte STX (0x02) err need
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	FRM.DLEN, err = rw.ReadByte() // read byte dataLen

	FRM.CMD = make([]byte, cmdLen)
	rw.Read(FRM.CMD) // read cmd bytes

	FRM.ERR, _ = rw.ReadByte() // read err byte

	FRM.DATA = make([]byte, FRM.DLEN-cmdLen-models.ErrLen)
	rw.Read(FRM.DATA) // read data bytes

	FRM.CRC, _ = rw.ReadByte() // read crc byte

	u.logger.Debug("<- end recive frame: \n",
		fmt.Sprintf("stx: %v, dlen: %v, crc: %v  \n", FRM.STX, FRM.DLEN, FRM.CRC),
		hex.Dump(FRM.Bytes()))

	dataCheck := bytes.NewBuffer([]byte{})
	dataCheck.Write(FRM.CMD)
	dataCheck.WriteByte(FRM.ERR)
	dataCheck.Write(FRM.DATA)
	if err := FRM.CheckCRC(); err != nil {
		return nil, err
	}

	return &FRM, nil
}
