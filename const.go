package plan9

// Plan 9 constants.

const (
	VERSION9P = "9P2000"
	MAXWELEM  = 16
	IOHDRSZ   = 24
	STATMAX   = 65535

	OREAD     = 0
	OWRITE    = 1
	ORDWR     = 2
	OEXEC     = 3
	OTRUNC    = 16
	OCEXEC    = 32
	ORCLOSE   = 64
	ODIRECT   = 128
	ONONBLOCK = 256
	OEXCL     = 0x1000
	OLOCK     = 0x2000
	OAPPEND   = 0x4000
)

const (
	AEXIST = 0 // accessible: exists
	AEXEC  = 1 // execute access
	AWRITE = 2 // write access
	AREAD  = 4 // read access
)

// Qid.Type bits
const (
	QTDIR     = 0x80
	QTAPPEND  = 0x40
	QTEXCL    = 0x20
	QTMOUNT   = 0x10
	QTAUTH    = 0x08
	QTTMP     = 0x04
	QTSYMLINK = 0x02
	QTFILE    = 0x00
)

// Dir.Mode bits
const (
	DMDIR    = 0x80000000 // directories
	DMAPPEND = 0x40000000 // append-only
	DMEXCL   = 0x20000000 // exclusvie use (only one open handle allowed)
	DMMOUNT  = 0x10000000 // mount points
	DMAUTH   = 0x08000000 // authentication file
	DMTMP    = 0x04000000 // non-backed-up files
	DMREAD   = 0x4
	DMWRITE  = 0x2
	DMEXEC   = 0x1
)

// 9P2000.u extensions
const (
	DMSYMLINK   = 0x02000000 // symbolic links
	DMDEVICE    = 0x00800000 // device files
	DMNAMEDPIPE = 0x00200000 // named pipe
	DMSOCKET    = 0x00100000 // socket
	DMSETUID    = 0x00080000 // setuid
	DMSETGID    = 0x00040000 // setgid
	// DMSETVTX = 0x00000000  sticky bit
)

const (
	NOTAG = 0xffff
	NOFID = 0xffffffff
	NOUID = 0xffffffff
)
