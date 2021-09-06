package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/fess932/shtrih-m-driver/examples/client/kkt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type KKTService struct {
	ks   map[string]*kkt.KKT
	addr string
}

func New(ks map[string]*kkt.KKT, addr string) *KKTService {
	return &KKTService{
		ks:   ks,
		addr: addr,
	}
}

func (k *KKTService) Run() {
	k.rest()
}

func (k *KKTService) rest() {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(time.Second * 30))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	r.Get("/status", k.status)
	r.Post("/printPackage", func(w http.ResponseWriter, r *http.Request) {
		k.printPackageHandler(w, r)
	})

	log.Println("server listen at: ", k.addr)
	log.Fatal(http.ListenAndServe(k.addr, r))
}

type Status struct {
	IP       string `json:"ip"`
	State    string `json:"state"`
	SubState string `json:"subState"`
}

func (k *KKTService) status(w http.ResponseWriter, r *http.Request) {
	s := make([]Status, 0, len(k.ks))

	for _, kk := range k.ks {
		s = append(s, Status{IP: kk.Addr, State: kk.State.Current(), SubState: kk.Substate.Current()})
	}

	sort.Slice(s, func(i, j int) bool {
		return false
	})

	if _, ok := r.URL.Query()["json"]; ok {
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	if _, err := fmt.Fprintf(w, "Время: %s \n\n", time.Now().Format(time.UnixDate)); err != nil {
		log.Println(err)
		return
	}

	for _, line := range s {
		if _, err := fmt.Fprintf(w, "Kkt ip: %v, state: %v, subState: %v \n", line.IP, line.State, line.SubState); err != nil {
			log.Println(err)
			return
		}
	}
}
