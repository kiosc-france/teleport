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

package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gravitational/trace"

	"github.com/gravitational/teleport/api/utils"
)

// DatabaseObject represents a single object in a database, e.g. a table.
type DatabaseObject interface {
	ResourceWithLabels

	// GetSpec returns the database object spec
	GetSpec() DatabaseObjectSpec

	// Copy returns a fresh copy of the database object
	Copy() DatabaseObject
}

var _ DatabaseObject = (*DatabaseObjectV1)(nil)

func NewDatabaseObject(metadata Metadata, spec DatabaseObjectSpec) (DatabaseObject, error) {
	o := &DatabaseObjectV1{
		ResourceHeader: ResourceHeader{
			Metadata: metadata,
		},
		Spec: spec,
	}
	err := o.CheckAndSetDefaults()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return o, nil
}

func (d *DatabaseObjectV1) CheckAndSetDefaults() error {
	d.setStaticFields()
	if d.Spec.ObjectKind == "" {
		return trace.BadParameter("empty ObjectKind in DatabaseObject")
	}
	return trace.Wrap(d.Metadata.CheckAndSetDefaults())
}

func (d *DatabaseObjectV1) setStaticFields() {
	d.Kind = KindDatabaseObject
	d.Version = V1
}

func (d *DatabaseObjectV1) Copy() DatabaseObject {
	return proto.Clone(d).(*DatabaseObjectV1)
}

func (d *DatabaseObjectV1) MatchSearch(searchValues []string) bool {
	values := append(utils.MapToStrings(d.GetAllLabels()), d.GetName())
	return MatchSearch(values, searchValues, nil)
}

func (d *DatabaseObjectV1) GetSpec() DatabaseObjectSpec {
	return d.Spec
}

// DatabaseObjectImportRule is a global database object import rule.
type DatabaseObjectImportRule interface {
	ResourceWithLabels

	// Priority represents the priority of the rule application. Lower numbered rules will be applied first.
	Priority() int

	// DatabaseLabels is a set of labels which must match the database for the rule to be applied.
	DatabaseLabels() map[string]string

	// Mappings is a list of matches that will map match conditions to labels.
	Mappings() []DatabaseObjectImportRuleMapping
}

var _ DatabaseObjectImportRule = (*DatabaseObjectImportRuleV1)(nil)

func NewDatabaseObjectImportRule(metadata Metadata, spec DatabaseObjectImportRuleSpec) (DatabaseObjectImportRule, error) {
	o := &DatabaseObjectImportRuleV1{
		ResourceHeader: ResourceHeader{
			Metadata: metadata,
		},
		Spec: spec,
	}
	err := o.CheckAndSetDefaults()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return o, nil
}

func (d *DatabaseObjectImportRuleV1) setStaticFields() {
	d.Kind = KindDatabaseObjectImportRule
	d.Version = V1
}

func (d *DatabaseObjectImportRuleV1) CheckAndSetDefaults() error {
	d.setStaticFields()
	if err := d.Metadata.CheckAndSetDefaults(); err != nil {
		return trace.Wrap(err)
	}

	if d.Spec.Priority < 0 {
		return trace.BadParameter("priority must be a positive number")
	}

	return nil
}

func (d *DatabaseObjectImportRuleV1) MatchSearch(searchValues []string) bool {
	values := append(utils.MapToStrings(d.GetAllLabels()), d.GetName())
	return MatchSearch(values, searchValues, nil)
}

func (d *DatabaseObjectImportRuleV1) Priority() int {
	return int(d.Spec.Priority)
}

func (d *DatabaseObjectImportRuleV1) Mappings() []DatabaseObjectImportRuleMapping {
	return d.Spec.Mappings
}

func (d *DatabaseObjectImportRuleV1) DatabaseLabels() map[string]string {
	return d.Spec.DatabaseLabels
}

type DatabasePermissions []DatabasePermission
