// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package permissions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api/types"
)

type mockGetter struct {
	allow []types.DatabasePermission
	deny  []types.DatabasePermission
}

func (m mockGetter) GetDatabasePermissions() (allow types.DatabasePermissions, deny types.DatabasePermissions) {
	return m.allow, m.deny
}

func TestCalculatePermissions(t *testing.T) {
	mkDatabaseObject := func(name string, spec types.DatabaseObjectSpec) types.DatabaseObject {
		out, err := types.NewDatabaseObject(types.Metadata{Name: name}, spec)
		require.NoError(t, err)
		return out
	}

	tests := []struct {
		name   string
		getter GetDatabasePermissions
		objs   []types.DatabaseObject
		want   PermissionSet
	}{
		{
			name: "Allow all permissions for all objects",
			getter: &mockGetter{
				allow: []types.DatabasePermission{
					{
						ObjectKind: "table",
						Permission: "SELECT",
						Match:      nil,
					},
					{
						ObjectKind: "schema",
						Permission: "DELETE",
						Match:      nil,
					},
				},
				deny: nil,
			},
			objs: []types.DatabaseObject{
				mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}),
				mkDatabaseObject("bar", types.DatabaseObjectSpec{ObjectKind: "schema"}),
			},
			want: PermissionSet{
				"SELECT": {
					mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}),
				},
				"DELETE": {
					mkDatabaseObject("bar", types.DatabaseObjectSpec{ObjectKind: "schema"}),
				},
			},
		},
		{
			name: "Deny removes permissions",
			getter: &mockGetter{
				allow: []types.DatabasePermission{
					{
						ObjectKind: "table",
						Permission: "SELECT",
					},
					{
						ObjectKind: "table",
						Permission: "INSERT",
					},
					{
						ObjectKind: "schema",
						Permission: "SELECT",
					},
					{
						ObjectKind: "schema",
						Permission: "DELETE",
					},
				},
				deny: []types.DatabasePermission{
					{
						ObjectKind: "table",
						Permission: "*",
					},
					{
						ObjectKind: "schema",
						Permission: "DELETE",
					},
				},
			},
			objs: []types.DatabaseObject{
				mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}),
				mkDatabaseObject("bar", types.DatabaseObjectSpec{ObjectKind: "schema"}),
			},
			want: PermissionSet{
				"SELECT": {
					mkDatabaseObject("bar", types.DatabaseObjectSpec{ObjectKind: "schema"}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePermissions(tt.getter, tt.objs)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
