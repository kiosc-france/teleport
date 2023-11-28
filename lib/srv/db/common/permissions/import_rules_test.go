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

func TestDatabaseObjectSpecToMap(t *testing.T) {
	tests := []struct {
		name string
		spec types.DatabaseObjectSpec
		want map[string]string
	}{
		{
			name: "full object with all attributes",
			spec: types.DatabaseObjectSpec{
				Name:        "my-table",
				Schema:      "public",
				Database:    "sales",
				ServiceName: "postgres",
				Protocol:    "postgres",
				ObjectKind:  "table",
				Attributes: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: map[string]string{
				"name":         "my-table",
				"schema":       "public",
				"database":     "sales",
				"service_name": "postgres",
				"protocol":     "postgres",
				"object_kind":  "table",
				"key1":         "value1",
				"key2":         "value2",
			},
		},
		{
			name: "empty attributes are skipped",
			spec: types.DatabaseObjectSpec{
				Name:        "my-table",
				Schema:      "public",
				Database:    "",
				ServiceName: "postgres",
				Protocol:    "",
				ObjectKind:  "table",
			},
			want: map[string]string{
				"name":         "my-table",
				"schema":       "public",
				"service_name": "postgres",
				"object_kind":  "table",
			},
		},
		{
			name: "standard attributes take precedence",
			spec: types.DatabaseObjectSpec{
				Name:        "my-table",
				Schema:      "public",
				Database:    "sales",
				ServiceName: "postgres",
				ObjectKind:  "table",
				Attributes: map[string]string{
					"key1":   "value1",
					"key2":   "value2",
					"name":   "override", // Overriding name
					"schema": "override", // Overriding schema
				},
			},
			want: map[string]string{
				"name":         "my-table", // Standard name takes precedence
				"schema":       "public",   // Standard schema takes precedence
				"database":     "sales",
				"service_name": "postgres",
				"object_kind":  "table",
				"key1":         "value1",
				"key2":         "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := databaseObjectSpecToMap(tt.spec)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDatabaseObjectToMap(t *testing.T) {
	type option func(db types.DatabaseObject) error

	mkDatabaseObject := func(name string, spec types.DatabaseObjectSpec, options ...option) types.DatabaseObject {
		out, err := types.NewDatabaseObject(types.Metadata{Name: name}, spec)
		require.NoError(t, err)
		for _, opt := range options {
			require.NoError(t, opt(out))
		}

		return out
	}

	tests := []struct {
		name string
		obj  types.DatabaseObject
		want map[string]string
	}{
		{
			name: "object without labels",
			obj:  mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table", Protocol: "postgres"}),
			want: map[string]string{
				"protocol":    "postgres",
				"object_kind": "table",
			},
		},
		{
			name: "object with some labels",
			obj: mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table", Protocol: "postgres"}, func(db types.DatabaseObject) error {
				db.SetStaticLabels(map[string]string{
					"dev_access":     "ro",
					"flag_from_prod": "dummy",
					"database":       "static-flag-value",
				})
				return nil
			}),
			want: map[string]string{
				"protocol":       "postgres",          // standard attribute
				"object_kind":    "table",             // standard attribute
				"database":       "static-flag-value", // static flag: standard attribute was empty, so static flag took precedence.
				"dev_access":     "ro",                // static flag
				"flag_from_prod": "dummy",             // static flag
			},
		},
		{
			name: "object with overlapping labels",
			obj: mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table", Protocol: "postgres", Attributes: map[string]string{"domain": "finance"}}, func(db types.DatabaseObject) error {
				db.SetStaticLabels(map[string]string{
					"protocol":       "mysql",       // ignored: overridden by a standard attribute
					"dev_access":     "ro",          // static flag, used
					"flag_from_prod": "dummy",       // static flag, used
					"domain":         "engineering", // ignored: overridden by custom attribute
				})
				return nil
			}),
			want: map[string]string{
				"protocol":       "postgres", // standard object attribute
				"object_kind":    "table",    // standard object attribute
				"domain":         "finance",  // custom object attribute, taking precedence over static label
				"dev_access":     "ro",       // static label
				"flag_from_prod": "dummy",    // static label
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := databaseObjectToMap(tt.obj)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMatchLabels(t *testing.T) {
	tests := []struct {
		name     string
		object   map[string]string
		matcher  map[string]string
		expected bool
	}{
		{
			name:     "empty mapping and matcher",
			object:   map[string]string{},
			matcher:  map[string]string{},
			expected: true,
		},
		{
			name:     "empty matcher",
			object:   map[string]string{"key1": "value1", "key2": "value2"},
			matcher:  map[string]string{},
			expected: true,
		},
		{
			name:     "matcher with all matching values",
			object:   map[string]string{"key1": "value1", "key2": "value2"},
			matcher:  map[string]string{"key1": "value1", "key2": "value2"},
			expected: true,
		},
		{
			name:     "matcher with missing key in mapping",
			object:   map[string]string{"key1": "value1"},
			matcher:  map[string]string{"key1": "value1", "key2": "value2"},
			expected: false,
		},
		{
			name:     "matcher with non-matching value",
			object:   map[string]string{"key1": "value1", "key2": "value3"},
			matcher:  map[string]string{"key1": "value1", "key2": "value2"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchLabels(tt.object, tt.matcher)
			require.Equal(t, tt.expected, result)

			result2 := types.MatchLabels(&allLabels{resource: tt.object}, tt.matcher)
			require.Equal(t, tt.expected, result2)
		})
	}
}

func TestApplyDatabaseObjectImportRules(t *testing.T) {
	mkDatabase := func(name string, labels map[string]string) *types.DatabaseV3 {
		db, err := types.NewDatabaseV3(types.Metadata{
			Name:   name,
			Labels: labels,
		}, types.DatabaseSpecV3{
			Protocol: "postgres",
			URI:      "localhost:5252",
		})
		require.NoError(t, err)
		return db
	}

	type option func(db types.DatabaseObject) error

	mkDatabaseObject := func(name string, spec types.DatabaseObjectSpec, options ...option) types.DatabaseObject {
		out, err := types.NewDatabaseObject(types.Metadata{Name: name}, spec)
		require.NoError(t, err)
		for _, opt := range options {
			require.NoError(t, opt(out))
		}

		return out
	}

	tests := []struct {
		name     string
		rules    []types.DatabaseObjectImportRule
		database types.Database
		objs     []types.DatabaseObject
		want     []types.DatabaseObject
	}{
		{
			name:     "empty inputs",
			rules:    []types.DatabaseObjectImportRule{},
			database: mkDatabase("dummy", map[string]string{"env": "prod"}),
			objs:     nil,
			want:     nil,
		},
		{
			name: "database labels are matched by the rules",
			rules: []types.DatabaseObjectImportRule{
				&types.DatabaseObjectImportRuleV1{
					Spec: types.DatabaseObjectImportRuleSpec{
						Priority:       10,
						DatabaseLabels: map[string]string{"env": "dev"},
						Mappings: []types.DatabaseObjectImportRuleMapping{
							{
								AddLabels: map[string]string{
									"dev_access":    "rw",
									"flag_from_dev": "dummy",
								},
							},
						},
					},
				},
				&types.DatabaseObjectImportRuleV1{
					Spec: types.DatabaseObjectImportRuleSpec{
						Priority:       20,
						DatabaseLabels: map[string]string{"env": "prod"},
						Mappings: []types.DatabaseObjectImportRuleMapping{
							{
								AddLabels: map[string]string{
									"dev_access":     "ro",
									"flag_from_prod": "dummy",
								},
							},
						},
					},
				},
			},
			database: mkDatabase("dummy", map[string]string{"env": "prod"}),
			objs: []types.DatabaseObject{
				mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}),
			},
			want: []types.DatabaseObject{
				mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}, func(db types.DatabaseObject) error {
					db.SetStaticLabels(map[string]string{
						"dev_access":     "ro",
						"flag_from_prod": "dummy",
					})
					return nil
				}),
			},
		},
		{
			name: "rule priorities are applied",
			rules: []types.DatabaseObjectImportRule{
				&types.DatabaseObjectImportRuleV1{
					Spec: types.DatabaseObjectImportRuleSpec{
						Priority: 10,
						Mappings: []types.DatabaseObjectImportRuleMapping{
							{
								AddLabels: map[string]string{
									"dev_access":    "rw",
									"flag_from_dev": "dummy",
								},
							},
						},
					},
				},
				&types.DatabaseObjectImportRuleV1{
					Spec: types.DatabaseObjectImportRuleSpec{
						Priority: 20,
						Mappings: []types.DatabaseObjectImportRuleMapping{
							{
								AddLabels: map[string]string{
									"dev_access":     "ro",
									"flag_from_prod": "dummy",
								},
							},
						},
					},
				},
			},
			database: mkDatabase("dummy", map[string]string{}),
			objs: []types.DatabaseObject{
				mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}),
			},
			want: []types.DatabaseObject{
				mkDatabaseObject("foo", types.DatabaseObjectSpec{ObjectKind: "table"}, func(db types.DatabaseObject) error {
					db.SetStaticLabels(map[string]string{
						"dev_access":     "ro",
						"flag_from_dev":  "dummy",
						"flag_from_prod": "dummy",
					})
					return nil
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := ApplyDatabaseObjectImportRules(tt.rules, tt.database, tt.objs)
			require.Equal(t, tt.want, out)
		})
	}
}
