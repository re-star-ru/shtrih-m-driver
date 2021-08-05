package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"golang.org/x/text/encoding/charmap"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
)

// KKT Commands
const (
	ShortStatus byte = 0x10
	ZReport     byte = 0x41
	CancelCheck byte = 0x88
	// сброс состояния сделать или
	// отмена чека
)

// Fn Commands Start with FF
const (
	FNCommand byte = 0xFF

	FnOperationV2      byte = 0x46
	FnCloseCheckV2     byte = 0x45
	FnCancelFNDocument byte = 0x08

	FnBeginOpenSession byte = 0x41 // start then
	FnWriteTLV         byte = 0x0C // send tlv then
	FnOpenSession      byte = 0x0B // end open
)

var defaultPassword = []byte{0x1E, 0x00, 0x00, 0x00}

func newBufWithDefaultPassword(cmdID byte, isFnCmd bool) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)

	if isFnCmd { // if is fn cmd write fn cmd byte
		buf.WriteByte(FNCommand)
	}
	buf.WriteByte(cmdID)
	buf.Write(defaultPassword)

	return
}

func CreateShortStatus() []byte {
	return newBufWithDefaultPassword(ShortStatus, false).Bytes()
}

func CreateCancelCheck() []byte {
	return newBufWithDefaultPassword(CancelCheck, false).Bytes()
}

func CreateCloseSession() []byte {
	return newBufWithDefaultPassword(ZReport, false).Bytes()
}

// Operation Операции в чеке
type Operation struct {
	Type    byte   `json:"type"`    // тип операции
	Subject byte   `json:"subject"` // Предмет рассчета
	Amount  int64  `json:"amount"`  // количество товара
	Price   int64  `json:"price"`   // цена в копейках
	Sum     int64  `json:"sum"`     // сумма товар * цену
	Name    string `json:"name"`    // Наименование продукта
}

func (o Operation) Validate() error {
	if (o.Type < 1) || (o.Type > 2) {
		return errors.New("wrong operation type")
	}

	return nil
}

func CreateFNOperationV2(o Operation) (cmdData []byte, err error) {
	buf := newBufWithDefaultPassword(FnOperationV2, true)

	buf.WriteByte(o.Type)

	// Запись количества товара
	// Количество записывается в миллиграммах
	amount, err := intToBytesWithLen(o.Amount*consts.Milligram, 6)
	if err != nil {
		return nil, err
	}

	buf.Write(amount)

	// запись цены товара
	// цена записывается в копейках
	price, err := intToBytesWithLen(o.Price, 5) // одна копейка
	if err != nil {
		return nil, err
	}

	buf.Write(price)

	// запись суммы товара
	// Сумма записывается в копейках
	sum, err := intToBytesWithLen(o.Sum, 5) // две копейки
	if err != nil {
		return nil, err
	}

	buf.Write(sum)
	buf.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff}) // если нет налога надо отправлять 0xff*6
	buf.WriteByte(consts.VAT0)                      // Запись налоговой ставки
	buf.WriteByte(1)                                // // Запись номера отдела
	buf.WriteByte(consts.FullPayment)               // Запись признака способа рассчета
	buf.WriteByte(o.Subject)                        // Запись признака предмета рассчета

	// Запись название товара 0 - 128 байт строка
	// кодировка win1251
	str, err := charmap.Windows1251.NewEncoder().String(o.Name)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 128)

	if _, err := bytes.NewBufferString(str).Read(b); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	buf.Write(b)

	if buf.Len() != 160 {
		return nil, errors.New("wrong len of cmd addOperationV2")
	}

	return buf.Bytes(), nil
}

func CreateFNCloseCheck(chk models.CheckPackage) (cmdData []byte, err error) {
	// writeCashierINN in check composition before create fn close check
	buf := newBufWithDefaultPassword(FnCloseCheckV2, true)
	cash, err := intToBytesWithLen(chk.Cash, 5)
	if err != nil {
		return nil, err
	}

	casheless, err := intToBytesWithLen(chk.Casheless, 5)
	if err != nil {
		return nil, err
	}

	str, err := charmap.Windows1251.NewEncoder().String(chk.BottomLine)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 64)
	if _, err := bytes.NewBufferString(str).Read(b); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	if chk.Rounding > 99 {
		return nil, fmt.Errorf("round penni biggest than 99: %v", chk.Rounding)
	}

	buf.Write(cash)              // 5 байт сумма наличных
	buf.Write(casheless)         // 5 байт сумма безналичных
	buf.Write(make([]byte, 70))  // 5 * 14 = 70 байт остальные пустые суммы
	buf.WriteByte(chk.Rounding)  // округление до рубля в копейках, макс 99коп
	buf.Write(make([]byte, 30))  // 5 * 6 = 30 байт налогов
	buf.WriteByte(chk.TaxSystem) // биты систем налогообложения
	buf.Write(b)                 // нижняя строка чека, 64 байта win1251 текста

	log.Println("buf len:", buf.Len())

	if buf.Len() != 182 {
		return nil, fmt.Errorf("wrong FNCloseCheck len command: %v", buf.Len())
	}

	return buf.Bytes(), nil
}

func intToBytesWithLen(val int64, bytesLen int64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}

	return buf.Bytes()[:bytesLen], nil
}
