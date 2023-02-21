package channel

import "github.com/CuteReimu/netty/util/concurrent"

type Future interface {
	concurrent.Future
	Id() Id
}
