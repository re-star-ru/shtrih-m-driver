package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
)

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
		return s[i].IP < s[j].IP
	})

	if _, ok := r.URL.Query()["json"]; ok {
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	if _, err := fmt.Fprintf(w, "Время: %s \n\n", time.Now().Format(time.UnixDate)); err != nil {
		log.Err(err).Send()
		return
	}

	for _, line := range s {
		if _, err := fmt.Fprintf(
			w, "Kkt ip: %v, state: %v, subState: %v \n", line.IP, line.State, line.SubState,
		); err != nil {
			log.Err(err).Send()
			return
		}
	}
}
