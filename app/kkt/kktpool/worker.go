package kktpool

import (
	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
)

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
			req.Resp <- "hello" // todo cmd routing
		}
	}(in)

	return in
}
