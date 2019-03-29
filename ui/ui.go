package ui

import (
	"github.com/jroimartin/gocui"
)

type View interface {
	SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error)
	CursorDown(g *gocui.Gui, v *gocui.View) error
	CursorUp(g *gocui.Gui, v *gocui.View) error
	NextView(g *gocui.Gui, v *gocui.View)
}

func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		nextLine, err := v.Line(cy + 1)
		if err != nil {
			return nil
		}
		if nextLine == "" {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	viewArr = []string{"Image", "Container"}
	active  = 0
)

func NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]
	_, err := SetCurrentViewOnTop(g, name)
	if err != nil {
		return err
	}
	next, err := g.View(name)
	if err != nil {
		return err
	}

	if next.SelBgColor == 0 {
		next.SelBgColor = gocui.ColorGreen
		next.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorDefault
		v.SelFgColor = gocui.ColorWhite
	}
	active = nextIndex
	return nil
}
