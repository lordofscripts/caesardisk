/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ crypto.ICaesarKeyService = (*CaesarKeyGadget)(nil)
var _ IGadget = (*CaesarKeyGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarKeyGadget struct {
	parent    IGadgetParent
	container *fyne.Container
	labelL    *widget.Label
	labelD    *widget.Label
	slider    *widget.Slider
	//slider   *ScrollableSlider
	alphabet *caesardisk.AlphabetModel
	mutex    sync.Mutex
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/
func init() {}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newCaesarKeyGadget(parent IGadgetParent) *CaesarKeyGadget {
	return &CaesarKeyGadget{
		parent: parent,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *CaesarKeyGadget) Define() IGadget {
	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	// construct gadget
	// · The alphabet's rune corresponding to the Main Key Shift on the slider
	labelAlpha := widget.NewLabel(string(sm.Alpha.FirstChar()))
	// · The alphabet's numerical Main Key Shift taken from the slider's value
	labelShift := widget.NewLabel("00")
	// · The slider in which the user selects the Main Key Shift
	slider := widget.NewSliderWithData(0, float64(sm.Alpha.Length()-1), BoundKeyShift)
	//slider := NewScrollableSliderWidthData(0, float64(sm.Alpha.Length()-1), BoundKeyShift)
	slider.Step = 1
	// · Arrange widgets in a container with a proper layout
	g.container = container.NewBorder(nil, nil, labelAlpha, labelShift, slider)
	g.labelL = labelAlpha
	g.labelD = labelShift
	g.slider = slider

	return g
}

// Lets the gadget bind data to the widgets
func (g *CaesarKeyGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *CaesarKeyGadget) PostRender() IGadget {
	g.slider.OnChanged = g.changeCallback    // updates letter
	g.slider.OnChangeEnded = g.endedCallback // regenerate disk

	return g
}

// Get the container that would be composited into another or a window
func (g *CaesarKeyGadget) Container() *fyne.Container {
	return g.container
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *CaesarKeyGadget) Update() {
	logx.OnUpdate()

	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	g.updateVisualOnly(sm.MainKey.Letter, sm.MainKey.Shift, false)
}

// Hide gadget
func (g *CaesarKeyGadget) Hide() {
	g.container.Hide()
}

// Show gadget
func (g *CaesarKeyGadget) Show() {
	g.container.Show()
}

// Enable slider gadget
func (g *CaesarKeyGadget) Enable() {
	g.slider.Enable()
}

// Disable slider gadget
func (g *CaesarKeyGadget) Disable() {
	g.slider.Disable()
}

// Clears all fields of a gadget
func (g *CaesarKeyGadget) Clear() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// changing the slider value programmatically does NOT trigger OnChange
	g.slider.SetValue(0)
	g.slider.OnChangeEnded(0)
}

// provides/implements ICaesarKeyService
func (g *CaesarKeyGadget) GetKey() (rune, int) {
	//g.mutex.Lock()
	//defer g.mutex.Unlock()

	return g.getKey()
}

// provides/implements ICaesarKeyService
func (g *CaesarKeyGadget) SetKey(letter rune, shift int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.updateVisualOnly(letter, shift, false)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *CaesarKeyGadget) getKey() (rune, int) {
	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	return sm.MainKey.Letter, sm.MainKey.Shift
}

// called every time the slider gets a new step value. Use for
// non-expensive updates. OnChanged is only triggered by user
// interactions, NOT by changes in the data-binding!
func (g *CaesarKeyGadget) changeCallback(sliderValue float64) {
	logx.OnChanged()

	mainKey, mainShift := g.getKey()
	if mainKey != rune(0) {
		g.updateVisualOnly(mainKey, mainShift, true)
	} else {
		logx.Print("getKey failed")
	}
}

// called when the slider gets its final value. Expensive updates
// go here. The global CaesarParams is updated with the new key.
func (g *CaesarKeyGadget) endedCallback(finalSliderValue float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	logx.OnChanged(finalSliderValue)
	logx.Result("key-slide %f", finalSliderValue)

	BoundKeyShift.Set(finalSliderValue)

	mainKey, mainShift := g.getKey()
	g.updateVisualOnly(mainKey, mainShift, true)
	g.parent.Cascade(GadgetMainKey, int(finalSliderValue))
}

func (g *CaesarKeyGadget) updateVisualOnly(let rune, shift int, labelsOnly bool) {
	if !labelsOnly {
		g.slider.Value = float64(shift)
	}
	g.labelD.SetText(fmt.Sprintf("%02d", shift))
	g.labelL.SetText(string(let))
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
