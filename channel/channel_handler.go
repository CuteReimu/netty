package channel

import (
	"context"
	"github.com/CuteReimu/netty/util/concurrent"
	"sync"
)

type Handler interface {
	HandlerAdded(ctx context.Context) error
	HandlerRemoved(ctx context.Context) error
	ErrorCaught(ctx context.Context, err error) error
}

type InboundHandler interface {
	Handler
	ChannelRegistered(ctx context.Context) error
	ChannelUnregistered(ctx context.Context) error
	ChannelActive(ctx context.Context) error
	ChannelInactive(ctx context.Context) error
	ChannelRead(ctx context.Context, msg any) error
	ChannelReadComplete(ctx context.Context) error
	UserEventTriggered(ctx context.Context, evt any) error
	ChannelWritabilityChanged(ctx context.Context) error
}

type Initializer[C Channel] struct {
	InitFunc func(ch C) error
	initMap  sync.Map
}

func (i *Initializer[C]) HandlerAdded(ctx context.Context) error {
	if ctx.(interface{ Channel() Channel }).Channel().IsRegistered() {
		if ok, err := i.initChannel(ctx); err != nil {
			return err
		} else if ok {
			i.removeState(ctx)
		}
	}
	return nil
}

func (i *Initializer[C]) HandlerRemoved(ctx context.Context) error {
	i.initMap.Delete(ctx)
	return nil
}

func (i *Initializer[C]) ErrorCaught(ctx context.Context, err error) error {
	return nil
}

func (i *Initializer[C]) ChannelRegistered(ctx context.Context) error {
	if ok, err := i.initChannel(ctx); err != nil {
		return err
	} else if ok {
		i.removeState(ctx)
	}
	return nil
}

func (i *Initializer[C]) ChannelUnregistered(ctx context.Context) error {
	return nil
}

func (i *Initializer[C]) ChannelActive(ctx context.Context) error {
	return nil
}

func (i *Initializer[C]) ChannelInactive(ctx context.Context) error {
	return nil
}

func (i *Initializer[C]) ChannelRead(ctx context.Context, msg any) error {
	return nil
}

func (i *Initializer[C]) ChannelReadComplete(ctx context.Context) error {
	return nil
}

func (i *Initializer[C]) UserEventTriggered(ctx context.Context, evt any) error {
	return nil
}

func (i *Initializer[C]) ChannelWritabilityChanged(ctx context.Context) error {
	return nil
}

func (i *Initializer[C]) initChannel(ctx context.Context) (bool, error) {
	if _, exists := i.initMap.Swap(ctx, struct{}{}); !exists {
		if err := i.InitFunc(ctx.(interface{ Channel() Channel }).Channel().(C)); err != nil {
			return false, i.ErrorCaught(ctx, err)
		}
		return true, nil
	}
	return false, nil
}

func (i *Initializer[C]) removeState(ctx context.Context) {
	if ctx.(interface{ IsRemoved() bool }).IsRemoved() {
		i.initMap.Delete(ctx)
	} else {
		ctx.(interface{ Executor() concurrent.Executor }).Executor().Execute(func() {
			i.initMap.Delete(ctx)
		})
	}
}
