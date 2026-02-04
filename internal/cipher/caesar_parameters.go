/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cipher

import (
	"fmt"

	"github.com/lordofscripts/caesardisk"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IParamProvider interface {
	// Get the current alphabet.
	GetAlpha() *caesardisk.AlphabetModel
	// Get the main key in character and key shift.
	GetKey() (rune, int)
	// Get the alternate key in character and key shift based on
	// the existing main key shift and offset.
	GetAltKey() (rune, int)
	// The the source encoding alphabet. If there was no previous alpha,
	// key & offset are reset. If the new alpha is shorter than the previous
	// and the key shift is out or range, it gets adjusted. It returns
	// the newly computed main key shift and computed alternate key.
	SetAlpha(*caesardisk.AlphabetModel) (newKeyShift, newAltKey int)
	// set the main keyShift and returns the adjusted (actual) effective
	// key shift corrected for alphabet length. Positive and Negative
	// values produce different results!
	SetKey(keyShift int) int
	// compute the new alternate key based on the existing main key
	// and the provided offset. Positive and Negative values produce
	// different results. It returns the computed alternate key shift.
	SetAltKeyOffset(offset int) int
}

var _ IParamProvider = (*CaesarParameters)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarParameters struct {
	Alphabet *caesardisk.AlphabetModel
	KeyValue int
	Offset   int // not used for plain Caesar, just Didimus & Fibonacci
	altKey   int // derived from key+offset not used in plain Caesar
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) a new instance of Caesar cipher parameters for the
// chosen alphabet.
func NewCaesarParameters(alphabet *caesardisk.AlphabetModel) *CaesarParameters {
	return &CaesarParameters{
		Alphabet: alphabet,
		KeyValue: 0,
		Offset:   -1,
		altKey:   0,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer and returns the current key in character and
// shift value together with the actual alphabet letters or symbols.
func (c *CaesarParameters) String() string {
	keyLetter, _ := c.Alphabet.Character(c.KeyValue)
	return fmt.Sprintf("(%02d|%c) %s", c.KeyValue, keyLetter, c.Alphabet)
}

/* ----------------------------------------------------------------
 *			I n t e r f a c e	M e t h o d s
 *-----------------------------------------------------------------*/

func (c *CaesarParameters) GetAlpha() *caesardisk.AlphabetModel {
	return c.Alphabet
}

func (c *CaesarParameters) GetKey() (rune, int) {
	char, _ := c.Alphabet.Character(c.KeyValue)
	return char, c.KeyValue
}

// get the actual alternate key (derived from key, offset & alphabet)
// in its letter and key shift forms.
func (c *CaesarParameters) GetAltKey() (rune, int) {
	char, _ := c.Alphabet.Character(c.altKey)
	return char, c.altKey
}

// The the source encoding alphabet. If there was no previous alpha,
// key & offset are reset. If the new alpha is shorter than the previous
// and the key shift is out or range, it gets adjusted. It returns
// the newly computed main key shift and computed alternate key.
func (c *CaesarParameters) SetAlpha(alpha *caesardisk.AlphabetModel) (int, int) {
	if c.Alphabet == nil {
		c.Alphabet = alpha
		c.KeyValue = 0
		c.Offset = 0
	} else {
		oldN := c.Alphabet.Length()
		newN := alpha.Length()
		c.Alphabet = alpha

		if oldN > newN {
			// new alphabet is shorter may need K&O adjustment
			if c.KeyValue >= newN {
				c.KeyValue = c.adjustKey(c.KeyValue)
			}
			c.altKey = c.adjustOffset(c.Offset)
		}
	}

	return c.KeyValue, c.altKey
}

// set the main keyShift and returns the adjusted (actual) effective
// key shift corrected for alphabet length. Positive and Negative
// values produce different results!
func (c *CaesarParameters) SetKey(shift int) int { // @audit deprecate!
	c.KeyValue = c.adjustKey(shift)
	return c.KeyValue
}

// compute the new alternate key based on the existing main key
// and the provided offset. Positive and Negative values produce
// different results. It returns the computed alternate key shift.
func (c *CaesarParameters) SetAltKeyOffset(offset int) int { // @audit deprecate!
	c.Offset = offset
	c.altKey = c.adjustOffset(offset)
	return c.altKey
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// with the provided key shift (main key) adjust it so that it falls
// within the range of valid keys (0..N where N is length of alphabet).
// This works with both positive and negative offsets.
// Returns the adjusted key.
func (c *CaesarParameters) adjustKey(shift int) int {
	return KeyAdjuster(c.Alphabet.Length(), shift)
}

// with the provided the offset compute the alternate key (shift).
// This works with positive and negative offsets.
// Returns the alternate key.
func (c *CaesarParameters) adjustOffset(offset int) int {
	return OffsetAdjuster(c.Alphabet.Length(), c.KeyValue, offset)
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// For an alphabet length N adjust the provided shift key value so
// that it falls within the valid range of keys 0..N-1 by taking
// the modulo operation.
func KeyAdjuster(N, shift int) int {
	// For N=26 k=-2 becomes k=24
	// k=-27 becomes k=25
	return shift % N
}

// For alphabet length N the valid range of shift is 0..N-1.
// It returns the alternate key which is computed from the
// main shift plus the offset modulo N. Note, it does not
// return an adjusted offset but the alternate key!
func OffsetAdjuster(N, shift, offset int) int { // @audit deprecate and use CaesarCorrection()
	// For N=26 (k:0..25)
	// 	k=1 ofs=5 ofs'=6
	// 	k=1 ofs=24 ofs'=25
	// 	k=1 ofs=25 ofs'=0
	// 	k=1 ofs=26 ofs'=1
	// 	k=1 ofs=-1 ofs'=0
	// 	k=1 ofs=-2 ofs'=25
	// 	k=1 ofs=-25 ofs'=2
	// 	k=1 ofs=-26 ofs'=1
	return (shift + offset) % N
}
