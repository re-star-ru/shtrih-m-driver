package printer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
)

func (p *Printer) TLVWriteCashierINN(INN string) error {
	cmdBinary, cmdLen := p.createCommandData(consts.FnWriteTLV)
	buf := bytes.NewBuffer(cmdBinary)
	p.logger.Debug(buf)

	if utf8.RuneCountInString(INN) > int(consts.CashierINN.Length) {
		return errors.New("cлишком длинное инн")
	}

	cpStr, err := charmap.CodePage866.NewEncoder().String(INN)
	if err != nil {
		return err
	}

	p.logger.Debug(cpStr)

	// tlv структура
	tlv, err := newTLV(consts.CashierINN.Code, consts.CashierINN.Length, []byte(cpStr))
	if err != nil {
		return err
	}
	buf.Write(tlv.Tag)
	buf.Write(tlv.Len)
	buf.Write(tlv.Value)

	p.logger.Debug("Код, длинна, значение:", tlv.Tag, tlv.Len, tlv.Value)
	p.logger.Debug("Команда с тлв структурой\n", hex.Dump(buf.Bytes()))

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.bytes()))

	//
	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		return err
	}

	return nil
}

type TLV struct {
	Tag   []byte
	Len   []byte
	Value []byte
}

func newTLV(Tag, Len uint16, value []byte) (TLV, error) {
	tlv := TLV{
		Tag:   make([]byte, 2),
		Len:   make([]byte, 2),
		Value: make([]byte, Len),
	}
	binary.LittleEndian.PutUint16(tlv.Tag, Tag) // код тега
	binary.LittleEndian.PutUint16(tlv.Len, Len) // длинна тега

	copy(tlv.Value, value) // значение тега

	if len(tlv.Value) != int(Len) {
		return TLV{}, errors.New("длинна не совпадает со значением")
	}

	return tlv, nil
}
