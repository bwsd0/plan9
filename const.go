package plan9

// Plan 9 constants.

const (
	VERSION9P = "9P2000"
	MAXWELEM  = 16
	IOHDRSZ   = 24
	STATMAX   = 65535
)

const (
	OREAD     = 0      // Open for Read
	OWRITE    = 1      // Write
	ORDWR     = 2      // Read and write
	OEXEC     = 3      // Execute, == read but check execute permission
	OTRUNC    = 16     // Truncate file first (except for exec)
	OCEXEC    = 32     // Close on exec
	ORCLOSE   = 64     // Remove on close
	ODIRECT   = 128    // Direct access
	ONONBLOCK = 256    // Non-blocking call
	OEXCL     = 0x1000 // Exclusive use (create only)
	OLOCK     = 0x2000 // Lock after opening
	OAPPEND   = 0x4000 // Append only
)

const (
	AEXIST = 0 // accessible: exists
	AEXEC  = 1 // execute access
	AWRITE = 2 // write access
	AREAD  = 4 // read access
)

// Qid.Type bits
const (
	QTDIR     = 0x80 // type bit for directories
	QTAPPEND  = 0x40 // type bit for append only files
	QTEXCL    = 0x20 // type bit for exclusive use files
	QTMOUNT   = 0x10 // type bit for mounted channel
	QTAUTH    = 0x08 // type bit for authentication file
	QTTMP     = 0x04 // type bit for non-backed-up file
	QTSYMLINK = 0x02 // type bit for symbolic link
	QTFILE    = 0x00 // type bit for plain file
)

// Dir.Mode bits
const (
	DMDIR    = 0x80000000 // directories
	DMAPPEND = 0x40000000 // append-only
	DMEXCL   = 0x20000000 // exclusvie use (only one open handle allowed)
	DMMOUNT  = 0x10000000 // mount points
	DMAUTH   = 0x08000000 // authentication file
	DMTMP    = 0x04000000 // non-backed-up files
	DMREAD   = 0x4        // mode bit for read permission
	DMWRITE  = 0x2        // mode bit for write permission
	DMEXEC   = 0x1        // mode bit for execute permission
)

// 9P2000.u extensions
const (
	DMSYMLINK   = 0x02000000 // symbolic links
	DMDEVICE    = 0x00800000 // device files
	DMNAMEDPIPE = 0x00200000 // named pipe
	DMSOCKET    = 0x00100000 // socket
	DMSETUID    = 0x00080000 // setuid
	DMSETGID    = 0x00040000 // setgid
	// sticky bit
	// DMSETVTX = 0x00000000
)

const (
	NOTAG = 0xffff     // Dummy tag
	NOFID = 0xffffffff // Dummy fid
	NOUID = 0xffffffff // Dummy uid
)
