package channel

import "net"

const (
	defaultLowWaterMark  = 32 * 1024
	defaultHighWaterMark = 64 * 1024
)

type Option struct {
	ConnectTimeoutMillis int32
	MaxMessagesPerWrite  int32

	WriteSpinCount int32

	AllowHalfClosure bool
	AutoRead         bool

	// If true then the Channel is closed automatically and immediately on write failure
	AutoClose bool

	SoBroadcast bool
	SoKeepAlive bool
	SoSndbuf    int32
	SoRcvbuf    int32
	SoReuseaddr bool
	SoLinger    int32
	SoBacklog   int32
	SoTimeout   int32

	IpTos                   int32
	IpMulticastAddr         net.Addr
	IpMulticastIf           *net.Interface
	IpMulticastTtl          int32
	IpMulticastLoopDisabled bool

	TcpNoDelay bool
	// Client-side TCP FastOpen. Sending data with the initial TCP handshake.
	TcpFastopenConnect bool

	// Server-side TCP FastOpen. Configures the maximum number of outstanding (waiting to be accepted) TFO connections.
	TcpFastopen int32

	SingleEventexecutorPreGroup bool
}
