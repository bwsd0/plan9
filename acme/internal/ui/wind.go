package ui

import (
	"bwsd.dev/plan9/acme/internal/adraw"
	"bwsd.dev/plan9/acme/internal/file"
	"bwsd.dev/plan9/acme/internal/wind"

	"bwsd.dev/plan9/draw"
)

func WinresizeAndMouse(w *wind.Window, r draw.Rectangle, safe, keepextra bool) int {
	mouseintag := Mouse.Point.In(w.Tag.All)
	mouseinbody := Mouse.Point.In(w.Body.All)

	y := wind.Winresize(w, r, safe, keepextra)

	// If mouse is in tag, pull up as tag closes.
	if mouseintag && !Mouse.Point.In(w.Tag.All) {
		p := Mouse.Point
		p.Y = w.Tag.All.Max.Y - 3
		adraw.Display.MoveCursor(p)
	}

	// If mouse is in body, push down as tag expands.
	if mouseinbody && Mouse.Point.In(w.Tag.All) {
		p := Mouse.Point
		p.Y = w.Tag.All.Max.Y + 3
		adraw.Display.MoveCursor(p)
	}

	return y
}

func Wintype(w *wind.Window, t *wind.Text, r rune) {
	Texttype(t, r)
	if t.What == wind.Body {
		for i := 0; i < len(t.File.Text); i++ {
			wind.Textscrdraw(t.File.Text[i])
		}
	}
	wind.Winsettag(w)
}

var fff *file.File
