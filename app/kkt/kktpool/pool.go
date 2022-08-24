package kktpool

import "github.com/re-star-ru/shtrih-m-driver/app/models"

type Req struct {
	Name string
	Err  error
	Resp chan interface{}
}

type KKTPool struct {
	KKTS map[string]chan Req
}

func NewPool(kkts []string) *KKTPool {
	workers := map[string]chan Req{}
	for _, v := range kkts {
		workers[v] = runWorker(v)
	}

	return &KKTPool{
		KKTS: workers,
	}
}

func (k *KKTPool) UpdateState() {
	for name := range k.KKTS {
		k.updateState(name)
	}
}

func (k *KKTPool) GetState(name string) {}

func (k *KKTPool) PrintCheck(name string, operation models.Operation) {

}

func (k *KKTPool) updateState(name string) interface{} {
	in := k.KKTS[name]
	req := Req{
		Name: "UpdateState",
		Resp: make(chan interface{}, 1),
	} // make request
	in <- req // send to pool

	return <-req.Resp // read reasponse
}
