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

	if _, err = fmt.Fprintf(w, "Время: %s \n\n", time.Now().Format(time.UnixDate)); err != nil {
		log.Err(err).Send()
		return
	}

	for _, line := range s {
		if _, err = fmt.Fprintf(
			w, "Адрес кассы ip: %v, статус: %v, подстатус: %v \n", line.IP, line.State, line.SubState,
		); err != nil {
			log.Err(err).Send()
			return
		}
	}
}
