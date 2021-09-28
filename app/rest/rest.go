package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
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
	r.Use(middleware.Timeout(time.Second * 120))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	{
		r.Get("/status", k.status)
		r.Post("/printPackage", k.printPackageHandler)
	}

	log.Print("server listen at: ", k.addr)
	log.Fatal().Err(http.ListenAndServe(k.addr, r)).Send()
}
