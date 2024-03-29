package ui

import (
	"unicode/utf8"

	"bwsd.dev/plan9/acme/internal/adraw"
	"bwsd.dev/plan9/acme/internal/wind"

	"bwsd.dev/plan9/draw"
)

func Rowdragcol(row *wind.Row, c *wind.Column, _0 int) {
	Clearmouse()
	adraw.Display.SwitchCursor2(&adraw.BoxCursor, &adraw.BoxCursor2)
	b := Mouse.Buttons
	op := Mouse.Point
	for Mouse.Buttons == b {
		Mousectl.Read()
	}
	adraw.Display.SwitchCursor(nil)
	if Mouse.Buttons != 0 {
		for Mouse.Buttons != 0 {
			Mousectl.Read()
		}
		return
	}

	wind.Rowdragcol1(row, c, op, Mouse.Point)
	Clearmouse()
	Colmousebut(c)
}

func Rowtype(row *wind.Row, r rune, p draw.Point) *wind.Text {
	if r == 0 {
		r = utf8.RuneError
	}

	Clearmouse()
	BigUnlock()
	row.Lk.Lock()
	BigLock()
	var t *wind.Text
	if Bartflag {
		t = wind.Barttext
	} else {
		t = wind.Rowwhich(row, p)
	}
	if t != nil && (t.What != wind.Tag || !p.In(t.ScrollR)) {
		w := t.W
		if w == nil {
			Texttype(t, r)
		} else {
			wind.Winlock(w, 'K')
			Wintype(w, t, r)
			// Expand tag if necessary
			if t.What == wind.Tag {
				t.W.Tagsafe = false
				if r == '\n' {
					t.W.Tagexpand = true
				}
				WinresizeAndMouse(w, w.R, true, true)
			}
			wind.Winunlock(w)
		}
	}
	row.Lk.Unlock()
	return t
}

var Bartflag bool
