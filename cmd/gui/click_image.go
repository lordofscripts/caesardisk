/* -----------------------------------------------------------------
 *
 *				  Copyright (C)2025
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type ClickImage struct {
	widget.BaseWidget
	img            *canvas.Image
	box            *fyne.Container
	OnTapped       func(evt *fyne.PointEvent)
	OnDoubleTapped func(evt *fyne.PointEvent)
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewClickImageFromImage(img *canvas.Image) *ClickImage {
	w := &ClickImage{}
	w.img = img
	w.box = container.NewStack(img)
	w.ExtendBaseWidget(w)
	return w
}

func NewClickImageFromResource(resource fyne.Resource) *ClickImage {
	w := &ClickImage{}
	w.img = canvas.NewImageFromResource(resource)
	w.box = container.NewStack(w.img)
	w.ExtendBaseWidget(w)
	return w
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (c *ClickImage) SetMinSize(size fyne.Size) {
	c.img.SetMinSize(size)
	c.img.Refresh()
}

func (c *ClickImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.box)
}

func (c *ClickImage) Tapped(evt *fyne.PointEvent) {
	if c.OnTapped != nil {
		c.OnTapped(evt)
	}
}

func (c *ClickImage) TappedSecondary(evt *fyne.PointEvent) {
	if c.OnDoubleTapped != nil {
		c.OnDoubleTapped(evt)
	}
}

func (c *ClickImage) SetImage(img *canvas.Image) {
	c.box.RemoveAll()
	c.box.Add(img)
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
