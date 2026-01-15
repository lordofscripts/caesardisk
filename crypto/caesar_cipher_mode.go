/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package crypto

import "fmt"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	CaesarMode CaesarCipherMode = iota
	DidimusMode
	FibonacciMode
	PrimusMode
)

var (
	cipherModeToString map[CaesarCipherMode]string = map[CaesarCipherMode]string{
		CaesarMode:    "Caesar",
		DidimusMode:   "Didimus",
		FibonacciMode: "Fibonacci",
		PrimusMode:    "Primus",
	}
	cipherModeFromString map[string]CaesarCipherMode = map[string]CaesarCipherMode{
		"Caesar":    CaesarMode,
		"Didimus":   DidimusMode,
		"Fibonacci": FibonacciMode,
		"Primus":    PrimusMode,
	}
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarCipherMode uint8

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (cm CaesarCipherMode) String() string {
	return cipherModeToString[cm]
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func ParseCipherMode(s string) (CaesarCipherMode, error) {
	if val, ok := cipherModeFromString[s]; ok {
		return val, nil
	} else {
		return CaesarMode, fmt.Errorf("invalid CipherMode name: %s", s)
	}
}
