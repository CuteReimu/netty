package bootstrap

import (
	"github.com/CuteReimu/netty/channel"
	"net"
)

type ServerBootstrap struct {
	ChildGroup   channel.EventLoopGroup
	ChildHandler channel.Handler

	ParentGroup    channel.EventLoopGroup
	newChannelFunc func() channel.Channel
	localAddress   net.Addr
	ParentHandler  channel.Handler
}

func (bs *ServerBootstrap) init(ch channel.Channel) {
	p := ch.Pipeline()

	currentChildGroup, currentChildHandler := bs.ChildGroup, bs.ChildHandler

	p.PushBack(channel.Initializer[channel.Channel]{
		InitFunc: func(ch channel.Channel) error {
			pipeline := ch.Pipeline()
			ch.EventLoop().Execute(func() {
				pipeline.PushBack(serverBootstrapAcceptor{
					childGroup:     currentChildGroup,
					channelHandler: currentChildHandler,
				})
			})
			return nil
		},
	})
}

type serverBootstrapAcceptor struct {
	childGroup     channel.EventLoopGroup
	channelHandler channel.Handler
}

// TODO: need to implement serverBootstrapAcceptor
