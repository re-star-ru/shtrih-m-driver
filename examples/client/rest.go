package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func rest(kkt *KKT) {
	r := chi.NewRouter()

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w,
			"state kkt: %v, substate: %v\ncan print: %v\n",
			kkt.state.Current(), kkt.substate.Current(), kkt.canPrintCheck()); err != nil {
			log.Println(err)
		}
	})

	r.Get("/print", func(w http.ResponseWriter, r *http.Request) {
		data := &models.CheckPackage{}
		if err := render.DecodeJSON(r.Body, data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := kkt.Do(printCheckHandler(data)); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
