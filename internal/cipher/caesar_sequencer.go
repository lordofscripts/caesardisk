/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A key sequencer for the plain Caesar cipher in which only ONE key
 * is used.
 * Version: 1
 * Class: Caesar (substitution cipher)
 * Mode : plain Caesar
 * Type : Monoalphabetic cipher (1)
 *-----------------------------------------------------------------*/
package cipher

import (
	"fmt"

	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*CaesarSequencer)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// A Caesar sequencer sequences the key value through the message
// encoding & decoding. It implements IKeySequencer. Only ONE
// Caesar table is needed.
type CaesarSequencer struct {
	params  *CaesarParameters
	isValid bool
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) A new instance of a plain Caesar key sequencer. For plain
// Caesar the same key is used throughout the message. The Caesar
// cipher is monosyllabic, one key and one transcode.
func NewCaesarSequencer(par *CaesarParameters) *CaesarSequencer {
	return &CaesarSequencer{
		params:  par,
		isValid: false,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (cs *CaesarSequencer) String() string {
	char, _ := cs.params.Alphabet.Character(cs.params.KeyValue)

	return fmt.Sprintf("Caesar(%c|%d)", char, cs.params.KeyValue)
}

func (cs *CaesarSequencer) Validate() error {
	if cs.params.KeyValue < 0 {
		return fmt.Errorf("cannot have negative keys")
	}
	if cs.params.KeyValue >= cs.params.Alphabet.Length() {
		logx.Print("WARN: keys greater than alphabet length are rewound")
		cs.params.KeyValue %= cs.params.Alphabet.Length()
	}
	if cs.params.KeyValue == 0 {
		logx.Print("WARN: A shift of zero does not transcode")
	}
	cs.isValid = true
	return nil
}

// the range of valid key values for Caesar cipher. For plain
// Caesar the prim parameter is ignored.
func (cs *CaesarSequencer) KeyRange() (min, max int) {
	min = 0
	max = cs.params.Alphabet.Length() - 1
	return
}

func (cs *CaesarSequencer) GetParams() *CaesarParameters {
	return cs.params
}

// Get the next key to use, should only be called if a message's
// character is not skipped. Characters are skipped if they are not
// part of the encoding alphabet, and thus do not participate in the
// key computation.
// Note: For plain Caesar we use the same key throughout the message.
func (cs *CaesarSequencer) NextKey() int {
	return cs.params.KeyValue
}

// Caesar is a monoalphabetic substitution cipher
func (cs *CaesarSequencer) IsPolyalphabetic() bool {
	return false
}

// whether the Offset parameter is used in key sequencing
func (cs *CaesarSequencer) IsOffsetRequired() bool {
	return false
}
