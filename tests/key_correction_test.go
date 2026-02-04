package tests

import (
	"testing"

	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/caesardisk/internal/cipher"
)

// Try the CaesarCorrection which for any given alphabet, corrects
// or normalizes the passed Key Shift. It corrects values that are
// out of range by returning their modulo, and if
// negative, it returns its equivalent complement.
func Test_CaesarCorrection(t *testing.T) {
	WithAlphabet := caesardisk.AlphabetFactory("English")
	N := WithAlphabet.Length()

	vectors := []struct {
		Original    int  // passed to corrector
		Expect      int  // expected after correction (if applied)
		WithWarning bool // whether a corrector warning is expected
	}{
		{8, 8, false},   // no corrections done
		{N, 0, true},    // out of range
		{-5, 21, true},  // negative within range
		{-28, 24, true}, // negative out of range
	}

	for i, v := range vectors {
		got, warn := cipher.CaesarCorrection(v.Original, WithAlphabet)
		if warn != nil && !v.WithWarning {
			t.Errorf("#%d Exp:%d Got:%d no warning expected but: %s", i+1, v.Expect, got, warn)
		}
		if warn == nil && v.WithWarning {
			t.Errorf("#%d Exp:%d Got:%d warning expected but got none", i+1, v.Expect, got)
		}
		if !v.WithWarning && v.Expect != got {
			t.Errorf("#%d Exp:%d Got:%d", i+1, v.Expect, got)
		}
		// now ensure the caesar encoded value is the same for both
		const DUMMY = "Hello"
		ctrl := crypto.NewCipherController(WithAlphabet, nil)
		enc1, _ := ctrl.Encrypt(crypto.CaesarMode, DUMMY, v.Expect)
		enc2, _ := ctrl.Encrypt(crypto.CaesarMode, DUMMY, got)
		if enc1 != enc2 {
			t.Errorf("#%d different result for '%s' != '%s'", i+1, enc1, enc2)
		}
	}
}

// Tests the Didimus corrector which already relies on the CaesarCorrector.
// we therefore do not need subvectors for invalid main key shifts.
func Test_DidimusCorrection(t *testing.T) {
	WithAlphabet := caesardisk.AlphabetFactory("English")
	N := WithAlphabet.Length()

	vectors := []struct {
		Main        int  // Main key
		Original    int  // Offset passed to corrector (Offset)
		Expect      int  // expected after correction (if applied)
		WithWarning bool // whether a corrector warning is expected
	}{
		{10, 5, 15, false},     // no corrections done, positive offset
		{10, -5, 5, false},     // negative offset corrected
		{N - 1, -1, 24, false}, // no corrections done, negative offset
		{N - 1, 1, 1, true},    // out of range above, 0 bounced to 1
		{N - 1, 2, 1, true},    // out of range above
	}

	for i, v := range vectors {
		_, _, got, warn := cipher.DidimusCorrection(v.Main, v.Original, WithAlphabet)
		if warn != nil && !v.WithWarning {
			t.Errorf("#%d Exp:%d Got:%d no warning expected but: %s", i+1, v.Expect, got, warn)
		}
		if warn == nil && v.WithWarning {
			t.Errorf("#%d Exp:%d Got:%d warning expected but got none", i+1, v.Expect, got)
		}
		if !v.WithWarning && v.Expect != got {
			t.Errorf("#%d Exp:%d Got:%d", i+1, v.Expect, got)
		}
		// now ensure the caesar encoded value is the same for both
		const DUMMY = "Hello"
		ctrl := crypto.NewCipherController(WithAlphabet, nil)
		enc1, _ := ctrl.Encrypt(crypto.DidimusMode, DUMMY, v.Main, v.Expect)
		enc2, _ := ctrl.Encrypt(crypto.DidimusMode, DUMMY, v.Main, got)
		if enc1 != enc2 {
			t.Errorf("#%d different result for '%s' != '%s'", i+1, enc1, enc2)
		}
	}
}
