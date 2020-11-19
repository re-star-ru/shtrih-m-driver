package printerUsecase

import (
	"encoding/hex"
	"errors"
	"unicode/utf8"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"golang.org/x/text/encoding/charmap"
)

func (p *printerUsecase) WriteCashierINN(INN string) error {
	buf, cmdLen := p.createCommandBuffer(models.FnWriteTLV, p.password)

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
		p.logger.Error(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

	//
	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		return err
	}

	return nil
}
