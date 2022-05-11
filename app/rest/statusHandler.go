package rest

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
)

func (k *KKTService) status(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("get status all")

	s, err := k.pool.GetStatusAll(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("get status")
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].IP < s[j].IP
	})

	if _, ok := r.URL.Query()["json"]; ok {
		render.JSON(w, r, s)

		return
	}

	render.PlainText(w, r, fmt.Sprintf("Время: %s \n\n", time.Now().Format(time.UnixDate)))

	for _, line := range s {
		render.PlainText(
			w, r,
			fmt.Sprintf("Адрес кассы ip: %v, статус: %v, подстатус: %v\n", line.IP, line.State, line.SubState),
		)
	}
}
