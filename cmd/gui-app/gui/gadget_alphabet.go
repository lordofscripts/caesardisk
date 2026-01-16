/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

// the currently selected alphabet model
type IAlphabetService interface {
	Update()
	// the selected alphabet model
	GetAlphabet() *caesardisk.AlphabetModel
	// length and characters in the alphabet
	GetAlphabetInfo() (int, string)
}

var _ IAlphabetService = (*AlphabetGadget)(nil)
var _ IGadget = (*AlphabetGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type AlphabetGadget struct {
	parent IGadgetParent

	catalogPtr   *AlphabetList
	selected     *caesardisk.AlphabetModel
	previous     int
	defaultKey   string
	langSelector *widget.Select
	langText     *widget.Label
	card         *widget.Card
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

func newAlphabetGadget(parent IGadgetParent, catalog *AlphabetList, defaultKey string) *AlphabetGadget {
	return &AlphabetGadget{
		parent:     parent,
		catalogPtr: catalog,
		defaultKey: defaultKey,
		selected:   (*catalog)[defaultKey],
		previous:   -1,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *AlphabetGadget) Define() IGadget {
	// (W) Alphabet content
	g.langText = widget.NewLabel((*g.catalogPtr)[g.defaultKey].String())

	// (W) Alphabet selector
	// slice of alphabet keys for the selector. This is the same list
	// of names used in AlphabetFactory!
	keys := make([]string, 0, len(*g.catalogPtr))
	for key := range *g.catalogPtr {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create the selector (dropdown)
	//g.langSelector = widget.NewSelect(keys, nil)
	g.langSelector = widget.NewSelectWithData(keys, BoundAlphaName)

	langSelectorLabel := widget.NewLabel("Alphabet")

	selectorContainer := container.New(layout.NewBorderLayout(nil, nil, langSelectorLabel, nil),
		langSelectorLabel,
		g.langSelector,
	)

	alphaCardContent := container.NewVBox(
		selectorContainer,
		g.langText,
	)

	g.card = widget.NewCard(
		"Source Alphabet",
		"for encoding & decoding",
		alphaCardContent,
	)

	return g
}

// Lets the gadget bind data to the widgets
func (g *AlphabetGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *AlphabetGadget) PostRender() IGadget {
	g.langSelector.SetSelected(g.defaultKey)

	g.langSelector.OnChanged = g.onSelectedIndexChanged

	return g
}

func (g *AlphabetGadget) Canvas() fyne.CanvasObject {
	return g.card
}

// Get the container that would be composited into another or a window
func (g *AlphabetGadget) Container() *fyne.Container {
	return nil
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *AlphabetGadget) Update() {
	logx.OnUpdate()

	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	if len(sm.Alpha.Name) != 0 {
		g.langSelector.SetSelected(sm.Alpha.Name)
	} else {
		selectionName := caesardisk.IdentifyAlphabet(&sm.Alpha)
		if len(selectionName) != 0 {
			g.langSelector.SetSelected(selectionName)
		} else {
			logx.Fatalf("unable to identify alphabet '%s'", sm.Alpha.String())
		}
	}
}

// Hide gadget
func (g *AlphabetGadget) Hide() {
	g.card.Hide()
}

// Show gadget
func (g *AlphabetGadget) Show() {
	g.card.Show()
}

// Enable gadget
func (g *AlphabetGadget) Enable() {
	g.langSelector.Enable()
}

// Disable gadget
func (g *AlphabetGadget) Disable() {
	g.langSelector.Disable()
}

// Clears all fields of a gadget
func (g *AlphabetGadget) Clear() {
	g.langSelector.SetSelectedIndex(0)
}

// implement IAlphabetService
func (g *AlphabetGadget) GetAlphabet() *caesardisk.AlphabetModel {
	return g.selected
}

// implement IAlphabetService
func (g *AlphabetGadget) GetAlphabetInfo() (length int, chars string) {
	return g.selected.Length(), g.selected.String()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *AlphabetGadget) onSelectedIndexChanged(selectedKey string) {
	logx.OnChanged(selectedKey)

	value := (*g.catalogPtr)[selectedKey] // Retrieve the alphabet corresponding to the selected key
	// update UI children
	g.langText.SetText(value.String())
	// update Model
	g.selected = value
	if g.previous != g.langSelector.SelectedIndex() {
		g.previous = g.langSelector.SelectedIndex()

		BoundAlphaName.Set(selectedKey)

		g.parent.Cascade(GadgetAlphabet, selectedKey)
	}
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
