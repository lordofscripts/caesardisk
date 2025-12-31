/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *					 	goCaesarDisk GUI
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Graphical User Interface with the Fyne toolkit that allows for
 * simple Caesar cipher encoding/decoding. It uses the same library
 * as the CaesarDisk CLI to generate the wheel images.
 *-----------------------------------------------------------------*/
package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/caesardisk/internal/cipher"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	WINDOW_WIDTH  = 380
	WINDOW_HEIGHT = 600
)

var (
	gadgetCaesarKey *CaesarKeyGadget   = nil
	gadgetImage     *CaesarWheelGadget = nil

	optionUsePDU binding.Bool
)

//go:embed caesar_disk_sample.png
var caesarDiskSample embed.FS

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type SliderCallback func(f float64)

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
func loadEmbeddedImage(path string) (image.Image, error) {
	data, err := caesarDiskSample.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func guiBuildMainTab() *container.TabItem {
	const TAB_NAME = "Caesar"

	if sample, err := loadEmbeddedImage("caesar_disk_sample.png"); err != nil {
		log.Print(err)
		return nil
	} else {
		gadgetImage = newCaesarWheelGadgetFrom(sample)
	}

	gadgetCaesarKey = newCaesarKeyGadget(DefaultAlphabet)

	// Multi-line text widget
	textEntry1 := widget.NewMultiLineEntry()
	textEntry1.SetPlaceHolder("Enter plain text here...")
	textEntry1.Wrapping = fyne.TextWrapBreak

	// Second multi-line text widget
	textEntry2 := widget.NewMultiLineEntry()
	textEntry2.SetPlaceHolder("Output will appear here...")
	textEntry2.Wrapping = fyne.TextWrapBreak

	// Action button
	actionExchangeButton := widget.NewButton(string(rune(0x21d5)), func() {
		temp := textEntry1.Text
		textEntry1.SetText(textEntry2.Text)
		textEntry2.SetText(temp)
	})

	actionClearButton := widget.NewButton("CLR", func() {
		textEntry1.SetText("")
		textEntry2.SetText("")
	})

	actionEncodeButton := widget.NewButton("Encode", func() {
		var encoded string
		var err error = nil
		engine := cipher.NewCaesarCipher(CaesarParams)
		if usePDU, _ := optionUsePDU.Get(); usePDU {
			encoded, err = engine.EncodeMessage(textEntry1.Text)
			if err != nil {
				log.Print(err)
			}
		} else {
			encoded = engine.Encode(textEntry1.Text)
		}

		textEntry2.SetText(encoded)
	})
	actionEncodeButton.Disable()

	actionDecodeButton := widget.NewButton("Decode", func() {
		var decoded string
		var err error = nil
		engine := cipher.NewCaesarCipher(CaesarParams)
		if usePDU, _ := optionUsePDU.Get(); usePDU {
			decoded, err = engine.DecodeMessage(textEntry1.Text)
			if err != nil {
				log.Print(err)
			}
		} else {
			decoded = engine.Decode(textEntry1.Text)
		}

		textEntry2.SetText(decoded)
	})
	actionDecodeButton.Disable()

	actionButtonContainer := container.NewHBox(actionEncodeButton, actionExchangeButton, actionClearButton, actionDecodeButton)

	textEntry1.OnChanged = func(s string) {
		if len(s) > 0 {
			actionEncodeButton.Enable()
			actionDecodeButton.Enable()
		} else {
			actionEncodeButton.Disable()
			actionDecodeButton.Disable()
		}
	}

	// Organize first tab content
	tab1Content := container.New(layout.NewVBoxLayout(),
		//img,
		gadgetImage.GetImage(),
		gadgetCaesarKey.GetContainer(),
		textEntry1,
		actionButtonContainer,
		textEntry2,
	)

	tab1 := container.NewTabItem(TAB_NAME, tab1Content)

	return tab1
}

func guiBuildParamTab() *container.TabItem {
	const TAB_NAME = "Options"

	// · Alphabet content
	langText := widget.NewLabel(DefaultAlphabet.String())

	// · Alphabet selector
	// slice of alphabet keys for the selector
	keys := make([]string, 0, len(MyAlphabets))
	for key := range MyAlphabets {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create the selector (dropdown)
	langSelector := widget.NewSelect(keys, func(selectedKey string) {
		value := MyAlphabets[selectedKey] // Retrieve the alphabet corresponding to the selected key
		// update UI children
		langText.SetText(value.String())
		// update Model
		CaesarParams.Alphabet = value // @todo MUX Lock
	})
	langSelector.SetSelected(defaultAlphabetKey)
	langSelectorLabel := widget.NewLabel("Alphabet")

	selectorContainer := container.New(layout.NewBorderLayout(nil, nil, langSelectorLabel, nil),
		langSelectorLabel,
		langSelector,
	)

	alphaCardContent := container.NewVBox(
		selectorContainer,
		langText,
	)
	alphaCard := widget.NewCard(
		"Source Alphabet",
		"for encoding & decoding",
		alphaCardContent,
	)

	// · Rendering: Orthogonality
	var checkOrtho *widget.Check
	checkOrtho = widget.NewCheck("Orthogonal", nil)
	checkOrtho.SetChecked(WheelOptions.Orthogonal) // we don't want to trigger just yet
	checkOrtho.OnChanged = func(b bool) {
		WheelOptions.Orthogonal = checkOrtho.Checked
	}

	optionUsePDU = binding.NewBool() // data source for PDU checkbox
	optionUsePDU.Set(false)
	checkPDU := widget.NewCheckWithData("Use PDU format", optionUsePDU)

	// -- Content Layout
	tab2Content := container.New(layout.NewVBoxLayout(),
		alphaCard,
		checkOrtho,
		checkPDU,
	)

	// -- Tab
	tab2 := container.NewTabItem(TAB_NAME, tab2Content)

	return tab2
}

// Builds the main application Graphical User Interface
func BuildGUI() *fyne.Window {
	const APP_TITLE = "Modern Caesar Cipher"

	myApp := app.New()
	myWindow := myApp.NewWindow(APP_TITLE)

	// Create main menu
	shortcutQuit := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierAlt,
	}
	menuQuit := fyne.NewMenuItem("Quit", func() {
		myApp.Quit()
	})
	menuQuit.Shortcut = shortcutQuit

	menuTopFile := fyne.NewMenu("File", menuQuit)
	menuMain := fyne.NewMainMenu(menuTopFile)
	myWindow.SetMainMenu(menuMain)

	// Create tabs
	tab1 := guiBuildMainTab()
	tab2 := guiBuildParamTab()
	tabs := container.NewAppTabs(
		tab1,
		tab2,
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT)) // Adjust window size
	myWindow.SetFixedSize(false)

	return &myWindow
}
