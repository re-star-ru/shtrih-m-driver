package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"
	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/go-chi/render"
)

type CheckReq map[string]CheckPackage

type CheckPackage struct {
	Place      string      `json:"place"`
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

func (k *KKTService) getPrinterByOrgAndPlace(organization, place string) *kkt.KKT {
	for _, kk := range k.ks {
		if kk.Place == place && kk.Organization == organization {
			return kk
		}
	}

	return nil
}

func (k *KKTService) printPackageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := CheckReq{}

	if err := render.DecodeJSON(r.Body, &data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// concurrent run print cmd
	var g errgroup.Group

	for key, chkPkg := range data {
		ky := key
		c := chkPkg
		g.Go(func() error {
			kk, ok := k.ks[ky]
			if !ok {
				notFoundKKT := fmt.Errorf("не найдена касса по ключу место-организация: %v", ky)
				return notFoundKKT
			}

			chkModelPkg, err := packageModelFromReq(c)
			if err != nil {
				return err
			}

			printCmd := kkt.PrintCheckHandler(chkModelPkg)
			log.Printf("cmd print : %v\n", printCmd)

			err = kk.Do(printCmd)
			return err
		})
	}

	err := g.Wait()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(data)

	if err := json.NewEncoder(w).Encode(data); err != nil {
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
