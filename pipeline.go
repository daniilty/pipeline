package pipeline

import (
	"context"
	"sync"
	"time"
)

type Runner func(context.Context, Provider)

type Provider interface {
	ProvideData(string, any, error)
	WaitData(string) (Result, bool)
}

type Result interface {
	Err() error
	Data() any
}

type Pipeline struct {
	data *kv
}

func NewPipeline(timeout time.Duration) *Pipeline {
	return &Pipeline{
		data: &kv{
			timeout: timeout,
			mux:     &sync.Mutex{},
		},
	}
}

func (p *Pipeline) Execute(ctx context.Context, runners ...Runner) {
	wg := &sync.WaitGroup{}
	for _, r := range runners {
		wg.Add(1)
		go func(r Runner) {
			r(ctx, p.data)
			wg.Done()
		}(r)
	}
	wg.Wait()
}
