package util

import (
	"encoding/binary"
	"fmt"
)

func Factorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	return uint64(n) * Factorial(n-1)
}

func GetMapKey(w, h, l int) string {
	return fmt.Sprintf("%d-%d-%d", w, h, l)
}

func Murmur3Hash32(key []byte, seed uint32) uint32 {
	const (
		c1 = 0xcc9e2d51
		c2 = 0x1b873593
		r1 = 15
		r2 = 13
		m  = 5
		n  = 0xe6546b64
	)

	var (
		length = len(key)
		h1     = seed
	)

	for i := 0; i+4 <= length; i += 4 {
		k1 := binary.LittleEndian.Uint32(key[i:])
		k1 *= c1
		k1 = (k1 << r1) | (k1 >> (32 - r1))
		k1 *= c2

		h1 ^= k1
		h1 = (h1 << r2) | (h1 >> (32 - r2))
		h1 = h1*m + n
	}

	var k1 uint32
	switch length & 3 {
	case 3:
		k1 ^= uint32(key[length&^3+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(key[length&^3+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(key[length&^3])
		k1 *= c1
		k1 = (k1 << r1) | (k1 >> (32 - r1))
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= uint32(length)
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}
