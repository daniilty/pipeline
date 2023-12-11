## Pipelines

### Goal

This package provides you with a simple way to process business pipelines in async way.
The idea is you don't have to wait for some function to finish, instead you wait for data this function provides, everything else can be async.

### Installation

```bash
$ go get github.com/daniilty/pipeline@latest
```

### Usage 

```go
// set timeout to 1 second(after timeout all waiting runners will get timeout error)
// set to 0 to disable
pipe := pipeline.NewPipeline(1 * time.Second)
pipe.Execute(context.Background(),
    func(_ context.Context, provider pipeline.Provider) {
		provider.ProvideData("d1", "somedata", nil)
	},
    func(_ context.Context, provider pipeline.Provider) {
		wd1, ok := provider.WaitData("d1")
		if !ok {
			log.Fatal("did not get d1")
		}

		if wd1.Err() != nil {
			log.Fatal("d1 fail")
		}

		provider.ProvideData("d2", nil, errors.New("shit happened"))
	},
    func(_ context.Context, provider pipeline.Provider) {
		wd2, ok := provider.WaitData("d2")
		if !ok {
			log.Fatal("did not get d2")
		}

		if wd2.Err() != nil {
			log.Fatalf("d2 fail: %s", wd.Err().Error())
		}
	},
)
```
