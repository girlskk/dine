package grace

import "context"

type Process interface {
	Run() error
	Shutdown()
}

func RunWaitContext(ctx context.Context, process Process) error {
	errCh := make(chan error)
	go func() {
		errCh <- process.Run()
	}()

	select {
	case <-ctx.Done():
		process.Shutdown()
		// wait for run to return
		return <-errCh
	case err := <-errCh:
		return err
	}
}

type Stoper struct {
	handle      func(ctx context.Context)
	stoppedChan chan struct{}
	cancel      context.CancelFunc
}

func NewStoper(handle func(ctx context.Context)) *Stoper {
	return &Stoper{
		handle:      handle,
		stoppedChan: make(chan struct{}),
	}
}

func (s *Stoper) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go func() {
		defer close(s.stoppedChan)
		s.handle(ctx)
	}()
}

func (s *Stoper) Stop() {
	s.cancel()
	<-s.stoppedChan
}

type BlockProcess struct {
	stopChan chan struct{}
}

func NewBlockProcess() BlockProcess {
	return BlockProcess{
		stopChan: make(chan struct{}),
	}
}

// Run .
func (bp BlockProcess) Run() error {
	<-bp.stopChan
	return nil
}

// Shutdown .
func (bp BlockProcess) Shutdown() {
	close(bp.stopChan)
}
