package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
}

type ctxCancelFunc struct{}

func NewApp() App {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, ctxCancelFunc{}, cancel)
	return App{
		ctx:       ctx,
		cancel:    cancel,
		waitGroup: new(sync.WaitGroup),
	}
}

func (app *App) Run(task func(context.Context)) {
	app.waitGroup.Add(1)
	go func() {
		task(app.ctx)
		app.waitGroup.Done()
	}()
}

func (app *App) Call(task func(context.Context)) {
	task(app.ctx)
}

func (app *App) CallAsync(task func(context.Context)) {
	go task(app.ctx)
}

func (app *App) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-app.ctx.Done():
	case <-quit:
	}

	app.cancel()
	app.waitGroup.Wait()
}

// 如果 ctx 是由 app.Run 传入的，此函数会取消 ctx
func StopApp(ctx context.Context) {
	cancel, ok := ctx.Value(ctxCancelFunc{}).(context.CancelFunc)
	if !ok {
		panic("this context is not created by app")
	}
	cancel()
}
