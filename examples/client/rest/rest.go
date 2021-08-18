package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fess932/shtrih-m-driver/pkg/consts"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type KKTService struct {
	ks   []*kkt.KKT
	addr string
}

func New(ks ...*kkt.KKT) *KKTService {
	return &KKTService{
		ks:   ks,
		addr: "",
	}
}

func (k *KKTService) Run() {
	k.rest()
}

func (k *KKTService) rest() {
	r := chi.NewRouter()

	r.Get("/status", k.status)
	r.Post("/printPackage", func(w http.ResponseWriter, r *http.Request) {
		k.printPackageHandler(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

type Status struct {
	IP    string `json:"ip"`
	State string `json:"state"`
}

func (k *KKTService) status(w http.ResponseWriter, r *http.Request) {
	s := make([]Status, 0, len(k.ks))

	// todo run concurrent with gorutines
	for _, kk := range k.ks {
		s = append(s, Status{IP: kk.Addr, State: kk.State.Current()})
	}

	if _, ok := r.URL.Query()["json"]; ok {
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	for _, line := range s {
		if _, err := fmt.Fprintf(w, "Kkt ip: %v, state: %v \n", line.IP, line.State); err != nil {
			log.Println(err)
			return
		}
	}
}

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
	Type    byte   `json:"type"`    // Тип операции
	Subject byte   `json:"subject"` // Предмет рассчета
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

	for organization, chkPkg := range data {
		kk := k.getPrinterByOrgAndPlace(organization, chkPkg.Place)
		if kk == nil {
			notFoundKKT := fmt.Errorf("не найдена касса для организации: %v, и места: %v", organization, chkPkg.Place)
			http.Error(w, notFoundKKT.Error(), http.StatusNotFound)
			return
		}

		chkModelPkg, err := packageModelFromReq(chkPkg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		printCmd := kkt.PrintCheckHandler(chkModelPkg)
		log.Printf("print cmd: %v", printCmd)
		//if err := kk.Do(printCmd); err != nil {
		//	log.Println(err)
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//	return
		//}
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

/*
0646420041024536
заводской номер уникальный


*/
