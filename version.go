/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							 goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Caesar Disk maker.
 *-----------------------------------------------------------------*/
package caesardisk

import (
	"fmt"
	"runtime"
	"strings"
	//_ "embed"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	MANUAL_VERSION string = "1.1" // in case vcsVersion not injected during link phase

	// Useful Unicode Characters
	CHR_COPYRIGHT       = '\u00a9'      // ©
	CHR_REGISTERED      = '\u00ae'      // ®
	CHR_GUILLEMET_L     = '\u00ab'      // «
	CHR_GUILLEMET_R     = '\u00bb'      // »
	CHR_TRADEMARK       = '\u2122'      // ™
	CHR_SAMARITAN       = '\u214f'      // ⅏
	CHR_PLACEOFINTEREST = '\u2318'      // ⌘
	CHR_HIGHVOLTAGE     = '\u26a1'      // ⚡
	CHR_TRIDENT         = rune(0x1f531) // 🔱
	CHR_SPLATTER        = rune(0x1fadf)
	CHR_WARNING         = '\u26a0' // ⚠
	CHR_EXCLAMATION     = '\u2757'
	CHR_SKULL           = '\u2620' // ☠

	CO1 = "odlamirG omidiD 5202)C("
	CO2 = "stpircS fO droL 5202)C("
	CO3 = "gnitirwnitsol"

	// Change these values accordingly
	NAME string = "Go Caesar Disk Maker"
	DESC string = "An extended Caesar Cipher for modern-day alphabets"
	// don't change
	statusAlpha    status = "Alpha"
	statusBeta     status = "Beta"
	statusRC       status = "RC" // Release Candidate
	statusReleased status = ""
)

var (
	vcsVersion  string // automatically injected with linker
	vcsCommit   string
	vcsDate     string
	vcsBuildNum string
	//NOT USEDgo:embed version.txt
)

var (
	// NOTE: Change these values accordingly
	appVersion version = version{NAME, MANUAL_VERSION, statusReleased, 0}

	// DO NOT CHANGE THESE!
	Version      string = appVersion.String()
	ShortVersion string = appVersion.Short()
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/
type status = string

type version struct {
	n  string // name
	v  string // version tag
	s  status // status
	sv int    // Alpha/Beta/RC-### sequence
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func BuildMeta() string {
	ver := vcsVersion
	if len(vcsVersion) == 0 {
		ver = "v" + MANUAL_VERSION
	}
	return fmt.Sprintf("\t\t%s-%s %s", ver, vcsBuildNum, vcsDate)
}

func (v version) BuildInfo() string {
	return fmt.Sprintf("Build #%s (%s)", vcsBuildNum, vcsCommit)
}

func (v version) Short() string {
	var ver string

	if len(vcsVersion) != 0 {
		v.v = vcsVersion
	}
	var buildInfo string = ""
	if vcsBuildNum != "" {
		buildInfo = fmt.Sprintf(" build #%s", vcsBuildNum)
	}

	switch v.s {
	case statusAlpha:
		fallthrough
	case statusBeta:
		fallthrough
	case statusRC:
		ver = fmt.Sprintf("v%s-%s-%d%s", v.v, v.s, v.sv, buildInfo)
	default:
		ver = fmt.Sprintf("v%s %s", v.v, buildInfo)
	}
	return ver
}

func (v version) String() string {
	var ver string

	if len(vcsVersion) != 0 {
		v.v = vcsVersion
	}
	var buildInfo string = ""
	if vcsBuildNum != "" {
		buildInfo = fmt.Sprintf(" build #%s", vcsBuildNum)
	}

	switch v.s {
	case statusAlpha:
		fallthrough
	case statusBeta:
		fallthrough
	case statusRC:
		ver = fmt.Sprintf("%s v%s-%s-%d%s", v.n, v.v, v.s, v.sv, buildInfo)
	default:
		ver = fmt.Sprintf("%s v%s %s", v.n, v.v, buildInfo)
	}
	return ver
}

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
func BuyMeCoffee(coffee4 ...string) {
	const (
		coffee rune = '\u2615' // ☕
	)

	var recipient string
	if len(coffee4) == 0 {
		recipient = Reverse(CO3)
	} else {
		recipient = coffee4[0]
	}

	fmt.Printf("\t%c Buy me a Coffee? https://www.buymeacoffee/%s\n", coffee, recipient)
}

func Copyright(owner string, withLogo bool) {
	//fmt.Printf("\t\u2720 %s %s \u269d\n", Version, Reverse(owner))
	fmt.Printf("\t%c %s %s %c\n", CHR_TRIDENT, Version, Reverse(owner), CHR_TRIDENT)
	fmt.Println("\t\t\t\t", Logo())
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
