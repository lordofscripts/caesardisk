/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesar key sequencer for Primus mode. It uses the main Caesar
 * key and for subsequent encodeable characters in the input, it uses
 * a shift composed of the Main key/shift plus the current term in the
 * list of 10 primes (2..31). After the last prime is used, it resets
 * and starts again. The effective
 * key shift is always normalized to the current alphabet.
 * Version: 1
 * Class: Caesar (substitution cipher)
 * Mode : Primus
 * Type : Poly-alphabetic cipher (up to 12)
 *-----------------------------------------------------------------*/
package cipher

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

// List of first primes. Although 0 is not a prime, we need to ensure
// that the main Caesar shift/key is always the first.
var primes []int = []int{0, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*PrimusSequencer)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// Primus is like Fibonacci but instead of a Fibonacci sequence, it
// uses a sequence of the first 11 prime numbers as an offset to the
// main Caesar key. Although 0 is not a prime, the first in the series
// is always zero, thus the main Caesar key.
type PrimusSequencer struct {
	CaesarSequencer
	// The current Prime-number term index as it progresses. (0..11)
	termIndex int
	// The maximum number of Prime-number terms to be used. Zero
	// is equivalent to plain Caesar, One is equivalent to Didimus
	// with offset set to 1. Suggested 2..11
	maxPrimes int
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (Ctor) a new instance of the Primus sequencer for the Caesar encoder.
// The Offset value is automatically corrected via modulo to the maximum
// number of prime values. When Offset is zero we select the maximum (11)
func NewPrimusSequencer(par *CaesarParameters) *PrimusSequencer {
	// normalize Offset, and if zero use the maximum set of primes
	var maxPrimeTerms int = par.Offset % len(primes)
	if par.Offset == 0 {
		maxPrimeTerms = len(primes)
	}

	return &PrimusSequencer{
		CaesarSequencer: *NewCaesarSequencer(par),
		termIndex:       0,
		maxPrimes:       maxPrimeTerms,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (ps *PrimusSequencer) String() string {
	charM, _ := ps.params.Alphabet.Character(ps.params.KeyValue)

	return fmt.Sprintf("Primus(%c|%d,P(%d))",
		charM, ps.params.KeyValue, ps.maxPrimes)
}

func (ps *PrimusSequencer) Validate() error {
	return ps.CaesarSequencer.Validate()
}

// The valid range of Primus keys.
func (ps *PrimusSequencer) KeyRange() (min, max int) {
	return ps.CaesarSequencer.KeyRange()
}

func (ps *PrimusSequencer) GetParams() *CaesarParameters {
	return ps.params
}

func (ps *PrimusSequencer) NextKey() int {
	// ensure we wrap correctly
	_, max := ps.CaesarSequencer.KeyRange()
	// new shift but wrapped to domain
	keyShift := (ps.params.KeyValue + primes[ps.termIndex]) % (max + 1)
	// now update term for next call
	newIndex := ps.termIndex + 1
	// rewind when we reach the maximum number of available primes,
	// or the maximum selected number of primes.
	if newIndex >= len(primes) || newIndex > ps.maxPrimes {
		newIndex = 0
	}
	ps.termIndex = newIndex

	return keyShift
}

// The internal key schedule
func (ps *PrimusSequencer) GetRawKeySchedule() []KeyScheduleItemInt {
	qty := len(primes)
	fakeSeq := NewPrimusSequencer(ps.params)
	schedule := make([]KeyScheduleItemInt, qty)
	for i := range qty {
		schedule[i] = KeyScheduleItemInt{
			KeyShift: fakeSeq.NextKey(),
			Comment:  fmt.Sprintf("#%d", i)}
	}
	return schedule
}

// Primus is a polyalphabetic substitution cipher
func (ps *PrimusSequencer) IsPolyalphabetic() bool {
	return true
}

// whether the Offset parameter is used in key sequencing.
// Primus uses the Offset value to indicate how many prime
// numbers (past the initial main key) should be used.
func (ps *PrimusSequencer) IsOffsetRequired() bool {
	return true
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// The maximum number of Prime-number terms available in the algorithm.
func PrimusMaximus() int {
	return len(primes)
}
