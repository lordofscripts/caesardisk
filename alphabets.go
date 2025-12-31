/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *						goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesardisk

import "strings"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Keep this in mind if you want an dual letter & symbol
	// disks that you can use with the same encoding key
	// (same length for both alphabets)
	// 0        1         2         3
	// 1234567890123456789012345678901
	// -------------------------------
	// ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚ
	// !"#$%&'()*+,-./ 0123456789:;<=>?
	Alpha_EN string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alpha_ES string = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚ"
	Alpha_CZ string = "ABCČDĎEFGHIJKLMNŇOPQRŘSŠTŤUVWXYÝZŽÁÉÍÓÚĚŮ"
	Alpha_DE string = "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ"
	Alpha_IT string = "ABCDEFGHILMNOPQRSTUVZÉÓÀÈÌÒÙ"
	Alpha_PT string = "ABCÇDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÀÂÊÔÃÕ"
	Alpha_GR string = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ"
	Alpha_RU string = "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
	Alpha_PU string = `!"#$%&'()*+,-./ 0123456789:;<=>?`

	Alpha_RUNES string = "ᚫᛒᚳᛞᛖᚠᚷᚻᛁᛃᛱᛚᛗᚾᚩᛈᛩᚱᛋᛏᚢᚡᚹᛪᛦᛎ"

	Alpha_ES_DUAL    = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZ"
	Alpha_PU_DUAL_ES = `!"#$%&()*+,-./ 0123456789=?` // to match length of ES_DUAL
	Alpha_PU_DUAL_EN = `!"#$%&()*+,-./ 0123456789?`  // to match length of EN
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func RuneString(latin string) string {
	const (
		LOOKUP_STD string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		RUNES      string = "\u16ab\u16d2\u16b3\u16de\u16d6\u16a0\u16b7\u16bb\u16c1\u16c3\u16f1\u16da\u16d7\u16be\u16a9\u16c8\u16e9\u16b1\u16cb\u16cf\u16a2\u16a1\u16b9\u16ea\u16e6\u16ce"
	)

	chars := []rune(strings.ToUpper(latin))
	runesLookup := []rune(RUNES)
	result := make([]rune, len(chars))

	for index, char := range chars {
		if strings.ContainsRune(LOOKUP_STD, char) {
			at := strings.IndexRune(LOOKUP_STD, char)
			result[index] = runesLookup[at]
		} else {
			result[index] = char
		}
	}

	return string(result)
}
