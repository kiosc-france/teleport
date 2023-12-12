// Teleport
// Copyright (C) 2023  Gravitational, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package multiplexer

import (
	"bufio"
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"

	"github.com/gravitational/teleport/lib/limiter"
	"github.com/gravitational/teleport/lib/utils"
)

type maybeSignedPROXYV2Conn struct {
	net.Conn

	remoteAddrOverride atomic.Value

	// verifyContext is used by Close to cancel signature verification
	// operations.
	verifyContext context.Context
	// verifyCancel cancels verifyContext.
	verifyCancel context.CancelCauseFunc

	caGetter    CertAuthorityGetter
	clusterName string

	limiter *limiter.Limiter
	// limiterToken
	limiterToken atomic.Pointer[string]

	mu     sync.Mutex
	reader *bufio.Reader
	parsed bool
}

var _ net.Conn = (*maybeSignedPROXYV2Conn)(nil)

// Read implements [io.Reader] and [net.Conn].
func (c *maybeSignedPROXYV2Conn) Read(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(b) < 1 {
		// don't try to parse the header if we never intended to read in the
		// first place
		return c.reader.Read(b)
	}

	if err := c.ensureParsedLocked(); err != nil {
		return 0, trace.Wrap(err)
	}

	return c.reader.Read(b)
}

func (c *maybeSignedPROXYV2Conn) ensureParsedLocked() (err error) {
	if c.parsed {
		return nil
	}
	c.parsed = true

	defer func() {
		if err != nil {
			c.Close()
		}
	}()

	hasProxyV2, err := readerHasPrefix(c.reader, ProxyV2Prefix)
	if err != nil {
		return trace.Wrap(err)
	}

	if !hasProxyV2 {
		return nil
	}

	proxyline, err := ReadProxyLineV2(c.reader)
	if err != nil {
		return trace.Wrap(err)
	}

	if proxyline == nil {
		return trace.Wrap(ErrNoSignature)
	}

	if err := proxyline.VerifySignature(
		c.verifyContext,
		c.caGetter,
		c.clusterName,
		clockwork.NewRealClock(),
	); err != nil {
		return trace.Wrap(err)
	}

	limiterToken := proxyline.Source.IP.String()
	if err := c.limiter.AcquireConnection(limiterToken); err != nil {
		return trace.Wrap(err)
	}

	if token := c.limiterToken.Swap(&limiterToken); token != nil {
		c.limiter.ReleaseConnection(*token)
	} else {
		// the connection was closed, we have to release the token we just
		// acquired (unless it has already been released by someone else)
		if token := c.limiterToken.Swap(nil); token != nil {
			c.limiter.ReleaseConnection(*token)
		}
	}

	c.remoteAddrOverride.Store(&proxyline.Source)
	c.verifyCancel(nil)

	return nil
}

func readerHasPrefix[B ~[]byte | ~string](r *bufio.Reader, prefix B) (bool, error) {
	for i, b := range []byte(prefix) {
		buf, err := r.Peek(i + 1)
		if err != nil {
			return false, trace.Wrap(err)
		}
		if buf[i] != b {
			return false, nil
		}
	}
	return true, nil
}

func (c *maybeSignedPROXYV2Conn) RemoteAddr() net.Addr {
	if a, ok := c.remoteAddrOverride.Load().(net.Addr); ok {
		return a
	}
	return c.Conn.RemoteAddr()
}

func (c *maybeSignedPROXYV2Conn) Close() error {
	c.verifyCancel(net.ErrClosed)
	if token := c.limiterToken.Swap(nil); token != nil {
		c.limiter.ReleaseConnection(*token)
	}
	return trace.Wrap(c.Conn.Close())
}

func MaybeSignedPROXYV2Listener(listener net.Listener, caGetter CertAuthorityGetter, clusterName string, limiter *limiter.Limiter) net.Listener {
	return &maybeSignedPROXYV2Listener{
		Listener:    listener,
		caGetter:    caGetter,
		clusterName: clusterName,
		limiter:     limiter,
	}
}

type maybeSignedPROXYV2Listener struct {
	net.Listener
	caGetter    CertAuthorityGetter
	clusterName string
	limiter     *limiter.Limiter
}

func (l *maybeSignedPROXYV2Listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	addr, _, err := utils.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
		return nil, trace.NewAggregate(err, c.Close())
	}

	if err := l.limiter.AcquireConnection(addr); err != nil {
		// preserve the LimitExceeded error, just in case
		_ = c.Close()
		return nil, trace.Wrap(err)
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	m := &maybeSignedPROXYV2Conn{
		Conn:          c,
		verifyContext: ctx,
		verifyCancel:  cancel,
		caGetter:      l.caGetter,
		clusterName:   l.clusterName,
		reader:        bufio.NewReader(c),

		limiter: l.limiter,
	}
	m.limiterToken.Store(&addr)
	return m, nil
}
