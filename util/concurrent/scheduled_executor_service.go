package concurrent

type Future interface {
	AddListeners(listener ...func() error) Future
	RemoveListeners(listener ...func() error) Future
	Sync() Future
	Await() Future
	IsSuccess() bool
}

type Executor interface {
	Execute(command func())
}

type ExecutorService interface {
	Executor
	Shutdown()
	ShutdownNow() []func()
	IsShutdown() bool
	IsTerminated() bool
	AwaitTerminated() bool
	Submit(command func())
}

type EventExecutorGroup interface {
	ExecutorService
	IsShuttingDown() bool
	ShutdownGracefully() Future
	Foreach(func(EventExecutor) bool)
	Next() EventExecutor
}

type EventExecutor interface {
	EventExecutorGroup
	Parent() EventExecutorGroup
}
