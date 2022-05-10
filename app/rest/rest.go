package rest

import (
	"github.com/re-star-ru/shtrih-m-driver/app/kkt/kktpool"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type KKTService struct {
	pool kktpool.KKTPool
	addr string
}

func New(pool kktpool.KKTPool, addr string) *KKTService {
	return &KKTService{
		pool: pool,
		addr: addr,
	}
}

func (k *KKTService) Run() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	// prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	{
		r.Get("/status", k.status)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Post("/printPackage", k.printPackageHandler)
	}

	log.Info().Msgf("server listen at: %v", k.addr)
	log.Fatal().Err(http.ListenAndServe(k.addr, r)).Send()
}
