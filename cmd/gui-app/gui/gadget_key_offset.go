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
	"fyne.io/fyne/v2/theme"
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

var _ crypto.IKeyOffsetService = (*KeyOffsetGadget)(nil)
var _ IGadget = (*KeyOffsetGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type KeyOffsetGadget struct {
	parent    IGadgetParent
	container *fyne.Container
	labelL    *widget.Label
	labelD    *widget.Label
	slider    *widget.Slider
	alphabet  *caesardisk.AlphabetModel
	mutex     sync.Mutex
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newKeyOffsetGadget(parent IGadgetParent) *KeyOffsetGadget {
	return &KeyOffsetGadget{
		parent: parent,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *KeyOffsetGadget) Define() IGadget {
	// construct gadget
	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	// · Slider title
	labelTitle := NewDynamicLabelWithStyle("Alternate Key Offset", fyne.TextAlignCenter, fyne.TextStyle{
		Bold:      true,
		Italic:    true,
		Monospace: false,
	}, nil)
	labelTitle.SizeName = theme.SizeNameCaptionText
	labelTitle.Wrapping = fyne.TextWrapOff
	// · The alphabet's rune corresponding to the COMPUTED selected key offset
	labelAlpha := widget.NewLabel(string(sm.Alpha.FirstChar()))
	// · The alphabet's key shift offset (it is applied over main key shift)
	labelShift := widget.NewLabel("00")
	// · Slider in which the user selects the key shift Offset
	slider := widget.NewSliderWithData(0,
		float64(sm.Alpha.Length()-1),
		BoundKeyOffset)
	slider.Step = 1
	// · container for all widgets
	g.container = container.NewBorder(labelTitle, nil, labelAlpha, labelShift, slider)

	g.labelL = labelAlpha
	g.labelD = labelShift
	g.slider = slider

	return g
}

// Lets the gadget bind data to the widgets
func (g *KeyOffsetGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *KeyOffsetGadget) PostRender() IGadget {
	g.slider.OnChanged = g.changeCallback    // updates letter
	g.slider.OnChangeEnded = g.endedCallback // regenerate disk

	g.container.Hide()
	return g
}

// Get the container that would be composited into another or a window
func (g *KeyOffsetGadget) Container() *fyne.Container {
	return g.container
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *KeyOffsetGadget) Update() {
	logx.OnUpdate()

	g.mutex.Lock()
	defer g.mutex.Unlock()

	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()
	g.updateVisualOnly(sm.AltKey.Letter, sm.Offset)
}

// Hide gadget
func (g *KeyOffsetGadget) Hide() {
	logx.Visit()
	g.container.Hide()
}

// Show gadget
func (g *KeyOffsetGadget) Show() {
	logx.Visit()
	g.container.Show()
}

// Enable slider gadget
func (g *KeyOffsetGadget) Enable() {
	g.slider.Enable()
}

// Disable slider gadget
func (g *KeyOffsetGadget) Disable() {
	g.slider.Disable()
}

// Clears all fields of a gadget
func (g *KeyOffsetGadget) Clear() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// changing the slider value programmatically does NOT trigger OnChange
	g.slider.SetValue(0)
	g.slider.OnChangeEnded(0)
}

// provides/implements ICaesarKeyService
func (g *KeyOffsetGadget) GetOffset() int {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.getOffset()
}

// get the derived alternate key based on the corrected offset
func (g *KeyOffsetGadget) GetAltKey() (rune, int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()
	return sm.AltKey.Letter, sm.AltKey.Shift
}

func (g *KeyOffsetGadget) SetOffsetRune(altKeyLetter rune) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.labelL.SetText(string(altKeyLetter))
}

func (g *KeyOffsetGadget) SetOffset(altKeyLetter rune, altShift int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.labelL.SetText(string(altKeyLetter))
	g.slider.Value = float64(altShift)
}

// resets the control's view to a slider value of zero
func (g *KeyOffsetGadget) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.labelL.SetText("  ")
	g.labelD.SetText("00")
	g.slider.Value = 0
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *KeyOffsetGadget) getOffset() int {
	shiftOffset := int(g.slider.Value)

	return shiftOffset
}

// called every time the slider gets a new step value. Use for
// non-expensive updates
func (g *KeyOffsetGadget) changeCallback(sliderValue float64) {
	logx.OnChanged()

	shiftOffset := g.getOffset()
	g.labelD.SetText(fmt.Sprintf("%02d", shiftOffset))
}

// called when the slider gets its final value. Expensive updates
// go here. The global CaesarParams is updated with the new key.
func (g *KeyOffsetGadget) endedCallback(finalSliderValue float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	logx.OnChanged(finalSliderValue)

	// @todo refactored use also updateVisualOnly
	// Update the selected offset value in the composite widget
	BoundKeyOffset.Set(g.slider.Value)
	g.labelD.SetText(fmt.Sprintf("%02d", int(g.slider.Value)))

	var ak crypto.CaesarKey
	sm := DataBindings.GetSessionModel()
	ak = sm.AltKey

	// Display the Alternate Key letter on the composite.
	// this key corresponds to the corrected/recomputed offset
	// on the selected alphabet
	g.labelL.SetText(string(ak.Letter))

	// relay value to parameter container and let it give us
	// the corresponding (recalculated) alternate key shift value.
	g.parent.Cascade(GadgetOffset, ak.Shift)
}

func (g *KeyOffsetGadget) updateVisualOnly(let rune, shift int) {
	g.slider.Value = float64(shift)
	g.labelD.SetText(fmt.Sprintf("%02d", shift))
	g.labelL.SetText(string(let))
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
