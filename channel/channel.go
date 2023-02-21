package channel

import (
	"container/list"
	"github.com/CuteReimu/netty/util/concurrent"
	"net"
)

type EventLoopGroup interface {
	concurrent.EventExecutorGroup
	Register(channel Channel) Future
	// TODO: Next() EventLoop
}

type EventLoop interface {
	EventLoopGroup
	Parent() EventLoopGroup
}

type Metadata struct {
	hasDisconnect             bool
	defaultMaxMessagesPerRead int32
}

func (m *Metadata) HasDisconnect() bool {
	return m.hasDisconnect
}

func (m *Metadata) DefaultMaxMessagesPerRead() int32 {
	return m.defaultMaxMessagesPerRead
}

type Channel interface {
	// Id returns the globally unique identifier of this Channel.
	Id() Id

	// EventLoop returns the EventLoop this Channel was registered to.
	EventLoop() EventLoop

	// Parent returns the parent of this channel, nil if this channel does not have a parent channel.
	Parent() Channel

	// IsOpen returns true if the Channel is open and may get active later
	IsOpen() bool

	// IsRegistered returns true if the Channel is registered.
	IsRegistered() bool

	// IsActive returns true if the Channel is active and so connected.
	IsActive() bool

	// Metadata returns the ChannelMetadata of the Channel which describe the nature of the Channel.
	Metadata() Metadata

	LocalAddress() net.Addr
	RemoteAddress() net.Addr

	// CloseFuture returns the Future which will be notified when this channel is closed. This method always returns the same future instance.
	CloseFuture() Future

	IsWriteable() bool
	BytesBeforeUnwritable() int64
	BytesBeforeWritable() int64

	Pipeline() *list.List

	Read() Channel
	Flush() Channel
}
