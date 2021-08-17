package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	CashierINN string      `json:"cashierINN"`
	Operations []Operation `json:"operations"` // Список операций в чеке
	Cash       int64       `json:"cash"`       // Сумма оплаты наличными
	Digital    int64       `json:"digital"`    // Сумма оплаты безналичными
	Rounding   byte        `json:"rounding"`   // Округление до рубля, макс 99 копеек
	TaxSystem  byte        `json:"taxSystem"`  // Система налогообложения
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

func (k *KKTService) printPackageHandler(w http.ResponseWriter, r *http.Request) {
	data := CheckReq{}

	if err := render.DecodeJSON(r.Body, &data); err != nil {
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

	//if err := kkt.Do(printCheckHandler(data)); err != nil {
	//	log.Println(err)
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//
	//	return
	//}
}

/*
0646420041024536
заводской номер уникальный


*/
