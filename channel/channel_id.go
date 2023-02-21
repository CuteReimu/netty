package channel

import (
	"encoding/binary"
	"fmt"
	"github.com/CuteReimu/netty/util"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type Id interface {
	// AsShortText returns the short but globally non-unique string representation of the ChannelId.
	AsShortText() string

	// AsLongText returns the long yet globally unique string representation of the ChannelId.
	AsLongText() string

	HashCode() uint32
}

var (
	machineId    = util.BestAvailableMac()
	processId    = uint32(os.Getpid())
	nextSequence uint32
)

const (
	processIdLen = 4
	sequenceLen  = 4
	timestampLen = 8
	randomLen    = 4
)

type defaultChannelId struct {
	data       []byte
	hashCode   uint32
	shortValue string
	longValue  string
}

func NewDefaultChannelId() Id {
	id := &defaultChannelId{
		data: make([]byte, len(machineId)+processIdLen+sequenceLen+timestampLen+randomLen),
	}
	var i int

	// machineId
	copy(id.data[i:], machineId)
	i += len(machineId)

	// processId
	i = id.writeUint32(i, processId)

	// sequence
	i = id.writeUint32(i, atomic.AddUint32(&nextSequence, 1))

	// timestamp (kind of)
	now := time.Now()
	i = id.writeUint64(i, util.ReverseUint64(uint64(now.UnixNano()))^uint64(now.UnixMilli()))

	// random
	i = id.writeUint32(i, rand.Uint32())
	if i != len(id.data) {
		panic(fmt.Sprint("incorrect len: ", i, " | ", len(id.data)))
	}

	id.hashCode = util.HashCode(id.data)

	return id
}

func (id *defaultChannelId) writeUint32(i int, value uint32) int {
	binary.BigEndian.PutUint32(id.data[i:], value)
	return i + 4
}

func (id *defaultChannelId) writeUint64(i int, value uint64) int {
	binary.BigEndian.PutUint64(id.data[i:], value)
	return i + 8
}

func (id *defaultChannelId) AsShortText() string {
	if len(id.shortValue) == 0 {
		id.shortValue = util.HexDump(id.data, len(id.data)-randomLen, randomLen)
	}
	return id.shortValue
}

func (id *defaultChannelId) AsLongText() string {
	if len(id.longValue) == 0 {
		id.longValue = id.newLongValue()
	}
	return id.longValue
}

func (id *defaultChannelId) newLongValue() string {
	buf := &strings.Builder{}
	buf.Grow(2*len(id.data) + 5)
	var i int
	i = id.appendHexDumpField(buf, i, len(machineId))
	i = id.appendHexDumpField(buf, i, processIdLen)
	i = id.appendHexDumpField(buf, i, sequenceLen)
	i = id.appendHexDumpField(buf, i, timestampLen)
	i = id.appendHexDumpField(buf, i, randomLen)
	if i != len(id.data) {
		panic(fmt.Sprint("incorrect len: ", i, " | ", len(id.data)))
	}
	s := buf.String()
	return s[:len(s)-1]
}

func (id *defaultChannelId) appendHexDumpField(buf *strings.Builder, i, length int) int {
	buf.WriteString(util.HexDump(id.data, i, length))
	buf.WriteByte('-')
	return i + length
}

func (id *defaultChannelId) HashCode() uint32 {
	return id.hashCode
}
