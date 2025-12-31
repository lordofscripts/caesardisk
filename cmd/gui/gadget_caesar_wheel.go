/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							goCaesarDisk GUI
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * The Caesar Wheel gadget encapsulates an image object in which a
 * Caesar cipher wheel is displayed. It is updated as the key is
 * changed, thus allowing the user to visually try to encode/decode.
 *-----------------------------------------------------------------*/
package main

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	wheel_WIDTH  = 300
	wheel_HEIGHT = 300
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarWheelGadget struct {
	image *canvas.Image
	mutex sync.Mutex
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newCaesarWheelGadget(initialFile string) *CaesarWheelGadget {
	me := &CaesarWheelGadget{}
	me.reload(initialFile)
	return me
}

func newCaesarWheelGadgetFrom(imagen image.Image) *CaesarWheelGadget {
	displaySize := fyne.NewSize(wheel_WIDTH, wheel_HEIGHT)

	img := canvas.NewImageFromImage(imagen)
	img.SetMinSize(displaySize)
	img.Resize(displaySize)
	img.FillMode = canvas.ImageFillContain

	return &CaesarWheelGadget{
		image: img,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (g *CaesarWheelGadget) GetImage() *canvas.Image {
	return g.image
}

func (g *CaesarWheelGadget) UpdateImageFromFile(filename string) {
	g.reload(filename)
}

func (g *CaesarWheelGadget) UpdateImage(img image.Image) {
	g.image.Image = img
	g.image.Refresh()
}

func (g *CaesarWheelGadget) UpdateImageForKey(key int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// @todo caesardisk. generate new image object reload without file
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *CaesarWheelGadget) reload(filename string) *canvas.Image {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	displaySize := fyne.NewSize(wheel_WIDTH, wheel_HEIGHT)

	img := canvas.NewImageFromFile(filename)
	img.SetMinSize(displaySize)
	img.Resize(displaySize)
	img.FillMode = canvas.ImageFillContain

	g.image = nil
	g.image = img

	return img
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
