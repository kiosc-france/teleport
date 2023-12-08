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
	"syscall"
)

// SocketType encodes the type of desired socket (datagram or stream).
type SocketType bool

const (
	SocketTypeDatagram SocketType = true
	SocketTypeStream   SocketType = false
)

// proto converts SocketType into the expected value for use as the
// 'proto' argument in syscall.Socketpair.
func (s SocketType) proto() int {
	var p int
	switch s {
	case SocketTypeDatagram:
		p = syscall.SOCK_DGRAM
	case SocketTypeStream:
		p = syscall.SOCK_STREAM
	}
	return p
}
