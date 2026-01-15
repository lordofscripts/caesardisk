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

// The AlphabetFactory function returns an instance of the
// requested alphabet. The nameOrCode can be the full name
// or its 2-5 letter id. It is case-sensitive. Returns nil
// if it cannot comply. It sets the Name if it is empty.
func AlphabetFactory(nameOrCode string) *AlphabetModel {
	var mdl *AlphabetModel
	switch nameOrCode {
	case "EN":
		fallthrough
	case "English":
		mdl = NewAlphabetModelCased(Alpha_EN)
		mdl.Name = "English"

	case "ES":
		fallthrough
	case "Español":
		mdl = NewAlphabetModelCased(Alpha_ES_DUAL)
		mdl.Name = "Español"

	case "ES-XTR":
		fallthrough
	case "Español con acentos":
		mdl = NewAlphabetModelCased(Alpha_ES)
		mdl.Name = "Español con acentos"

	case "CZ":
		fallthrough
	case "Czech":
		mdl = NewAlphabetModelCased(Alpha_CZ)
		mdl.Name = "Czech"

	case "DE":
		fallthrough
	case "Deutsch":
		mdl = NewAlphabetModelCased(Alpha_DE)
		mdl.Name = "Deutsch"

	case "IT":
		fallthrough
	case "Italiano":
		mdl = NewAlphabetModelCased(Alpha_IT)
		mdl.Name = "Italiano"

	case "PT":
		fallthrough
	case "Portuguese", "Português":
		mdl = NewAlphabetModelCased(Alpha_PT)
		mdl.Name = "Português"

	case "RU":
		fallthrough
	case "Cyrillic", "Russian", "Ukrainian":
		mdl = NewAlphabetModelCased(Alpha_RU)
		mdl.Name = "Russian"

	case "GR":
		fallthrough
	case "Greek":
		mdl = NewAlphabetModelCased(Alpha_GR)
		mdl.Name = "Greek"

	case "PU":
		fallthrough
	case "Punctuation (all)":
		mdl = NewAlphabetModelCased(Alpha_PU)
		mdl.Name = "Punctuation (all)"

	case "PU-ES":
		fallthrough
	case "Puntuacion para Español":
		mdl = NewAlphabetModelCased(Alpha_PU_DUAL_ES)
		mdl.Name = "Puntuacion para Español"

	case "PU-EN":
		fallthrough
	case "Punctuation for English":
		mdl = NewAlphabetModelCased(Alpha_PU_DUAL_EN)
		mdl.Name = "Punctuation for English"

	case "Runes":
		mdl = NewAlphabetModelCased(Alpha_RUNES)
		mdl.Name = "Runes"

	default:
		println("Unrecognized alphabet in factory: ", nameOrCode)
	}

	return mdl
}

// compare the alphabet contents with the predefined list.
// If found, set its Name and return that name.
func IdentifyAlphabet(α *AlphabetModel) string {
	if α.Name != "" {
		return α.Name
	}

	switch string(α.alphabet) {
	case Alpha_EN:
		α.Name = "English"
	case Alpha_ES_DUAL:
		α.Name = "Español"
	case Alpha_ES:
		α.Name = "Español con acentos"
	case Alpha_CZ:
		α.Name = "Czech"
	case Alpha_DE:
		α.Name = "Deutsch"
	case Alpha_IT:
		α.Name = "Italiano"
	case Alpha_PT:
		α.Name = "Português"
	case Alpha_RU:
		α.Name = "Russian"
	case Alpha_GR:
		α.Name = "Greek"
	case Alpha_PU:
		α.Name = "Punctuation (all)"
	case Alpha_PU_DUAL_ES:
		α.Name = "Puntuacion para Español"
	case Alpha_PU_DUAL_EN:
		α.Name = "Punctuation for English"
	case Alpha_RUNES:
		α.Name = "Runes"
	}

	return α.Name
}
