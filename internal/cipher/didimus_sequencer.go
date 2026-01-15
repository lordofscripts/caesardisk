/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesar key sequencer for Didimus mode. It uses the main Caesar
 * key for even-positioned characters that are encodeable, and the
 * alternate key made by adding the Offset to the Main key for all
 * odd-positioned characters that are encodeable. The odd/even condition
 * is not determined by the position in the input text, but by their
 * position after SKIPPING non-encodeable characters. The effective
 * key shift is always normalized to the current alphabet.
 * Version: 1
 * Class: Caesar (substitution cipher)
 * Mode : Didimus
 * Type : Bi-alphabetic cipher (2)
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

var _ IKeySequencer = (*DidimusSequencer)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// Didimus uses a bi-alphabetic sequencer with the same Caesar Key
// for encodeable even characters, and using a key+offset for the
// odd encodeable characters. Therefore, it uses 2 Caesar tables.
type DidimusSequencer struct {
	CaesarSequencer
	isEvenPosition bool
	altKey         int
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) A new instance of a Didimus key sequencer. Didimus is part
// of the Caesar cipher family but is bi-syllabic, the normal Caesar
// key is used for even characters (provided they are in the encoding
// alphabet), an
func NewDidimusSequencer(par *CaesarParameters) *DidimusSequencer {
	return &DidimusSequencer{
		CaesarSequencer: *NewCaesarSequencer(par),
		isEvenPosition:  false,
		altKey:          0,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (ds *DidimusSequencer) String() string {
	var charM, charA rune
	charM, _ = ds.params.Alphabet.Character(ds.params.KeyValue)
	charA, _ = ds.params.Alphabet.Character(ds.altKey)

	return fmt.Sprintf("Didimus(%c|%d,%c|%d)",
		charM, ds.params.KeyValue,
		charA, ds.altKey)
}

func (ds *DidimusSequencer) Validate() error {
	err := ds.CaesarSequencer.Validate()
	if err != nil {
		return err
	}
	ds.isValid = false

	// now validate Offset
	if ds.params.Offset <= 0 {
		err = fmt.Errorf("Didimus needs a positive non-zero offset: %d", ds.params.Offset)
	} else {
		alternate := ds.params.KeyValue + ds.params.Offset
		if alternate >= ds.params.Alphabet.Length() {
			logx.Print("WARN: offsets greater than alphabet length are rewound")
			alternate %= ds.params.Alphabet.Length()
		}
		if alternate == 0 {
			alternate = 1
		}
		ds.altKey = alternate
		ds.isValid = true
	}

	return nil
}

// The valid range of Didimus keys.
func (ds *DidimusSequencer) KeyRange() (min, max int) {
	return ds.CaesarSequencer.KeyRange()
}

func (ds *DidimusSequencer) GetParams() *CaesarParameters {
	return ds.params
}

func (ds *DidimusSequencer) NextKey() int {
	var keyShift int = 0
	ds.isEvenPosition = !ds.isEvenPosition // toggle, first (0) is Even

	if ds.isEvenPosition {
		keyShift = ds.params.KeyValue
	} else {
		keyShift = ds.altKey
	}

	return keyShift
}

// Didimus is a bi-alphabetic substitution cipher
func (ds *DidimusSequencer) IsPolyalphabetic() bool {
	return true
}

// whether the Offset parameter is used in key sequencing.
// Didimus uses the Offset to generate an alternate key by
// adding the Offset to the Main Key's shift.
func (ds *DidimusSequencer) IsOffsetRequired() bool {
	return true
}
