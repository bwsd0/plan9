package client

/*
For these routines, addr is a network address of the form:

	network!netaddr!service
	network!netaddr
	netaddr

Network is tcp, udp, unix or the special token, net. Net is a free variable that
stands for any network in common between the source and the host netaddr.
Netaddr can be a host name, a domain name, or a network address.

On Plan 9, the dir argument is a path name to a line directory that has files
for accessing the connection. To keep the same function signatures, the Unix
port of these routines uses strings of the form /dev/fd/n instead of line
directory paths. These strings should be treated as opaque data and ignored.
*/

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

// Dial makes a call to destination addr on a multiplexed network.
//
// If the network in addr is net, dial will try all networks in succession that
// are common between the source and destination until the call succeeds.
func Dial(network, addr string) (*Conn, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(c)
}

// DialService is a convenience function that wraps Dial by calling the named
// service.
func DialService(service string) (*Conn, error) {
	ns := Namespace()
	return Dial("unix", ns+"/"+service)
}

// Mount mounts a 9P server's files  into the file system.
func Mount(network, addr string) (*Fsys, error) {
	c, err := Dial(network, addr)
	if err != nil {
		return nil, err
	}
	fsys, err := c.Attach(nil, getuser(), "")
	if err != nil {
		c.Close()
	}
	return fsys, err
}

func MountService(service string) (*Fsys, error) {
	c, err := DialService(service)
	if err != nil {
		return nil, err
	}
	fsys, err := c.Attach(nil, getuser(), "")
	if err != nil {
		c.Close()
	}
	return fsys, err
}

func MountServiceAname(service, aname string) (*Fsys, error) {
	c, err := DialService(service)
	if err != nil {
		return nil, err
	}
	fsys, err := c.Attach(nil, getuser(), aname)
	if err != nil {
		c.Close()
	}
	return fsys, err
}

// Namespace returns the path to the name space directory.
func Namespace() string {
	disp := os.Getenv("DISPLAY")
	if disp == "" && runtime.GOOS == "darwin" {
		disp = ":0.0"
	}
	if i := strings.LastIndex(disp, ":"); i < 0 {
		log.Fatalf("bad display: %s", disp)
	}
	// canonicalize $host:$display.$screen => $host:$display
	if i := strings.LastIndex(disp, "."); i > 0 {
		disp = disp[0 : i-1]
	}
	if runtime.GOOS == "darwin" {
		// Turn /tmp/launch/:0 into _tmp_launch_:0 (OS X 10.5).
		disp = strings.Replace(disp, "/", "_", -1)
	}
	// NOTE: plan9port creates this directory on demand.
	// Maybe someday we'll need to do that.

	ns := os.Getenv("NAMESPACE")
	if ns == "" {
		ns = fmt.Sprintf("/tmp/ns/%s.%s", getuser(), disp)
	}

	return ns
}
