package pipeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daniilty/pipeline"
)

func TestPipeline(t *testing.T) {
	pipe := pipeline.NewPipeline(1 * time.Second)
	d1 := "shit"
	d2 := "fuck"

	f1 := func(_ context.Context, provider pipeline.Provider) {
		time.Sleep(100 * time.Millisecond)
		provider.ProvideData("d1", d1, nil)
	}

	f2 := func(_ context.Context, provider pipeline.Provider) {
		wd1, ok := provider.WaitData("d1")
		if !ok {
			t.Error("did not get d1")
			return
		}

		if wd1.Err() != nil {
			t.Errorf("d1 err: %s", wd1.Err().Error())
			return
		}

		time.Sleep(100 * time.Millisecond)
		provider.ProvideData("d2", d2, nil)
	}

	f3 := func(_ context.Context, provider pipeline.Provider) {
		wd1, ok := provider.WaitData("d1")
		if !ok {
			t.Fatal("did not get d1")
		}

		if wd1.Err() != nil {
			t.Fatal("d1 fail")
		}

		wd2, ok := provider.WaitData("d2")
		if !ok {
			t.Fatal("did not get d1")
		}

		if wd2.Err() != nil {
			t.Fatal("d2 fail")
		}

		time.Sleep(100 * time.Millisecond)
		provider.ProvideData("d3", nil, errors.New("sample"))
	}

	f4 := func(_ context.Context, provider pipeline.Provider) {
		wd3, ok := provider.WaitData("d3")
		if !ok {
			t.Fatal("did not get d3")
		}

		if wd3.Err() == nil {
			t.Fatal("d3 did not fail")
		}
	}

	f5 := func(_ context.Context, provider pipeline.Provider) {
		time.Sleep(2 * time.Second)

		wd1, ok := provider.WaitData("d1")
		if !ok {
			t.Fatal("did not get d1")
		}

		if wd1.Err() != nil {
			t.Fatal("d1 fail")
		}
	}

	pipe.Execute(context.Background(), f1, f2, f3, f4, f5)
}
