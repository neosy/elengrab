package inmemory

import (
	"encoding/binary"
	"hash/maphash"
)

var hashSeed = maphash.MakeSeed()

// dataHash returns a fast non-cryptographic hash of data
func dataHash(data []byte) uint64 {
	return maphash.Bytes(hashSeed, data)
}

// dataHashString returns the hash of data as an 8-byte binary string (for map keys)
func dataHashString(data []byte) string {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], dataHash(data))
	return string(buf[:])
}
