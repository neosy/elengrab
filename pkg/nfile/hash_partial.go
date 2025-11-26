package nfile

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
)

// HashPartial computes a single combined hash from several file blocks.
// It is very fast because it does not read the entire file.
//
// blocks — number of blocks to read (e.g. 3: start, middle, end)
// blockSize — size of each block in bytes (e.g. 1*1024*1024 = 1MB)
func HashPartial(path string, blocks int, blockSize int64) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	size := fi.Size()
	if size == 0 {
		// Empty file → deterministic hash
		empty := sha256.Sum256([]byte("empty-file"))
		return hex.EncodeToString(empty[:]), nil
	}

	h := sha256.New()

	// Write file size to the hash to reduce collisions
	bufSize := make([]byte, 8)
	binary.LittleEndian.PutUint64(bufSize, uint64(size))
	h.Write(bufSize)

	// If the file is small, we read the whole thing
	if size <= int64(blocks)*blockSize {
		buf := make([]byte, size)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		h.Write(buf)
		return hex.EncodeToString(h.Sum(nil)), nil
	}

	// Large file — read only blocks
	for i := range blocks {
		// calculate offset of the block
		offset := (size - blockSize) * int64(i) / int64(blocks-1)
		if offset < 0 {
			offset = 0
		}
		if offset > size-blockSize {
			offset = size - blockSize
		}

		// seek to block position
		_, err := f.Seek(offset, io.SeekStart)
		if err != nil {
			return "", err
		}

		// read block
		buf := make([]byte, blockSize)
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}

		// hash block
		h.Write(buf[:n])
	}

	// return combined final hash
	return hex.EncodeToString(h.Sum(nil)), nil
}
