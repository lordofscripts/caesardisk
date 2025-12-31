/* -----------------------------------------------------------------
 *					P u b l i c   D o m a i n
 *				  			Copyleft
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * An XXHash64 implementation. It is present here in my BIP39 package
 * to provide a 64-bit hash that has less collisions than CRC64, and
 * thus more suitable for the Caesarium.
 *-----------------------------------------------------------------*/
package hash

import (
	"encoding/binary"
	"fmt"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	prime1 = uint64(11400714785074694791)
	prime2 = uint64(14029467366897019727)
	prime3 = uint64(16095879293928366613)
	//prime4 = uint64(9650029242287838081)
	prime5 = uint64(2870177450012600261)
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IDigest interface {
	// Render the final Hash or Checksume as a string.
	fmt.Stringer
	// Update the digest with new data
	Update(data []byte)
	// The length of the digest in bytes
	Length() int
}

var _ IDigest = (*XXH64)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// XXH64 represents the XXH64 state
type XXH64 struct {
	v1, v2, v3, v4 uint64
	totalLength    uint64
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// NewXXH64 initializes a new hash state
func NewXXH64(seed uint64) *XXH64 {
	h := &XXH64{
		v1: seed + prime1 + prime2,
		v2: seed + prime2,
		v3: seed,
		v4: seed - prime1,
	}
	return h
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implement fmt.Stringer on IDigest rendering the XXH64 checksum
// as a 16-digit hex string.
func (h *XXH64) String() string {
	return fmt.Sprintf("%16X", h.Digest())
}

// would return 8 which means 8 bytes (64-bits)
func (h *XXH64) Length() int {
	return 8
}

// Update hashes the input data
func (h *XXH64) Update(data []byte) {
	h.totalLength += uint64(len(data))
	for len(data) >= 32 {
		h.v1 += binary.BigEndian.Uint64(data[:8]) * prime2
		h.v1 = (h.v1<<31 | h.v1>>33) * prime1
		h.v1 += h.v2

		h.v2 += binary.BigEndian.Uint64(data[8:16]) * prime2
		h.v2 = (h.v2<<31 | h.v2>>33) * prime1
		h.v2 += h.v3

		h.v3 += binary.BigEndian.Uint64(data[16:24]) * prime2
		h.v3 = (h.v3<<31 | h.v3>>33) * prime1
		h.v3 += h.v4

		h.v4 += binary.BigEndian.Uint64(data[24:32]) * prime2
		h.v4 = (h.v4<<31 | h.v4>>33) * prime1
		h.v4 += h.v1

		data = data[32:]
	}

	// Handle remaining bytes
	for _, b := range data {
		h.v1 += uint64(b) * prime5
		h.v1 = (h.v1<<11 | h.v1>>53) * prime1
	}
}

// Digest returns the final hash value
func (h *XXH64) Digest() uint64 {
	h.v1 += h.v2 + h.v3 + h.v4
	h.v1 = (h.v1 ^ (h.v1 >> 33)) * prime2
	h.v1 ^= h.v1 >> 29
	h.v1 *= prime3
	h.v1 += prime5
	return h.v1
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/*
func demo() {
	hasher := NewXXH64(0) // Initialize with a seed
	data := []byte("Hello, World!")
	hasher.Update(data)          // Update the hash with your data
	hashValue := hasher.Digest() // Get the final hash value
	fmt.Printf("Hash: %x\n", hashValue)
}
*/
