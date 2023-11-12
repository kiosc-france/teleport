/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"fmt"
	"os"

	"github.com/gravitational/trace"

	apievents "github.com/gravitational/teleport/api/types/events"
	"github.com/gravitational/teleport/lib/client/terminal"
	"github.com/gravitational/teleport/lib/events"
	"github.com/gravitational/teleport/lib/session"
)

// playFromFileStreamer implements [player.Streamer] for
// streaming from a local file.
type playFromFileStreamer struct {
	filename string
}

func (p *playFromFileStreamer) StreamSessionEvents(
	ctx context.Context,
	sessionID session.ID,
	startIndex int64,
) (chan apievents.AuditEvent, chan error) {
	evts := make(chan apievents.AuditEvent)
	errs := make(chan error, 1)

	go func() {
		f, err := os.Open(p.filename)
		if err != nil {
			errs <- trace.ConvertSystemError(err)
			return
		}
		defer f.Close()

		pr := events.NewProtoReader(f)
		for i := int64(0); ; i++ {
			evt, err := pr.Read(ctx)
			if err != nil {
				errs <- trace.Wrap(err)
				return
			}

			if i >= startIndex {
				evts <- evt
			}
		}
	}()

	return evts, errs
}

// timestampFrame prints 'event timestamp' in the top right corner of the
// terminal after playing every 'print' event
func timestampFrame(term *terminal.Terminal, message string) {
	const (
		saveCursor    = "7"
		restoreCursor = "8"
	)
	width, _, err := term.Size()
	if err != nil {
		return
	}
	esc := func(s string) {
		os.Stdout.Write([]byte("\x1b" + s))
	}
	esc(saveCursor)
	defer esc(restoreCursor)

	// move cursor to -10:0
	// TODO(timothyb89): message length does not account for unicode characters
	// or ANSI sequences.
	esc(fmt.Sprintf("[%d;%df", 0, int(width)-len(message)))
	os.Stdout.WriteString(message)
}
