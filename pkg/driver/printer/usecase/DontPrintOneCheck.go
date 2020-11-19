package printerUsecase

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

func (p *printerUsecase) DontPrintOneCheck() {
	buf := new(bytes.Buffer)
	buf.WriteByte(1) // значение 1

	p.writeTable(17, 1, 7, buf)
}

func (p *printerUsecase) writeTable(tableNumber byte, rowNumber uint16, fieldNumber byte, fieldValue models.FieldValue) {
	buf, cmdLen := p.createCommandBuffer(models.WriteTable, p.password)

	buf.WriteByte(tableNumber) // номер таблицы

	rowNumBin := make([]byte, 2)
	binary.LittleEndian.PutUint16(rowNumBin, rowNumber)
	buf.Write(rowNumBin) // номер ряда

	buf.WriteByte(fieldNumber) // номер поля

	buf.Write(fieldValue.Bytes()) // запись поля

	rFrame, err := p.send(buf.Bytes(), cmdLen)
	if err != nil {
		p.logger.Error(err)
		return
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

}
