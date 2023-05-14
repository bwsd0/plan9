package client

import (
	"strings"

	"bwsd.dev/plan9"
)

// Fsys represents a connection to a 9P server.
type Fsys struct {
	root *Fid
}

func (c *Conn) Auth(uname, aname string) (*Fid, error) {
	conn, err := c.conn()
	if err != nil {
		return nil, err
	}
	afidnum, err := conn.newfidnum()
	if err != nil {
		return nil, err
	}
	tx := &plan9.Fcall{Type: plan9.Tauth, Afid: afidnum, Uname: uname, Aname: aname}
	rx, err := conn.rpc(tx, nil)
	if err != nil {
		conn.putfidnum(afidnum)
		return nil, err
	}
	return conn.newFid(afidnum, rx.Qid), nil
}

// Attach establishes a 9P fileserver connection for a given user.
//
// The afid argument specifies a fid to reuse from a previous auth message.  To
// connect without authentication, the afid field should be set to NOFID.
func (c *Conn) Attach(afid *Fid, user, aname string) (*Fsys, error) {
	conn, err := c.conn()
	if err != nil {
		return nil, err
	}
	fidnum, err := conn.newfidnum()
	if err != nil {
		return nil, err
	}
	tx := &plan9.Fcall{Type: plan9.Tattach, Afid: plan9.NOFID, Fid: fidnum, Uname: user, Aname: aname}
	if afid != nil {
		tx.Afid = afid.fid
	}
	rx, err := conn.rpc(tx, nil)
	if err != nil {
		conn.putfidnum(fidnum)
		return nil, err
	}
	return &Fsys{conn.newFid(fidnum, rx.Qid)}, nil
}

var accessOmode = [8]uint8{
	0,
	plan9.OEXEC,
	plan9.OWRITE,
	plan9.ORDWR,
	plan9.OREAD,
	plan9.OEXEC, // only approximate
	plan9.ORDWR,
	plan9.ORDWR, // only approximate
}

func (fs *Fsys) Access(name string, mode int) error {
	if mode == plan9.AEXIST {
		_, err := fs.Stat(name)
		return err
	}
	fid, err := fs.Open(name, accessOmode[mode&7])
	if fid != nil {
		fid.Close()
	}
	return err
}

// Create creates a new file or prepares to rewrite an existinf file.
//
// The filename is opened according to omode (as described for open), and
// returns the an associated file descriptor. If the file is new, the owner is
// set to the userid of the creating process group; the group to that of the
// containing directory; the permissions to perm ANDed with the permissions of
// the containing directory.
//
// If the file already exists, it is truncated to 0 length, and the permissions,
// owner, and group remain unchanged. The created file is a directory if the
// DMDIR bit is set in perm, an exclusive-use file if the DMEXCL bit is set, and
// append-only file if the DMAPPEND bit is set.
//
// Exclusive-use files may be open for I/O by only one client at a time, but the
// file descriptor may become invalid if no I/O is done for an extended period.
//
// Create fails if the path up to the last element of a file cannot be
// evaluated, if the user doesn't have write permission in the final directory,
// if the file already exists and does not permit the access defined by omode,
// or if there are no free file descriptors. In the last case, the file may be
// created even when an error is returned.
//
// Since create may succeed even if the file exists, a special mechanism is
// necessary for those applications that require an atomic create operation. If
// th OEXCL (0x1000) bit is set in the mode for a create, the call succeeds only
// if the file does not already exist.
//
// See: open(9p)
func (fs *Fsys) Create(name string, mode uint8, perm uint32) (*Fid, error) {
	i := strings.LastIndex(name, "/")
	var dir, elem string
	if i < 0 {
		elem = name
	} else {
		dir, elem = name[0:i], name[i+1:]
	}
	fid, err := fs.root.Walk(dir)
	if err != nil {
		return nil, err
	}
	err = fid.Create(elem, mode, perm)
	if err != nil {
		fid.Close()
		return nil, err
	}
	return fid, nil
}

// Open opens the file for I/O and returns an associated file descriptor.
//
// Omode is one of OREAD, OWRITE, ORDWR or OEXEC, asking for permssion to read,
// write, read and write, or execute, respecitvely. In addition, there are three
// values that can be ORed with omode: OTRUN says to truncate the file to zero
// length before opening it; OCEXEC says to close the file when an exec(3) or
// execl system is made; ORCLOSE says to remove the file when it is closed (by
// everyone who has a copy of the file descriptor); and OAPPEND says to open the
// file in append-only mode, so that writes are always appened to the end of the
// file.
//
// Open fails if the file does not exist or the user does not have permission to
// open it for the requested purpose. See: stat(3). The user must have write
// permission on the file if the OTRUNC bit is set.  For the open system call,
// unlike the implicit open in exec(3), OEXEC is actually identical to OREAD.
func (fs *Fsys) Open(name string, mode uint8) (*Fid, error) {
	fid, err := fs.root.Walk(name)
	if err != nil {
		return nil, err
	}
	if err := fid.Open(mode); err != nil {
		fid.Close()
		return nil, err
	}
	return fid, nil
}

func (fs *Fsys) Remove(name string) error {
	fid, err := fs.root.Walk(name)
	if err != nil {
		return err
	}
	return fid.Remove()
}

func (fs *Fsys) Stat(name string) (*plan9.Dir, error) {
	fid, err := fs.root.Walk(name)
	if err != nil {
		return nil, err
	}
	d, err := fid.Stat()
	fid.Close()
	return d, err
}

func (fs *Fsys) Wstat(name string, d *plan9.Dir) error {
	fid, err := fs.root.Walk(name)
	if err != nil {
		return err
	}
	err = fid.Wstat(d)
	fid.Close()
	return err
}

// Close closes the file associated with a file descriptor.
//
// Provide the file descriptor is a valid open descriptor, close is guaranteed
// to close it; there will be no error. Files are closed automatically upon
// termination of a process; close allows the file descriptor to be reused.
func (fs *Fsys) Close() error {
	return fs.root.Close()
}
