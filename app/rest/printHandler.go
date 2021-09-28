package rest

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"net/http"
	"sync"

	"github.com/go-chi/render"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
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

// Operation Операции в чеке
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
	w.Header().Set("Content-Type", "application/json")

	data := CheckReq{}

	if err := render.DecodeJSON(r.Body, &data); err != nil {
		log.Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(data)

	e := &errGroup{
		Mutex: sync.Mutex{},
		txs:   make(map[string]Tx),
	}

	var wg sync.WaitGroup

	for key, chkPkg := range data {
		wg.Add(1)
		go func(key string, chkPkg CheckPackage) {
			defer wg.Done()
			kk, ok := k.ks[key]
			if !ok {
				notFoundKKT := fmt.Errorf("не найдена касса по ключу место-организация: %v", key)
				e.addError(notFoundKKT, key)
				return
			}

			chkModelPkg, err := chkPkg.toPackageModel()
			if err != nil {
				e.addError(err, key)
				return
			}

			printCmd := kkt.PrintCheckHandler(chkModelPkg)

			e.addError(kk.Do(printCmd), key)
		}(key, chkPkg)
	}

	wg.Wait()

	log.Print(e)

	if err := json.NewEncoder(w).Encode(e.txs); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

///////////////////////////////////////////////////////

func getTaxSystemByte(tax string) (byte, error) {
	switch tax {
	case "PSN":
		return consts.PSN, nil
	case "USNIncome":
		return consts.USNIncome, nil
	default:
		return 0, fmt.Errorf("неправильная система налогообложения: %v", tax)
	}
}

func getTypeOperationByte(typ string) (byte, error) {
	switch typ {
	case "income":
		return consts.Income, nil
	case "returnIncome":
		return consts.ReturnIncome, nil
	default:
		return 0, fmt.Errorf("неправильный тип операции: %v", typ)
	}
}

func getSubByte(sub string) (byte, error) {
	switch sub {
	case "goods":
		return consts.Goods, nil
	case "service":
		return consts.Service, nil

	default:
		return 0, fmt.Errorf("неправильный признак рассчета %v", sub)
	}
}
