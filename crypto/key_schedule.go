/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package crypto

import "fmt"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// A key schedule is a slice of KeyScheduleItem.
type KeySchedule = []KeyScheduleItem

// an object to convey a single substitution cipher key.
type KeyScheduleItem struct {
	KeyShift int    `json:"keyShift"`
	KeyChar  rune   `json:"keyChar"`
	Comment  string `json:"comment,omitempty"`
	Tabula   string `json:"tabula,omitempty"`
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer displaying the KeySchedule object in a single line.
func (ks KeyScheduleItem) String() string {
	return fmt.Sprintf("%02d %c %s (%s)", ks.KeyShift, ks.KeyChar, ks.Tabula, ks.Comment)
}
