package kktpool

import (
	"context"
	"fmt"
	"github.com/re-star-ru/shtrih-m-driver/app/configs"
	"github.com/re-star-ru/shtrih-m-driver/app/kkt"
	"github.com/re-star-ru/shtrih-m-driver/app/models"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var ErrInvalidTypeAssertion = fmt.Errorf("invalid type assertion")
var ErrKKTNotFound = fmt.Errorf("kkt not found")

type Req struct {
	Name string
	Ctx  context.Context
	Body interface{}
	Resp chan Resp
}

type Resp struct {
	Err  error
	Body interface{}
}

type KKTPool interface {
	GetKKTNames() []string

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

	ctx := context.Background()
	pool.startUpdater(ctx) // state updater & healtchecker

	return pool, nil
}

type Pool struct {
	KKTS map[string]chan Req
}

func (p *Pool) GetStatus(ctx context.Context, name string) (models.Status, error) {
	in, ok := p.KKTS[name]
	if !ok {
		return models.Status{}, ErrKKTNotFound
	}

	req := Req{
		Name: "GetStatus",
		Ctx:  ctx,
		Resp: make(chan Resp, 1),
	}

	log.Debug().Msg("send request to kkt")
	in <- req
	log.Debug().Msg("request sent")

	resp := <-req.Resp
	if resp.Err != nil {
		return models.Status{}, resp.Err
	}

	respn, ok := resp.Body.(models.Status)
	if !ok {
		return models.Status{}, ErrInvalidTypeAssertion
	}

	return respn, nil
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

	for name := range p.KKTS {
		wg.Add(1)

		go func(name string) {
			defer wg.Done()

			log.Debug().Msgf("get status for %s", name)
			status, err := p.GetStatus(ctx, name)
			log.Debug().Msgf("after got status for %s", name)

			stats.Lock()
			stats.stats = append(stats.stats, status)
			stats.errors = append(stats.errors, err)
			stats.Unlock()
		}(name)
	}

	log.Debug().Msg("waiting for kkts statuses")
	wg.Wait()
	log.Debug().Msg("kkts statuses received")

	for _, err := range stats.errors {
		if err != nil {
			return stats.stats, err
		}
	}

	return stats.stats, nil
}

const updateTimeout = time.Second * 30

// run healthcheck
func (p *Pool) startUpdater(ctx context.Context) {
	for _, name := range p.GetKKTNames() {
		go func(name string) {
			for {
				err := p.UpdateStatus(ctx, name)
				if err != nil {
					log.Error().Err(err).Msg("cant update status")
				}

				time.Sleep(updateTimeout)
			}
		}(name)
	}
}

func (p *Pool) UpdateStatus(ctx context.Context, name string) error {
	in := p.KKTS[name]
	req := Req{
		Name: "UpdateState",
		Ctx:  ctx,
		Resp: make(chan Resp, 1),
	} // make request
	in <- req // send to pool
	resp := <-req.Resp

	return resp.Err // read reasponse
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

func (p *Pool) PrintCheck(ctx context.Context, name string, check models.CheckPackage) error {
	in, ok := p.KKTS[name]
	if !ok {
		return fmt.Errorf("%w: %v", ErrKKTNotFound, name)
	}

	req := Req{
		Name: "PrintCheck",
		Ctx:  ctx,
		Body: check,
		Resp: make(chan Resp, 1),
	}
	in <- req
	resp := <-req.Resp

	return resp.Err
}
