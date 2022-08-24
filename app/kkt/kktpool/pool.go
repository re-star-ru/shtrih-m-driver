package kktpool

import (
	"context"
	"fmt"
	"github.com/re-star-ru/shtrih-m-driver/app/configs"
	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
	"github.com/re-star-ru/shtrih-m-driver/app/models"
	"sync"
)

var ErrInvalidTypeAssertion = fmt.Errorf("invalid type assertion")
var ErrKKTNotFound = fmt.Errorf("kkt not found")

type Req struct {
	Name string
	Ctx  context.Context
	Err  chan error
	Resp chan interface{}
}

type KKTPool interface {
	GetKKTNames() []string
	GetKKT(name string) (kkt.KKT, error)

	GetStatus(ctx context.Context, name string) (models.Status, error)
	GetStatusAll(ctx context.Context) ([]models.Status, error)

	UpdateStatus(ctx context.Context, name string) error
	UpdateStatusAll(ctx context.Context) error

	PrintCheck(ctx context.Context, name string, check models.CheckPackage) error
}

func NewPool(config configs.ConfKKT) (*Pool, error) {
	pool := &Pool{
		KKTS: make(map[string]chan Req),
	}

	for key, v := range config {
		newKKT, err := kkt.NewKKT(key, v.Addr, v.Inn, configs.DefaultTimeout)
		if err != nil {
			return nil, fmt.Errorf("cant init kkt: %w", err)
		}

		pool.KKTS[key] = runWorker(newKKT)
	}

	return pool, nil
}

type Pool struct {
	KKTS map[string]chan Req
}

func (p *Pool) GetStatus(ctx context.Context, name string) (models.Status, error) {
	in := p.KKTS[name]
	req := Req{
		Name: "GetStatus",
		Ctx:  ctx,
		Err:  make(chan error, 1),
		Resp: make(chan interface{}, 1),
	}
	in <- req

	select {
	case err := <-req.Err:
		return models.Status{}, err
	case resp := <-req.Resp:
		respn, ok := resp.(models.Status)
		if !ok {
			return models.Status{}, ErrInvalidTypeAssertion
		}

		return respn, nil
	}
}
func (p *Pool) GetStatusAll(ctx context.Context) ([]models.Status, error) {
	var (
		wg sync.WaitGroup

		stats = struct {
			sync.Mutex
			stats  []models.Status
			errors []error
		}{
			stats:  []models.Status{},
			errors: []error{},
		}
	)

	wg.Add(len(p.KKTS))

	for name := range p.KKTS {
		go func(name string) {
			defer wg.Done()
			status, err := p.GetStatus(ctx, name)

			stats.Lock()
			stats.stats = append(stats.stats, status)
			stats.errors = append(stats.errors, err)
			stats.Unlock()
		}(name)
	}

	wg.Wait()

	for _, err := range stats.errors {
		if err != nil {
			return stats.stats, err
		}
	}

	return stats.stats, nil
}

func (p *Pool) UpdateStatus(ctx context.Context, name string) error {
	in := p.KKTS[name]
	req := Req{
		Name: "UpdateState",
		Ctx:  ctx,
		Err:  make(chan error, 1),
	} // make request
	in <- req // send to pool

	return <-req.Err // read reasponse
}
func (p *Pool) UpdateStatusAll(ctx context.Context) error {
	var (
		wg sync.WaitGroup

		updates = struct {
			sync.Mutex
			errors []error
		}{
			errors: []error{},
		}
	)

	wg.Add(len(p.KKTS))

	for name := range p.KKTS {
		go func(name string) {
			defer wg.Done()

			err := p.UpdateStatus(ctx, name)

			updates.Lock()
			updates.errors = append(updates.errors, err)
			updates.Unlock()
		}(name)
	}

	wg.Wait()

	for _, err := range updates.errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pool) GetKKTNames() []string {
	kktNames := make([]string, 0)
	for k := range p.KKTS {
		kktNames = append(kktNames, k)
	}
	return kktNames
}
func (p *Pool) GetKKT(name string) (chan Req, error) {
	if _, ok := p.KKTS[name]; !ok {
		return nil, ErrKKTNotFound
	}

	return p.KKTS[name], nil
}

func (p *Pool) PrintCheck(ctx context.Context, name string, operation models.CheckPackage) error {
	return nil
}
