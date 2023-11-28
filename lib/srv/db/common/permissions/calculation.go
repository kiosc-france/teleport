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
	"strings"

	"github.com/gravitational/teleport/api/types"
)

type PermissionSet map[string][]types.DatabaseObject

type GetDatabasePermissions interface {
	GetDatabasePermissions() (allow types.DatabasePermissions, deny types.DatabasePermissions)
}

// CalculatePermissions calculates the effective permissions for a set of database objects based on the provided allow and deny permissions.
func CalculatePermissions(getter GetDatabasePermissions, objs []types.DatabaseObject) (PermissionSet, error) {
	allow, deny := getter.GetDatabasePermissions()

	out := map[string][]types.DatabaseObject{}

	for _, obj := range objs {
		// we use both object attributes (from spec) and labels (from metadata) for matching.
		objMap := databaseObjectToMap(obj)

		permsToAdd := map[string]string{}

		// add allowed permissions
		for _, perm := range allow {
			if obj.GetSpec().ObjectKind == perm.ObjectKind {
				if matchLabels(objMap, perm.Match) {
					permsToAdd[strings.TrimSpace(strings.ToUpper(perm.Permission))] = perm.Permission
				}
			}
		}

		// remove denied permissions
		for _, perm := range deny {
			if obj.GetSpec().ObjectKind == perm.ObjectKind {
				if matchLabels(objMap, perm.Match) {
					// wildcard clears the permissions
					if perm.Permission == types.Wildcard {
						clear(permsToAdd)
						break
					}

					delete(permsToAdd, strings.TrimSpace(strings.ToUpper(perm.Permission)))
				}
			}
		}

		for _, perm := range permsToAdd {
			out[perm] = append(out[perm], obj)
		}
	}

	return out, nil
}
