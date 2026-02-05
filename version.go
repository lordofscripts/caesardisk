/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							 goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Caesar Disk maker (CLI) & Encoder (GUI).
 *-----------------------------------------------------------------*/
package caesardisk

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/lordofscripts/goapp/app"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	// Actual version. Gets injected into vcs* vars below by the linker.
	// And gets injected in FyneApp.toml by "make version"
	MANUAL_VERSION string        = "1.5.0"
	VERSION_STATUS app.DevStatus = app.DevStatusReleased

	// Change these values accordingly
	NAME string = "Go Caesar Disk Maker"
	DESC string = "An extended Caesar Cipher for modern-day alphabets"
)

const (
	// Useful Unicode Characters
	CHR_COPYRIGHT       = '\u00a9'      // ©
	CHR_REGISTERED      = '\u00ae'      // ®
	CHR_GUILLEMET_L     = '\u00ab'      // «
	CHR_GUILLEMET_R     = '\u00bb'      // »
	CHR_TRADEMARK       = '\u2122'      // ™
	CHR_SAMARITAN       = '\u214f'      // ⅏
	CHR_PLACEOFINTEREST = '\u2318'      // ⌘
	CHR_HIGHVOLTAGE     = '\u26a1'      // ⚡
	CHR_MALTA_CROSS     = '\u2720'      // ✠
	CHR_TRIDENT         = rune(0x1f531) // 🔱
	CHR_SPLATTER        = rune(0x1fadf)
	CHR_WARNING         = '\u26a0' // ⚠
	CHR_EXCLAMATION     = '\u2757'
	CHR_SKULL           = '\u2620' // ☠

	CO1 = "odlamirG omidiD 5202)C("
	CO2 = "stpircS fO droL 5202)C("
	CO3 = "gnitirwnitsol"
)

var (
	vcsVersion  string = MANUAL_VERSION // automatically injected with linker
	vcsCommit   string
	vcsDate     string
	vcsBuildNum string
)

var (
	// NOTE: Change these values accordingly
	appVersion app.PackageVersion = app.NewPackageVersion(NAME, DESC, MANUAL_VERSION, VERSION_STATUS)

	// DO NOT CHANGE THESE!
	Version      string = appVersion.String()
	ShortVersion string = appVersion.Short()
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Funny LordOfScripts logo
func Logo() string {
	const (
		whiteStar rune = '\u269d' // ⚝
		unisex    rune = '\u26a5' // ⚥
		hotSpring rune = '\u2668' // ♨
		leftConv  rune = '\u269e' // ⚞
		rightConv rune = '\u269f' // ⚟
		eye       rune = '\u25d5' // ◕
		mouth     rune = '\u035c' // ͜	‿ \u203f
		skull     rune = '\u2620' // ☠
	)
	return fmt.Sprintf("%c%c%c %c%c", leftConv, eye, mouth, eye, rightConv)
	//fmt.Sprintf("(%c%c %c)", eye, mouth, eye)
}

// Hey! My time costs money too!
func BuyMeCoffee(coffee4 string) {
	appVersion.BuyMeCoffee(coffee4)
}

func Copyright(owner string) {
	appVersion.Copyright(Reverse(owner), CHR_TRIDENT)
	fmt.Println("\t\t\t\t", Logo())
}

// Injected build information
func BuildInfo() string {
	if len(vcsVersion) == 0 {
		vcsVersion = MANUAL_VERSION
	}
	return fmt.Sprintf("Build #%s %s %s", vcsBuildNum, vcsDate, vcsCommit)
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// get the current GO language version
func GoVersion() string {
	ver := strings.Replace(runtime.Version(), "go", "", -1)
	return ver
}

// retrieve the current GO language version and compare it
// to the minimum required. It returns the current version
// and whether the condition current >= min is fulfilled or not.
func GoVersionMin(min string) (string, bool) {
	current := GoVersion()
	ok := current >= min
	return current, ok
}
