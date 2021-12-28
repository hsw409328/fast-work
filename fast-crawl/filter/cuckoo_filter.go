package filter

import (
	"fmt"
	metro "github.com/dgryski/go-metro"
	"math/bits"
	"math/rand"
)

var (
	cuckooCrawlKey = "cuckooCrawlKey"
	altHash        = [256]uint{}
	masks          = [65]uint{}
	maxCuckooCount = 500
)

type CuckooFilter struct {
	// 一维数组
	buckets []bucket
	// Filter 中已插入的元素
	count uint
	// 数组buckets长度中对应二进制包含0的个数
	bucketPow uint
}

func init() {
	for i := 0; i < 256; i++ {
		altHash[i] = (uint(metro.Hash64([]byte{byte(i)}, 1337)))
	}
	for i := uint(0); i < 64; i++ {
		masks[i] = (1 << i) - 1
	}
}

func NewFilter(capacity uint) *CuckooFilter {
	capacity = getNextPow2(uint64(capacity)) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	buckets := make([]bucket, capacity)
	return &CuckooFilter{
		buckets:   buckets,
		count:     0,
		bucketPow: uint(bits.TrailingZeros(capacity)),
	}
}

// Lookup returns true if data is in the counter
func (cf *CuckooFilter) Lookup(data []byte) bool {
	i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
	if cf.buckets[i1].getFingerprintIndex(fp) > -1 {
		return true
	}
	i2 := getAltIndex(fp, i1, cf.bucketPow)
	return cf.buckets[i2].getFingerprintIndex(fp) > -1
}

// Reset ...
func (cf *CuckooFilter) Reset() {
	for i := range cf.buckets {
		cf.buckets[i].reset()
	}
	cf.count = 0
}

func randi(i1, i2 uint) uint {
	if rand.Intn(2) == 0 {
		return i1
	}
	return i2
}

// Insert inserts data into the counter and returns true upon success
func (cf *CuckooFilter) Insert(data []byte) bool {
	i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
	if cf.insert(fp, i1) {
		return true
	}
	i2 := getAltIndex(fp, i1, cf.bucketPow)
	if cf.insert(fp, i2) {
		return true
	}
	return cf.reinsert(fp, randi(i1, i2))
}

// InsertUnique inserts data into the counter if not exists and returns true upon success
func (cf *CuckooFilter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	return cf.Insert(data)
}

func (cf *CuckooFilter) insert(fp fingerprint, i uint) bool {
	if cf.buckets[i].insert(fp) {
		cf.count++
		return true
	}
	return false
}

func (cf *CuckooFilter) reinsert(fp fingerprint, i uint) bool {
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		oldfp := fp
		fp = cf.buckets[i][j]
		cf.buckets[i][j] = oldfp

		// look in the alternate location for that random element
		i = getAltIndex(fp, i, cf.bucketPow)
		if cf.insert(fp, i) {
			return true
		}
	}
	return false
}

// Delete data from counter if exists and return if deleted or not
func (cf *CuckooFilter) Delete(data []byte) bool {
	i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
	if cf.delete(fp, i1) {
		return true
	}
	i2 := getAltIndex(fp, i1, cf.bucketPow)
	return cf.delete(fp, i2)
}

func (cf *CuckooFilter) delete(fp fingerprint, i uint) bool {
	if cf.buckets[i].delete(fp) {
		if cf.count > 0 {
			cf.count--
		}
		return true
	}
	return false
}

func (cf *CuckooFilter) Count() uint {
	return cf.count
}

func (cf *CuckooFilter) Encode() []byte {
	bytes := make([]byte, len(cf.buckets)*bucketSize)
	for i, b := range cf.buckets {
		for j, f := range b {
			index := (i * len(b)) + j
			bytes[index] = byte(f)
		}
	}
	return bytes
}

func Decode(bytes []byte) (*CuckooFilter, error) {
	var count uint
	if len(bytes)%bucketSize != 0 {
		return nil, fmt.Errorf("expected bytes to be multiple of %d, got %d", bucketSize, len(bytes))
	}
	buckets := make([]bucket, len(bytes)/4)
	for i, b := range buckets {
		for j := range b {
			index := (i * len(b)) + j
			if bytes[index] != 0 {
				buckets[i][j] = fingerprint(bytes[index])
				count++
			}
		}
	}
	return &CuckooFilter{
		buckets:   buckets,
		count:     count,
		bucketPow: uint(bits.TrailingZeros(uint(len(buckets)))),
	}, nil
}

func getAltIndex(fp fingerprint, i uint, bucketPow uint) uint {
	mask := masks[bucketPow]
	hash := altHash[fp] & mask
	return (i & mask) ^ hash
}

func getFingerprint(hash uint64) byte {
	// Use least significant bits for fingerprint.
	fp := byte(hash%255 + 1)
	return fp
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndexAndFingerprint(data []byte, bucketPow uint) (uint, fingerprint) {
	hash := metro.Hash64(data, 1337)
	fp := getFingerprint(hash)
	// Use most significant bits for deriving index.
	i1 := uint(hash>>32) & masks[bucketPow]
	return i1, fingerprint(fp)
}

func getNextPow2(n uint64) uint {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return uint(n)
}
