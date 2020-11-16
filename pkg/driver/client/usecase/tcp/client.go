package tcp

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/fess932/shtrih-m-driver/pkg/consts"

	"github.com/fess932/shtrih-m-driver/pkg/driver/client"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/fess932/shtrih-m-driver/pkg/logger"
)

type Usecase struct {
	host   string
	logger logger.Logger
}

func (u *Usecase) Send(frameToSend []byte, cmdLen int) (*models.Frame, error) {
	con, err := net.Dial("tcp", u.host)
	if err != nil {
		return nil, err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))

	defer func() {
		if err := con.Close(); err != nil {
			u.logger.Error(err)
			return
		}
	}()

	if err := u.sendFrame(frameToSend, con, rw); err != nil {
		return &models.Frame{}, err
	}

	// TODO: send close connection
	return u.receiveFrame(con, byte(cmdLen), rw)
}

func (u *Usecase) sendFrame(frame []byte, con net.Conn, rw *bufio.ReadWriter) error {
	u.ping(rw, con)

	if _, err := rw.Write(frame); err != nil {
		u.logger.Error(err)
	}
	if err := rw.Flush(); err != nil {
		u.logger.Error(err)
	}

	u.logger.Debug("-> send frame: \n", hex.Dump(frame))

	b, err := rw.ReadByte() // read control byte
	u.logger.Debug("<- recive control byte:", b)

	if err != nil {
		return err
	}

	switch b {
	case consts.ACK:
		return nil
	case consts.NAK:
		return errors.New("ошибка интерфейса либо неверная контрольная сумма")
	default:
		return errors.New("сообщение не принято либо не верные данные")
	}
}

func (u *Usecase) ping(rw *bufio.ReadWriter, con net.Conn) {
	u.logger.Debug("-> send ENQ")

	if err := rw.WriteByte(consts.ENQ); err != nil {
		u.logger.Error(err)
	}

	if err := rw.Flush(); err != nil {
		u.logger.Error(err)
	}

	b, _ := rw.ReadByte()
	u.logger.Debug("<- recive control byte:", b)

	switch b {
	case consts.ACK:
		u.logger.Debug("OK, ACK, wait for recive now")
		rw.WriteByte(consts.ACK)
		rw.Flush()
	case consts.NAK:
		u.logger.Debug("OK, NAK, wait for cmd now")
		rw.WriteByte(consts.ACK)
		rw.Flush()
	default:
		rw.WriteByte(consts.ACK)
		rw.Flush()
		u.logger.Fatal("ERR, ping byte:", b)
	}
}

func (u *Usecase) receiveFrame(con net.Conn, cmdLen byte, rw *bufio.ReadWriter) (*models.Frame, error) {
	u.logger.Debug("<- Receive frame")

	defer func() {
		rw.WriteByte(consts.ACK)
		rw.Flush()
	}()

	var FRM models.Frame
	var err error

	FRM.STX, err = rw.ReadByte() // read byte STX (0x02) err need
	if err != nil {
		u.logger.Fatal(err)
	}

	FRM.DLEN, err = rw.ReadByte() // read byte dataLen

	FRM.CMD = make([]byte, cmdLen)
	rw.Read(FRM.CMD) // read cmd bytes

	FRM.ERR, _ = rw.ReadByte() // read err byte

	FRM.DATA = make([]byte, FRM.DLEN-cmdLen-models.ErrLen)
	rw.Read(FRM.DATA) // read data bytes

	FRM.CRC, _ = rw.ReadByte() // read crc byte

	u.logger.Debug("<- recive frame: \n",
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

func NewClientUsecase(host string, logger logger.Logger) client.Usecase {
	return &Usecase{host: host, logger: logger}
}
