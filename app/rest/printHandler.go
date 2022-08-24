package rest

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/render"

	"github.com/re-star-ru/shtrih-m-driver/app/models"
	"github.com/re-star-ru/shtrih-m-driver/app/models/consts"
)

type CheckReq map[string]CheckPackage

type CheckPackage struct {
	CashierINN string      `json:"cashierINN"`
	Operations []Operation `json:"operations"` // Список операций в чеке
	Cash       uint64      `json:"cash"`       // Сумма оплаты наличными
	Digital    uint64      `json:"digital"`    // Сумма оплаты безналичными
	TaxSystem  string      `json:"taxSystem"`  // Система налогообложения
	Rounding   byte        `json:"rounding"`   // Округление до рубля, макс 99 копеек
	NotPrint   bool        `json:"notPrint"`   // Не печатать чек на бумаге
}

func (chk *CheckPackage) toPackageModel() (cp models.CheckPackage, err error) {
	cp = models.CheckPackage{
		CashierINN: chk.CashierINN,
		Cash:       chk.Cash,
		Digital:    chk.Digital,
		Rounding:   chk.Rounding,
		NotPrint:   chk.NotPrint,
	}

	cp.TaxSystem, err = getTaxSystemByte(chk.TaxSystem)
	if err != nil {
		return
	}

	for _, v := range chk.Operations {
		typ, err := getTypeOperationByte(v.Type)
		if err != nil {
			return models.CheckPackage{}, err
		}

		sub, err := getSubByte(v.Subject)
		if err != nil {
			return models.CheckPackage{}, err
		}

		op := models.Operation{
			Type:    typ,
			Subject: sub,
			Amount:  v.Amount,
			Price:   v.Price,
			Sum:     v.Sum,
			Name:    v.Name,
		}

		cp.Operations = append(cp.Operations, op)
	}

	return cp, nil
}

// Operation Операции в чеке.
type Operation struct {
	Type    string `json:"type"`    // Тип операции
	Subject string `json:"subject"` // Предмет рассчета
	Amount  uint64 `json:"amount"`  // Количество товара
	Price   uint64 `json:"price"`   // Цена в копейках
	Sum     uint64 `json:"sum"`     // сумма товар * цену
	Name    string `json:"name"`    // Наименование продукта
}

type Tx struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type errGroup struct {
	sync.Mutex
	txs map[string]Tx
}

func (e *errGroup) addError(err error, key string) {
	e.Lock()
	defer e.Unlock()

	if err != nil {
		e.txs[key] = Tx{
			Status: "error",
			Error:  err.Error(),
		}

		return
	}

	e.txs[key] = Tx{
		Status: "done",
		Error:  "",
	}
}

func (k *KKTService) printPackageHandler(w http.ResponseWriter, r *http.Request) {
	data := CheckReq{}
	if err := render.DecodeJSON(r.Body, &data); err != nil {
		log.Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	log.Debug().Msgf("input data %+v", data)

	e := &errGroup{
		Mutex: sync.Mutex{},
		txs:   make(map[string]Tx),
	}

	var wg sync.WaitGroup

	wg.Add(len(data))

	for name, check := range data {
		go func(name string, check CheckPackage) {
			defer wg.Done()

			cp, err := check.toPackageModel()
			if err != nil {
				e.addError(fmt.Errorf("cant convert check: %w", err), name)

				return
			}

			if err = k.pool.PrintCheck(r.Context(), name, cp); err != nil {
				e.addError(fmt.Errorf("cant print check: %w", err), name)

				return
			}
		}(name, check)
	}

	wg.Wait()

	log.Debug().Msgf("result %+v", e)
	render.JSON(w, r, e.txs)
}

///////////////////////////////////////////////////////

var (
	ErrWrongTaxSystem     = errors.New("неправильная система налогообложения")
	ErrWrongOperationType = errors.New("неправильный тип операции")
	ErrWrongSubject       = errors.New("неправильный предмет расчета")
)

func getTaxSystemByte(tax string) (byte, error) {
	switch tax {
	case "PSN":
		return consts.PSN, nil
	case "USNIncome":
		return consts.USNIncome, nil
	default:
		return 0, fmt.Errorf("ошибка %w: %v", ErrWrongTaxSystem, tax)
	}
}

func getTypeOperationByte(typ string) (byte, error) {
	switch typ {
	case "income":
		return consts.Income, nil
	case "returnIncome":
		return consts.ReturnIncome, nil
	default:
		return 0, fmt.Errorf("ошибка %w: %v", ErrWrongOperationType, typ)
	}
}

func getSubByte(sub string) (byte, error) {
	switch sub {
	case "goods":
		return consts.Goods, nil
	case "service":
		return consts.Service, nil

	default:
		return 0, fmt.Errorf("ошибка %w: %v", ErrWrongSubject, sub)
	}
}
