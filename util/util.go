package util

// HashCode returns a hash code based on the contents of the specified array.
func HashCode(a []byte) uint32 {
	if len(a) == 0 {
		return 0
	}

	var result uint32 = 1
	for _, element := range a {
		result = 31*result + uint32(element)
	}

	return result
}

// ReverseUint64 returns the value obtained by reversing the order of the bits
// in the two's complement binary representation of the specified long value.
func ReverseUint64(i uint64) uint64 {
	// HD, Figure 7-1
	i = (i&0x5555555555555555)<<1 | (i>>1)&0x5555555555555555
	i = (i&0x3333333333333333)<<2 | (i>>2)&0x3333333333333333
	i = (i&0x0f0f0f0f0f0f0f0f)<<4 | (i>>4)&0x0f0f0f0f0f0f0f0f
	i = (i&0x00ff00ff00ff00ff)<<8 | (i>>8)&0x00ff00ff00ff00ff
	return (i << 48) | ((i & 0xffff0000) << 16) | ((i >> 16) & 0xffff0000) | (i >> 48)
}
