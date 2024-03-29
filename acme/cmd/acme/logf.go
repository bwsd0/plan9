// State for global log file.

package main

import (
	"fmt"
	"sync"

	"bwsd.dev/plan9/acme/internal/wind"

	"bwsd.dev/plan9"
)

type Log struct {
	lk    sync.Mutex
	r     sync.Cond
	start int64
	ev    []string
	f     []*Fid
	read  []*Xfid
}

var eventlog Log

func init() {
	// Using eventlog.lk with a sync.Cond means we have
	// to reacquire big after sync.Wait has locked eventlog.lk.
	// Therefore every acquisition of eventlog.lk must drop big
	// before eventlog.lk.Lock and then reacquire it afterward,
	// or else the two different lock orders will deadlock.
	eventlog.r.L = &eventlog.lk
}

func xfidlogopen(x *Xfid) {
	bigUnlock()
	eventlog.lk.Lock()
	bigLock()
	eventlog.f = append(eventlog.f, x.f)
	x.f.logoff = eventlog.start + int64(len(eventlog.ev))
	eventlog.lk.Unlock()
}

func xfidlogclose(x *Xfid) {
	bigUnlock()
	eventlog.lk.Lock()
	bigLock()
	for i := 0; i < len(eventlog.f); i++ {
		if eventlog.f[i] == x.f {
			eventlog.f[i] = eventlog.f[len(eventlog.f)-1]
			eventlog.f = eventlog.f[:len(eventlog.f)-1]
			break
		}
	}
	eventlog.lk.Unlock()
}

func xfidlogread(x *Xfid) {
	bigUnlock()
	eventlog.lk.Lock()
	bigLock()
	eventlog.read = append(eventlog.read, x)

	x.flushed = false
	for x.f.logoff >= eventlog.start+int64(len(eventlog.ev)) && !x.flushed {
		bigUnlock()
		eventlog.r.Wait()
		bigLock()
	}
	var i int

	for i = 0; i < len(eventlog.read); i++ {
		if eventlog.read[i] == x {
			eventlog.read[i] = eventlog.read[len(eventlog.read)-1]
			eventlog.read = eventlog.read[:len(eventlog.read)-1]
			break
		}
	}

	if x.flushed {
		eventlog.lk.Unlock()
		return
	}

	i = int(x.f.logoff - eventlog.start)
	p := eventlog.ev[i]
	x.f.logoff++
	eventlog.lk.Unlock()

	var fc plan9.Fcall
	fc.Data = []byte(p)
	fc.Count = uint32(len(fc.Data))
	respond(x, &fc, "")
}

func xfidlogflush(x *Xfid) {
	bigUnlock()
	eventlog.lk.Lock()
	bigLock()
	for i := 0; i < len(eventlog.read); i++ {
		rx := eventlog.read[i]
		if rx.fcall.Tag == x.fcall.Oldtag {
			rx.flushed = true
			eventlog.r.Broadcast()
		}
	}
	eventlog.lk.Unlock()
}

/*
 * add a log entry for op on w.
 * expected calls:
 *
 * op == "new" for each new window
 *	- caller of coladd or makenewwindow responsible for calling
 *		xfidlog after setting window name
 *	- exception: zerox
 *
 * op == "zerox" for new window created via zerox
 *	- called from zeroxx
 *
 * op == "get" for Get executed on window
 *	- called from get
 *
 * op == "put" for Put executed on window
 *	- called from put
 *
 * op == "del" for deleted window
 *	- called from winclose
 */
func xfidlog(w *wind.Window, op string) {
	bigUnlock()
	eventlog.lk.Lock()
	bigLock()
	if len(eventlog.ev) >= cap(eventlog.ev) {
		// Remove and free any entries that all readers have read.
		min := eventlog.start + int64(len(eventlog.ev))
		for i := 0; i < len(eventlog.f); i++ {
			if min > eventlog.f[i].logoff {
				min = eventlog.f[i].logoff
			}
		}
		if min > eventlog.start {
			n := int(min - eventlog.start)
			copy(eventlog.ev, eventlog.ev[n:])
			eventlog.ev = eventlog.ev[:len(eventlog.ev)-n]
			eventlog.start += int64(n)
		}

		// Otherwise grow (in append below).
	}

	f := w.Body.File
	name := string(f.Name())
	eventlog.ev = append(eventlog.ev, fmt.Sprintf("%d %s %s\n", w.ID, op, name))
	eventlog.r.Broadcast()
	eventlog.lk.Unlock()
}
