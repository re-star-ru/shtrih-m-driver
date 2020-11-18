package printerUsecase

import (
	"encoding/hex"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"golang.org/x/text/encoding/charmap"
)

func (p *printerUsecase) CloseCheck(chk models.CheckPackage) {
	p.logger.Debug("Send command CloseCheck")

	buf, cmdLen := p.createCommandBuffer(models.CloseCheckV2, p.password)
	p.logger.Debug("cmdlen:", cmdLen)

	// запись суммы наличных - типа оплаты 1
	cash, err := p.intToBytesWithLen(chk.Cash, 5)
	if err != nil {
		p.logger.Error(err)
		return
	}
	buf.Write(cash)

	// запись суммы типа оплаты 2 - безнал
	casheless, err := p.intToBytesWithLen(chk.Casheless, 5)
	if err != nil {
		p.logger.Error(err)
		return
	}
	buf.Write(casheless)

	// запись остальных с 3 по 16 видоов оплаты, длинна вида 5 байт
	for i := 2; i < 16; i++ {
		buf.Write(make([]byte, 5)) // 3 - 16
	}

	buf.WriteByte(0) // округление до рубля

	for i := 0; i < 5; i++ {
		buf.Write(make([]byte, 5)) // налог 1-6
	}

	buf.WriteByte(chk.TaxSystem) //система налогоообложения, биты а не байт

	// Запись нижней линии чека 0 - 128 байт строка
	str, err := charmap.Windows1251.NewEncoder().String(chk.BottomLine)
	if err != nil {
		p.logger.Error(err)
		return
	}
	rStrBytes := make([]byte, 64)
	copy(rStrBytes, str)

	buf.Write(rStrBytes[:64])

	p.logger.Debug("len: ", buf.Len())

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Error(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))
}
