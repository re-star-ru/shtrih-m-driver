package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"

	"github.com/go-chi/chi/v5"
)

type KKTService struct {
	ks   map[string]*kkt.KKT
	addr string
}

func New(ks map[string]*kkt.KKT) *KKTService {
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
