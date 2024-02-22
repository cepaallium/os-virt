package exit

import "context"

func SetupCtxWithStop(ctx context.Context, stop <-chan struct{}) context.Context {
	_ctx, cancel := context.WithCancel(ctx)

	go func() {
		// cancel sub goroutine and wait exit
		<-stop
		cancel()
	}()

	return _ctx
}
