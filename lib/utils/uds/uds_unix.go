//go:build unix

/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package uds

import (
	"net"
	"os"
	"syscall"

	"github.com/gravitational/trace"
)

// assert that *Conn implements net.Conn.
var _ net.Conn = (*Conn)(nil)

// Conn extends net.UnixConn with additional useful methods.
type Conn struct {
	*net.UnixConn
}

// FromFile attempts to create a [Conn] from a file. The returned conn
// is independent from the file and closing one does not close the other.
func FromFile(fd *os.File) (*Conn, error) {
	fconn, err := net.FileConn(fd)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	uconn, ok := fconn.(*net.UnixConn)
	if !ok {
		return nil, trace.Errorf("unexpected conn type %T (expected %T)", fconn, (*net.UnixConn)(nil))
	}

	return &Conn{uconn}, nil
}

// WriteWithFDs performs a write that may also send associated FDs. Note that unless the
// underlying socket is a datagram socket it is not guaranteed that the exact bytes written
// will match the bytes received with the fds on the reader side.
func (c *Conn) WriteWithFDs(b []byte, fds []*os.File) (n, fdn int, err error) {
	// prefer stack-allocated arrays for small numbers of fds to reduce
	// per-write allocations.
	var sfbuf [64]int
	var fbuf []int
	if len(fds) > 64 {
		fbuf = make([]int, len(fds))
	} else {
		fbuf = sfbuf[:len(fds)]
	}

	for i, fd := range fds {
		fbuf[i] = int(fd.Fd())
	}

	n, _, err = c.WriteMsgUnix(b, syscall.UnixRights(fbuf...), nil)
	if err != nil {
		return n, 0, trace.Wrap(err)
	}

	return n, len(fds), nil
}

const (
	// sizeOfInt is the size of an fd/c_int in bytes. Ints are 32-bit on all
	// platforms supported by teleport.
	sizeOfInt = 4
	// receivedFileName is the name assigned to files received via
	// a unix socket.
	receivedFileName = "uds-received-file"
)

// ReadWithFDs reads data and associated fds. Note that the underlying socket must be a datagram socket
// if you need the bytes read to exactly match the bytes sent with the FDs.
func (c *Conn) ReadWithFDs(b []byte, fds []*os.File) (n, fdn int, err error) {
	// prefer reading oob via stack buffer for smaller expected fd sizes to
	// reduce per-read allocations.
	var sobuf [256]byte
	var obuf []byte
	oobSize := syscall.CmsgSpace(sizeOfInt * len(fds))
	if oobSize > 256 {
		obuf = make([]byte, oobSize)
	} else {
		obuf = sobuf[:oobSize]
	}

	n, oobn, _, _, err := c.ReadMsgUnix(b, obuf)
	if err != nil {
		return n, 0, trace.Wrap(err)
	}

	// of out of band data was sent we need to parse it and extract
	// any fds that were sent across.
	if oobn != 0 {
		scms, err := syscall.ParseSocketControlMessage(obuf[:oobn])
		if err != nil {
			return n, 0, trace.Wrap(err)
		}

		for _, scm := range scms {
			if scm.Header.Level != syscall.SOL_SOCKET || scm.Header.Type != syscall.SCM_RIGHTS {
				// unsupported control message
				continue
			}

			rawFDs, err := syscall.ParseUnixRights(&scm)
			if err != nil {
				continue
			}

			for _, rawFD := range rawFDs {
				if fdn < len(fds) {
					fds[fdn] = os.NewFile(uintptr(rawFD), receivedFileName)
					fdn++
				} else {
					// we got more files than expected, close the excess
					syscall.Close(rawFD)
				}
			}
		}
	}

	return n, fdn, nil
}
