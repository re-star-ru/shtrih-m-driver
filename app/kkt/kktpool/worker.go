package kktpool

import (
	"errors"
	"fmt"
	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
	"github.com/re-star-ru/shtrih-m-driver/app/models"
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
		for req := range in {
			req.Resp <- w.handle(req)
		}
	}(in)

	return in
}

func (w *worker) handle(req Req) Resp {
	switch req.Name {
	case "GetStatus":
		status := models.Status{
			IP:       w.kkt.Addr,
			State:    w.kkt.State.Current(),
			SubState: w.kkt.Substate.Current(),
		}

		return Resp{
			Err:  nil,
			Body: status,
		}

	case "UpdateState":
		if err := w.kkt.Do(kkt.UpdateState); err != nil {
			return Resp{
				fmt.Errorf("req name: %s: %w", req.Name, err),
				nil,
			}
		}

		return Resp{}

	case "PrintCheck":
		check, ok := req.Body.(models.CheckPackage)
		if !ok {
			return Resp{Err: fmt.Errorf("PrintCheck %w", ErrInvalidTypeAssertion)}
		}

		if err := w.kkt.Do(kkt.PrintCheckHandler(check)); err != nil {
			if err != nil {
				return Resp{Err: fmt.Errorf("PrintCheck %w", err)}
			}
		}

		return Resp{}

	default:
		return Resp{Err: ErrCommandNotFound}
	}
}
