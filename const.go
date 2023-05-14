package plan9

// Plan 9 constants.

// TODO: document constants
const (
	VERSION9P = "9P2000"
	MAXWELEM  = 16 // Maximum number of elements or Qids allowed in a single message
	IOHDRSZ   = 24 // Buffer size to reserve for a 9P header
	STATMAX   = (1 << 16) - 1
)

const (
	OREAD     = 0      // open for Read
	OWRITE    = 1      // write
	ORDWR     = 2      // read and write
	OEXEC     = 3      // execute, == read but check execute permission
	OTRUNC    = 16     // truncate file first (except for exec)
	OCEXEC    = 32     // close on exec
	ORCLOSE   = 64     // remove on close
	ODIRECT   = 128    // direct access
	ONONBLOCK = 256    // non-blocking call
	OEXCL     = 0x1000 // exclusive use (create only)
	OLOCK     = 0x2000 // lock after opening
	OAPPEND   = 0x4000 // append only
)

const (
	AEXIST = 0 // accessible: exists
	AEXEC  = 1 // execute access
	AWRITE = 2 // write access
	AREAD  = 4 // read access
)

// Qid.Type bits
const (
	QTDIR    = 0x80 // directories
	QTAPPEND = 0x40 // append only files
	QTEXCL   = 0x20 // exclusive use files
	QTMOUNT  = 0x10 // mounted channel
	QTAUTH   = 0x08 // authentication file
	QTTMP    = 0x04 // non-backed-up file
	QTFILE   = 0x00 // plain file
)

// Dir.Mode bits
const (
	DMDIR    = 0x80000000 // directories
	DMAPPEND = 0x40000000 // append-only
	DMEXCL   = 0x20000000 // exclusive use (only one open handle allowed)
	DMMOUNT  = 0x10000000 // mount points
	DMAUTH   = 0x08000000 // authentication file
	DMTMP    = 0x04000000 // non-backed-up files
	DMREAD   = 0x4        // read permission
	DMWRITE  = 0x2        // write permission
	DMEXEC   = 0x1        // execute permission
)

// 9P2000.u extensions
const (
	DMSYMLINK   = 0x02000000 // symbolic links
	DMLINK      = 0x01000000 // hard link
	DMDEVICE    = 0x00800000 // device files
	DMNAMEDPIPE = 0x00200000 // named pipe
	DMSOCKET    = 0x00100000 // socket
	DMSETUID    = 0x00080000 // setuid
	DMSETGID    = 0x00040000 // setgid
	// Unimplemented
	// DMSETVTX = 0x00000000 // sticky bit
)

const (
	NOTAG = 0xffff     // Dummy tag
	NOFID = 0xffffffff // Dummy fid
	NOUID = 0xffffffff // Dummy uid
)
