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

package provider

import (
	"context"
	"strings"

	"github.com/gravitational/trace"
	"github.com/jackc/pgx/v4"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/srv/db/common"
)

// SchemaInfo represents information about all schemas in a database.
type schemaInfo map[string]schema

// schema represents single schema.
type schema struct {
	tables []string
}

const ObjectKindTable = "table"

func FetchDatabaseObjects(ctx context.Context, session *common.Session, conn *pgx.Conn) ([]types.DatabaseObject, error) {
	s, err := getSchemaInfo(ctx, conn)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var out []types.DatabaseObject

	for schemaName, schemaVal := range s {
		for _, table := range schemaVal.tables {
			name := strings.Join([]string{
				session.Database.GetProtocol(),
				session.Database.GetType(),
				session.Database.GetName(),
				ObjectKindTable,
				session.DatabaseName,
				schemaName,
				table,
			}, "/")

			obj := &types.DatabaseObjectV1{
				ResourceHeader: types.ResourceHeader{
					Kind:     types.KindDatabaseObject,
					Version:  types.V1,
					Metadata: types.Metadata{Name: name},
				},
				Spec: types.DatabaseObjectSpec{
					ObjectKind:  ObjectKindTable,
					ServiceName: session.Database.GetName(),
					Protocol:    session.Database.GetProtocol(),
					Database:    session.DatabaseName,
					Schema:      schemaName,
					Name:        table,
				},
			}

			out = append(out, obj)
		}
	}

	return out, nil
}

// getSchemaInfo returns SchemaInfo for given database.
func getSchemaInfo(ctx context.Context, conn *pgx.Conn) (schemaInfo, error) {
	type row struct {
		SchemaName string
		TableName  string
	}

	schemaRows, err := conn.Query(ctx, "SELECT schemaname, tablename FROM pg_tables")
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var tables []row

	for schemaRows.Next() {
		var r row
		if err := schemaRows.Scan(&r.SchemaName, &r.TableName); err != nil {
			return nil, err
		}
		tables = append(tables, r)
	}

	if err := schemaRows.Err(); err != nil {
		return nil, err
	}

	schemas := map[string]schema{}
	for _, table := range tables {
		sch := schemas[table.SchemaName]
		sch.tables = append(sch.tables, table.TableName)

		schemas[table.SchemaName] = sch
	}

	return schemas, nil
}
