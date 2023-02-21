package util

import "fmt"

var hexdumpTable = make([]byte, 256*4)

func init() {
	digits := []byte("0123456789abcdef")
	for i := 0; i < 256; i++ {
		hexdumpTable[i<<1] = digits[i>>4&0x0F]
		hexdumpTable[(i<<1)+1] = digits[i&0x0F]
	}
}

// HexDump returns a hex dump  of the specified byte array's sub-region.
func HexDump(array []byte, fromIndex, length int) string {
	if length < 0 {
		panic(fmt.Sprint("length : ", length, " (expected: >= 0)"))
	}
	if length == 0 {
		return ""
	}

	endIndex := fromIndex + length
	buf := make([]byte, length<<1)

	srcIdx, dstIdx := fromIndex, 0
	for ; srcIdx < endIndex; srcIdx++ {
		copy(buf[dstIdx:dstIdx+2], hexdumpTable[(array[srcIdx]&0xFF)<<1:])
		dstIdx += 2
	}

	return string(buf)
}
