//go:build exclude

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * An extension to widget.Slider that implements fyne.Scrollable to
 * capture the mousewheel movement (forward/backward) and use its DY
 * to increase/decrease the slider value accordingly.
 * BROKEN: odd behavior, the slider ball doesn't quite follow the value.
 *-----------------------------------------------------------------*/
package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type ScrollableSlider struct {
	*widget.Slider

	OnScrolled func(*fyne.ScrollEvent)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewScrollableSlider(min, max float64) *ScrollableSlider {
	ss := &ScrollableSlider{
		Slider: widget.NewSlider(min, max),
	}

	return ss
}

func NewScrollableSliderWidthData(min, max float64, data binding.Float) *ScrollableSlider {
	ss := &ScrollableSlider{
		Slider: widget.NewSliderWithData(min, max, data),
	}

	return ss
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fyne.Scrollable
func (g *ScrollableSlider) Scrolled(event *fyne.ScrollEvent) {
	// @note FWD scroll gives Dx:0 Dy:25 and BACK scroll Dx:0 Dy:-25
	// I wonder if that 25 is the same for all platforms or does it vary?
	const DIVIDER int = 25
	// dx := int(float64(event.Scrolled.DX))
	dy := int(float64(event.Scrolled.DY)) / DIVIDER
	//fmt.Printf("Slider mousewheel dX:%d dY:%d\n", dx, dy)

	newValue := g.Value + float64(dy)
	if newValue > g.Max {
		newValue = g.Max
	} else if newValue < g.Min {
		newValue = g.Min
	}
	fmt.Printf("Slider mousewheel dY:%d new:%f\n", dy, newValue)
	g.SetValue(newValue)
	g.Refresh()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
