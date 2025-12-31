/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"image"
	"log"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/caesardisk"
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarKeyGadget struct {
	container *fyne.Container
	labelL    *widget.Label
	labelD    *widget.Label
	slider    *widget.Slider
	alphabet  *caesardisk.AlphabetModel
	mutex     sync.Mutex
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newCaesarKeyGadget(alphabet *caesardisk.AlphabetModel) *CaesarKeyGadget {
	// construct gadget
	labelAlpha := widget.NewLabel(string(alphabet.FirstChar()))
	labelShift := widget.NewLabel("00")
	slider := widget.NewSlider(0, float64(alphabet.Length()-1))
	//horizontal := container.New(layout.NewBorderLayout(nil, nil, labelAlpha, nil), labelAlpha, slider)
	horizontal := container.New(layout.NewBorderLayout(nil, nil, labelAlpha, labelShift), labelAlpha, slider, labelShift)

	me := &CaesarKeyGadget{
		labelL:    labelAlpha,
		labelD:    labelShift,
		slider:    slider,
		container: horizontal,
		alphabet:  alphabet,
	}
	slider.OnChanged = me.changeCallback    // updates letter
	slider.OnChangeEnded = me.endedCallback // regenerate disk

	return me
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (g *CaesarKeyGadget) GetContainer() *fyne.Container {
	return g.container
}

func (g *CaesarKeyGadget) SetKey(key int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.slider.SetValue(float64(key))
}

func (g *CaesarKeyGadget) SetKeyChar(key rune) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if index := g.alphabet.Find(key); index != -1 {
		g.slider.SetValue(float64(index))
	} else {
		log.Print("CaesarKeyGadget not a valid key to set")
	}
}

// get the current Caesar key selected in the Gadget
func (g *CaesarKeyGadget) GetKey() (int, rune) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// clean the Display value and get it as rune
	return g.getKey()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *CaesarKeyGadget) getKey() (int, rune) {
	keyRune := ([]rune(strings.Trim(g.labelL.Text, " ")))[0]

	return int(g.slider.Value), keyRune
}

// called every time the slider gets a new step value. Use for
// non-expensive updates
func (g *CaesarKeyGadget) changeCallback(sliderValue float64) {
	caesarKey := int(sliderValue)
	if letter, err := g.alphabet.Character(caesarKey); err == nil {
		g.labelL.SetText(fmt.Sprintf(" %c ", letter))
		g.labelD.SetText(fmt.Sprintf("%02d", caesarKey))
	} else {
		log.Print("CaesarKeyGadget", err)
	}
}

// called when the slider gets its final value. Expensive updates
// go here. The global CaesarParams is updated with the new key.
func (g *CaesarKeyGadget) endedCallback(finalSliderValue float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	const GENERATE_DUAL_ALPHABET_DISK bool = false
	var imgBase, imgOverlay, imgComposite image.Image
	var err error = nil

	keyShift, _ := g.getKey()
	CaesarParams.KeyValue = keyShift // global!

	// base/outer
	if imgBase, err = caesardisk.GenerateCaesarWheelImage(
		CaesarParams.Alphabet.String(), false, *WheelOptions); err == nil {
		// overlay/inner
		if imgOverlay, err = caesardisk.GenerateCaesarWheelImage(
			CaesarParams.Alphabet.String(), true, *WheelOptions); err == nil {
			if imgComposite, err = caesardisk.SuperimposeDisksByShiftImage(
				keyShift,
				CaesarParams.Alphabet.Length(),
				imgBase,
				imgOverlay,
				GENERATE_DUAL_ALPHABET_DISK,
				caesardisk.DefaultCaesarWheelOptions,
			); err == nil {
				gadgetImage.UpdateImage(imgComposite)
			}
		}
	}

	if err != nil {
		log.Print("CaesarKeyGadget ERR", err)
	}

	log.Print("sliderChangeEnded exit")
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
