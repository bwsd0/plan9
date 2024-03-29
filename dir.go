package plan9

// Plan 9 directory marshalling. See intro(5).

import (
	"fmt"
	"strconv"
)

type ProtocolError string

func (e ProtocolError) Error() string {
	return string(e)
}

// A Dir contains the metadata for a file.
type Dir struct {
	// system-modified data
	Type uint16 // server type
	Dev  uint32 // server subtype

	// file data
	Qid  Qid    // Unique ID from the server
	Mode uint32 // Permission bits
	// BUG(bwsd): vulnerable to Y2K38
	Atime uint32 // last read time
	Mtime uint32 // last write time

	Length uint64 // file length
	Name   string // last element of path
	Uid    string // owner name
	Gid    string // group name
	Muid   string // last modifier time

	// 9P2000.u extension fields
	// Plan 9 represents user identifiers using strings whereas Unix-like and POSIX
	// environments use numeric identifiers.
	// Uidnum uint
	// Gidnum uint
	// Muidnum uint
	// Ext []byte extended info
}

var nullDir = Dir{
	^uint16(0),
	^uint32(0),
	Qid{^uint64(0), ^uint32(0), ^uint8(0)},
	^uint32(0),
	^uint32(0),
	^uint32(0),
	^uint64(0),
	"",
	"",
	"",
	"",
}

// Null assigns special "don't touch" values to members of d to avoid modifying
// them during plan9.Wstat.
func (d *Dir) Null() {
	*d = nullDir
}

// pdir encodes a 9P stat call on dir d into buffer b.
func pdir(b []byte, d *Dir) []byte {
	n := len(b)
	b = pbit16(b, 0) // length, filled in later
	b = pbit16(b, d.Type)
	b = pbit32(b, d.Dev)
	b = pqid(b, d.Qid)
	b = pbit32(b, uint32(d.Mode))
	b = pbit32(b, d.Atime)
	b = pbit32(b, d.Mtime)
	b = pbit64(b, d.Length)
	b = pstring(b, d.Name)
	b = pstring(b, d.Uid)
	b = pstring(b, d.Gid)
	b = pstring(b, d.Muid)
	pbit16(b[0:n], uint16(len(b)-(n+2)))
	return b
}

func (d *Dir) Bytes() ([]byte, error) {
	return pdir(nil, d), nil
}

// UnmarshalDir decodes a single 9P stat message from b and returns the
// resulting Dir.
//
// If b is too small to hold a valid stat message, ErrShortStat is returned.  If
// the stat message itself is invalid, ErrBadStat is returned.
func UnmarshalDir(b []byte) (d *Dir, err error) {
	defer func() {
		if v := recover(); v != nil {
			d = nil
			err = ProtocolError("malformed Dir")
		}
	}()

	n, b := gbit16(b)
	if int(n) != len(b) {
		panic(1)
	}

	d = new(Dir)
	d.Type, b = gbit16(b)
	d.Dev, b = gbit32(b)
	d.Qid, b = gqid(b)
	d.Mode, b = gbit32(b)
	d.Atime, b = gbit32(b)
	d.Mtime, b = gbit32(b)
	d.Length, b = gbit64(b)
	d.Name, b = gstring(b)
	d.Uid, b = gstring(b)
	d.Gid, b = gstring(b)
	d.Muid, b = gstring(b)

	if len(b) != 0 {
		panic(1)
	}
	return d, nil
}

// String returns the string representation of dir
func (d *Dir) String() string {
	return fmt.Sprintf("'%s' '%s' '%s' '%s' q %v m %#o at %d mt %d l %d t %d d %d",
		d.Name, d.Uid, d.Gid, d.Muid, d.Qid, d.Mode,
		d.Atime, d.Mtime, d.Length, d.Type, d.Dev)
}

// dumpsome returns a string literal representing b quoting unprintable
// characters.
func dumpsome(b []byte) string {
	if len(b) > 64 {
		b = b[0:64]
	}

	// TODO(bwsd): optimize for Latin-1 case using stdlib strconv/quote.go
	printable := true
	for _, c := range b {
		if (c != 0 && c != '\n' && c != '\t' && c < ' ') || c > 127 {
			printable = false
			break
		}
	}

	if printable {
		return strconv.Quote(string(b))
	}
	return fmt.Sprintf("%x", b)
}

type Perm uint32

type permChar struct {
	bit Perm
	c   rune
}

var permChars = []permChar{
	{DMDIR, 'd'},
	{DMAPPEND, 'a'},
	{DMAUTH, 'A'},
	{DMDEVICE, 'D'},
	{DMSOCKET, 'S'},
	{DMNAMEDPIPE, 'P'},
	{0, '-'},
	{DMEXCL, 'l'},
	{DMSYMLINK, 'L'},
	{0, '-'},
	{0400, 'r'},
	{0, '-'},
	{0200, 'w'},
	{0, '-'},
	{0100, 'x'},
	{0, '-'},
	{0040, 'r'},
	{0, '-'},
	{0020, 'w'},
	{0, '-'},
	{0010, 'x'},
	{0, '-'},
	{0004, 'r'},
	{0, '-'},
	{0002, 'w'},
	{0, '-'},
	{0001, 'x'},
	{0, '-'},
}

func (p Perm) String() string {
	s := ""
	did := false
	for _, pc := range permChars {
		if p&pc.bit != 0 {
			did = true
			s += string(pc.c)
		}
		if pc.bit == 0 {
			if !did {
				s += string(pc.c)
			}
			did = false
		}
	}
	return s
}

// Qid represent's a 9P server's unique identification for a file.
type Qid struct {
	Path uint64 // the file server's unique identification for the file
	Vers uint32 // version number for the given path
	Type uint8  // the type of the file (plan9.QDIR for example)
}

func (q Qid) String() string {
	t := ""
	if q.Type&QTDIR != 0 {
		t += "d"
	}
	if q.Type&QTAPPEND != 0 {
		t += "a"
	}
	if q.Type&QTEXCL != 0 {
		t += "l"
	}
	if q.Type&QTAUTH != 0 {
		t += "A"
	}
	return fmt.Sprintf("(%.16x %d %s)", q.Path, q.Vers, t)
}

func gqid(b []byte) (Qid, []byte) {
	var q Qid
	q.Type, b = gbit8(b)
	q.Vers, b = gbit32(b)
	q.Path, b = gbit64(b)
	return q, b
}

func pqid(b []byte, q Qid) []byte {
	b = pbit8(b, q.Type)
	b = pbit32(b, q.Vers)
	b = pbit64(b, q.Path)
	return b
}
