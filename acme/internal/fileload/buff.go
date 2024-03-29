package fileload

import (
	"io"
	"os"
	"unicode/utf8"

	"bwsd.dev/plan9/acme/internal/alog"
	"bwsd.dev/plan9/acme/internal/bufs"
	"bwsd.dev/plan9/acme/internal/runes"
	"bwsd.dev/plan9/acme/internal/util"
	"bwsd.dev/plan9/acme/internal/wind"
)

func Loadfile(fd *os.File, q0 int, nulls *bool, f func(int, []rune) int, h io.Writer) int {
	p := make([]byte, bufs.Len+utf8.UTFMax+1)
	r := make([]rune, bufs.Len)
	m := 0
	n := 1
	q1 := q0
	/*
	 * At top of loop, may have m bytes left over from
	 * last pass, possibly representing a partial rune.
	 */
	for n > 0 {
		var err error
		n, err = fd.Read(p[m : m+bufs.Len])
		if err != nil && err != io.EOF {
			alog.Printf("read error in Buffer.load: %v", err)
			break
		}
		if h != nil {
			h.Write(p[m : m+n])
		}
		m += n
		nb, nr, nulls1 := runes.Convert(p[:m], r, err == io.EOF)
		if nulls1 {
			*nulls = true
		}
		copy(p, p[nb:m])
		m -= nb
		q1 += f(q1, r[:nr])
	}
	return q1 - q0
}

func fileloader(f *wind.File) func(pos int, data []rune) int {
	return func(pos int, data []rune) int {
		f.Insert(pos, data)
		return len(data)
	}
}

func fileload1(f *wind.File, pos int, fd *os.File, nulls *bool, h io.Writer) int {
	if pos > f.Len() {
		util.Fatal("internal error: fileload1")
	}
	return Loadfile(fd, pos, nulls, fileloader(f), h)
}
