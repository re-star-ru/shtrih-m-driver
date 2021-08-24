package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"
	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/go-chi/render"
)

type CheckReq map[string]CheckPackage

type CheckPackage struct {
	CashierINN string      `json:"cashierINN"`
	Operations []Operation `json:"operations"` // Список операций в чеке
	Cash       int64       `json:"cash"`       // Сумма оплаты наличными
	Digital    int64       `json:"digital"`    // Сумма оплаты безналичными
	TaxSystem  string      `json:"taxSystem"`  // Система налогообложения
	Rounding   byte        `json:"rounding"`   // Округление до рубля, макс 99 копеек
	NotPrint   bool        `json:"notPrint"`   // Не печатать чек на бумаге
}

// Operation Операции в чеке
type Operation struct {
	Type    string `json:"type"`    // Тип операции
	Subject string `json:"subject"` // Предмет рассчета
	Amount  int64  `json:"amount"`  // Количество товара
	Price   int64  `json:"price"`   // Цена в копейках
	Sum     int64  `json:"sum"`     // сумма товар * цену
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

			chkModelPkg, err := packageModelFromReq(chkPkg)
			if err != nil {
				e.addError(err, key)
				return
			}

			printCmd := kkt.PrintCheckHandler(chkModelPkg)
			log.Printf("cmd print : %v\n", printCmd)

			err = kk.Do(printCmd)
			log.Println("error in print handler ", err)
			e.addError(err, key)
		}(key, chkPkg)
	}

	wg.Wait()

	log.Println(e)

	if err := json.NewEncoder(w).Encode(e.txs); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func packageModelFromReq(chk CheckPackage) (cp models.CheckPackage, err error) {
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

// ///////////////////////////////////////////////////
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
	case "outcome":
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
