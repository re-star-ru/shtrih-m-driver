package kktpool

import (
	"errors"
	"fmt"
	"github.com/re-star-ru/shtrih-m-driver/app/models"

	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
)

var ErrCommandNotFound = errors.New("command not found")

type worker struct {
	name string
	kkt  *kkt.KKT
}

func runWorker(k *kkt.KKT) chan Req {
	w := &worker{
		name: k.Addr,
		kkt:  k,
	}

	return w.run()
}

func (w *worker) run() chan Req {
	in := make(chan Req)
	go func(in chan Req) {
		var (
			resp interface{}
			err  error
		)
		for req := range in {
			resp, err = w.handle(req) // todo cmd routing
			req.Resp <- resp
			req.Err <- err
		}
	}(in)

	return in
}

func (w *worker) handle(req Req) (interface{}, error) {
	switch req.Name {
	case "GetStatus":
		status := models.Status{
			IP:       w.kkt.Addr,
			State:    w.kkt.State.Current(),
			SubState: w.kkt.Substate.Current(),
		}

		return status, nil

	case "UpdateState":
		if err := w.kkt.Do(kkt.UpdateState); err != nil {
			return nil, fmt.Errorf("req name: %s: %w", req.Name, err)
		}

		return nil, nil

	case "PrintCheck":
		return nil, nil

	default:
		return nil, ErrCommandNotFound
	}
}
