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

package version

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/gravitational/trace"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/integrations/kube-agent-updater/pkg/basichttp"
	"github.com/gravitational/teleport/integrations/kube-agent-updater/pkg/constants"
)

const basicHTTPTestPath = "/v1/cloud-stable"

func Test_basicHTTPVersionClient_Get(t *testing.T) {
	mock := basichttp.NewServerMock(basicHTTPTestPath + "/" + constants.VersionPath)
	t.Cleanup(mock.Srv.Close)
	serverURL, err := url.Parse(mock.Srv.URL)
	serverURL.Path = basicHTTPTestPath
	require.NoError(t, err)
	ctx := context.Background()

	tests := []struct {
		name       string
		statusCode int
		response   string
		expected   string
		headers    map[string]string
		assertErr  require.ErrorAssertionFunc
	}{
		{
			name:       "all good",
			statusCode: http.StatusOK,
			response:   "12.0.3",
			expected:   "v12.0.3",
			assertErr:  require.NoError,
		},
		{
			name:       "all good with newline",
			statusCode: http.StatusOK,
			response:   "12.0.3\n",
			expected:   "v12.0.3",
			assertErr:  require.NoError,
		},
		{
			name:       "non-semver",
			statusCode: http.StatusOK,
			response:   "hello",
			expected:   "",
			assertErr: func(t2 require.TestingT, err2 error, _ ...interface{}) {
				require.IsType(t2, &trace.BadParameterError{}, trace.Unwrap(err2))
			},
		},
		{
			name:       "empty",
			statusCode: http.StatusOK,
			response:   "",
			expected:   "",
			assertErr: func(t2 require.TestingT, err2 error, _ ...interface{}) {
				require.IsType(t2, &trace.BadParameterError{}, trace.Unwrap(err2))
			},
		},
		{
			name:       "non-200 response",
			statusCode: http.StatusInternalServerError,
			response:   "ERROR - SOMETHING WENT WRONG",
			expected:   "",
			assertErr:  require.Error,
		},
		{
			name:       "agent metadata headers",
			statusCode: http.StatusOK,
			response:   "14.2.0",
			expected:   "v14.2.0",
			headers: map[string]string{
				constants.AgentVersionHeader:   "14.0.0",
				constants.UpdaterVersionHeader: "14.1.0",
			},
			assertErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basicHTTPVersionClient{
				baseURL:      serverURL,
				client:       &basichttp.Client{Client: mock.Srv.Client()},
				extraHeaders: tt.headers,
			}
			mock.SetResponse(t, tt.statusCode, tt.response, tt.headers)
			result, err := b.Get(ctx)
			tt.assertErr(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func Test_basicHTTPVersionClient_SetHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
	}{
		{
			name: "set agent metadata",
			headers: map[string]string{
				constants.AgentVersionHeader:   "14.0.0",
				constants.UpdaterVersionHeader: "14.1.0",
			},
		},
		{
			name:    "no exta headers",
			headers: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basicHTTPVersionClient{
				extraHeaders: tt.headers,
			}
			for header, value := range tt.headers {
				b.SetHeader(header, value)
			}
			require.Equal(t, tt.headers, b.extraHeaders)
		})
	}
}
