package pipeline

import (
	"sync"
	"time"
)

type result struct {
	err error
	v   any
}

type data struct {
	cond  *sync.Cond
	mux   *sync.Mutex
	timer *time.Timer
	k     string
	res   Result
}

type kv struct {
	timeout time.Duration
	mux     *sync.Mutex
	entries []*data
}

func (r *result) Err() error {
	return r.err
}

func (r *result) Data() any {
	return r.v
}

func (k *kv) WaitData(key string) (Result, bool) {
	k.mux.Lock()
	for _, e := range k.entries {
		if e.k == key {
			v := e.res
			k.mux.Unlock()

			if v == nil {
				e.cond.L.Lock()
				e.cond.Wait()
				e.cond.L.Unlock()
				v = e.res
			}

			return v, v != nil
		}
	}

	d := &data{
		cond: sync.NewCond(&sync.Mutex{}),
		mux:  &sync.Mutex{},
		k:    key,
	}

	if k.timeout > 0 {
		d.timer = time.AfterFunc(k.timeout, func() {
			d.mux.Lock()
			d.res = &result{
				err: ErrTimeout,
			}
			d.mux.Unlock()

			d.cond.L.Lock()
			d.cond.Broadcast()
			d.cond.L.Unlock()
		})
	}

	k.entries = append(k.entries, d)
	k.mux.Unlock()

	d.cond.L.Lock()
	d.cond.Wait()
	d.cond.L.Unlock()

	d.mux.Lock()
	v := d.res
	d.mux.Unlock()

	return v, v != nil
}

func (k *kv) ProvideData(key string, v any, err error) {
	k.mux.Lock()
	defer k.mux.Unlock()

	for _, e := range k.entries {
		if e.k == key {
			emptyRes := e.res == nil
			if !emptyRes && e.res.Err() != nil {
				return
			}

			e.res = &result{
				err: err,
				v:   v,
			}

			if e.cond != nil {
				if emptyRes {
					e.cond.L.Lock()
					e.cond.Broadcast()
					e.cond.L.Unlock()

					if e.timer != nil {
						e.timer.Stop()
					}
				}
			}

			return
		}
	}

	k.entries = append(k.entries, &data{
		k: key,
		res: &result{
			err: err,
			v:   v,
		},
		mux: &sync.Mutex{},
	})
}
