/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import (
	"image"
	"image/color"
	"log"

	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ fyne.Widget = (*ClickableImage)(nil)
var _ fyne.Tappable = (*ClickableImage)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type ClickableImage struct {
	widget.BaseWidget
	image *canvas.Image
	url   *url.URL
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewClickableImage(res fyne.Resource, targetURL string, size fyne.Size) *ClickableImage {
	u, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
	}

	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(size)

	c := &ClickableImage{
		image: img,
		url:   u,
	}

	c.ExtendBaseWidget(c)
	return c
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// CreateRenderer implements the fyne.Widget interface.
func (c *ClickableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.image)
}

// Tapped implements the fyne.Tappable interface to handle click events.
func (c *ClickableImage) Tapped(_ *fyne.PointEvent) {
	if c.url != nil {
		fyne.CurrentApp().OpenURL(c.url)
	}
}

func (c *ClickableImage) SetImage(img *canvas.Image) {
	c.image = img
	c.Refresh()
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// returns a minimal transparent pixel that can be stretched
func TransparentImage() *canvas.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0, 0, 0, 0}) // transparent pixel
	fyneImg := canvas.NewImageFromImage(img)
	fyneImg.FillMode = canvas.ImageFillStretch
	return fyneImg
}

// returns a theme's icon scaled to the given size
func IconImage(res fyne.Resource, size fyne.Size) *canvas.Image {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain // keep aspect
	img.SetMinSize(size)
	return img
}
