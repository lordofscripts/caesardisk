/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesar key sequencer for Fibonacci mode. It uses the main Caesar
 * key and for each subsequent encodeable character in the input, it
 * uses a new key made by adding the main key to the current Fibonacci
 * term. The algorithm uses a 10-term Fibonacci series. The effective
 * key shift is always normalized to the current alphabet.
 * Version: 1
 * Class: Caesar (substitution cipher)
 * Mode : Fibonacci
 * Type : Polyalphabetic cipher (up to 9)
 *-----------------------------------------------------------------*/
package cipher

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

// we use a 10-term Fibonacci sequence for the Fibonacci Key Sequencer
var fibonacci []int = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*FibonacciSequencer)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// Fibonacci uses a poly-alphabetic sequencer in which it starts
// with the main Caesar key, and for subsequent encodeable characters
// it uses key+F(x) where F(x) is the Xth term of a Fibonacci series.
// Our Fibonacci series has 10 terms.
type FibonacciSequencer struct {
	CaesarSequencer
	termIndex int
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (Ctor) a new instance of the Caesar encoder using Fibonacci mode.
func NewFibonacciSequencer(par *CaesarParameters) *FibonacciSequencer {
	return &FibonacciSequencer{
		CaesarSequencer: *NewCaesarSequencer(par),
		termIndex:       0,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (fs *FibonacciSequencer) String() string {
	charM, _ := fs.params.Alphabet.Character(fs.params.KeyValue)

	return fmt.Sprintf("Fibonacci(%c|%d,F(10))",
		charM, fs.params.KeyValue)
}

func (fs *FibonacciSequencer) Validate() error {
	return fs.CaesarSequencer.Validate()
}

// The valid range of Didimus keys.
func (fs *FibonacciSequencer) KeyRange() (min, max int) {
	return fs.CaesarSequencer.KeyRange()
}

func (fs *FibonacciSequencer) GetParams() *CaesarParameters {
	return fs.params
}

func (fs *FibonacciSequencer) NextKey() int {
	// ensure we wrap correctly
	_, max := fs.CaesarSequencer.KeyRange()
	// new shift but wrapped to domain
	keyShift := (fs.params.KeyValue + fibonacci[fs.termIndex]) % (max + 1)
	// now update term for next call
	newIndex := fs.termIndex + 1
	if newIndex >= len(fibonacci) {
		newIndex = 0
	}
	fs.termIndex = newIndex

	return keyShift
}

// Fibonacci is a polyalphabetic substitution cipher
func (fs *FibonacciSequencer) IsPolyalphabetic() bool {
	return true
}

// whether the Offset parameter is used in key sequencing.
// Fibonacci v1 does not use Offset.
func (fs *FibonacciSequencer) IsOffsetRequired() bool {
	return false
}
