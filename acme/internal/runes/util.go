// #include "dat.h"
// #include "fns.h"

package runes

type Text interface {
	Len() int
	RuneAt(pos int) rune
}
