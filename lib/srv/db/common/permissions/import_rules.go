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
	"sort"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/utils"
)

// ApplyDatabaseObjectImportRules applies the given set of rules onto a set of objects coming from a same database.
// Returns a fresh copy of a subset of supplied objects, filtered and modified.
// For the object to be returned, it must match at least one rule.
// The modification consists of application of extra labels, per matching mappings.
func ApplyDatabaseObjectImportRules(rules []types.DatabaseObjectImportRule, database types.Database, objs []types.DatabaseObject) []types.DatabaseObject {
	// sort: rules with higher priorities are applied last.
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority() < rules[j].Priority()
	})

	// filter rules: keep those with matching labels
	// we only need mappings from the rules, so extract those.
	var mappings []types.DatabaseObjectImportRuleMapping
	for _, rule := range rules {
		if types.MatchLabels(database, rule.DatabaseLabels()) {
			mappings = append(mappings, rule.Mappings()...)
		}
	}

	// anything to do?
	if len(mappings) == 0 {
		return nil
	}

	var out []types.DatabaseObject

	// find all objects that match any of the rules
	for _, obj := range objs {
		var objClone types.DatabaseObject

		// apply each mapping in order.
		for _, mapping := range mappings {
			// a mapping contains multiple matchers; an object must match at least one.
			// if there are no matchers, the mapping will always match.
			// the matching is applied to the object spec, so disregarding any labels on the object itself.
			if utils.Any(mapping.ObjectMatches, matchDatabaseObjectSpec(obj.GetSpec())) {
				if objClone == nil {
					objClone = obj.Copy()
				}

				// mapping applies additional (static) labels
				labels := objClone.GetAllLabels()
				if labels == nil {
					labels = map[string]string{}
				}
				for k, v := range mapping.AddLabels {
					labels[k] = v
				}
				objClone.SetStaticLabels(labels)
			}
		}

		if objClone != nil {
			out = append(out, objClone)
		}
	}

	return out
}

// matchDatabaseObjectSpec is a helper function to match object's spec with a matcher.
// It returns a function suitable to be passed to utils.Any.
func matchDatabaseObjectSpec(spec types.DatabaseObjectSpec) func(matcher types.DatabaseObjectSpec) bool {
	mapping := databaseObjectSpecToMap(spec)
	return func(matcher types.DatabaseObjectSpec) bool {
		return matchLabels(mapping, databaseObjectSpecToMap(matcher))
	}
}

type allLabels struct {
	types.ResourceWithLabels
	resource map[string]string
}

func (a *allLabels) GetAllLabels() map[string]string { return a.resource }

// matchLabels is a helper to pass a custom set of object labels to the types.MatchLabels function.
func matchLabels(resource map[string]string, labels map[string]string) bool {
	return types.MatchLabels(&allLabels{resource: resource}, labels)
}

// databaseObjectToMap combines object labels and attributes, with the latter having preference.
func databaseObjectToMap(d types.DatabaseObject) map[string]string {
	out := make(map[string]string)
	for k, v := range d.GetAllLabels() {
		if v != "" {
			out[k] = v
		}
	}
	for k, v := range databaseObjectSpecToMap(d.GetSpec()) {
		out[k] = v
	}
	return out
}

// databaseObjectSpecToMap summarizes object attributes. Only non-empty values are returned.
func databaseObjectSpecToMap(d types.DatabaseObjectSpec) map[string]string {
	m := make(map[string]string)
	add := func(key, val string) {
		if val != "" {
			m[key] = val
		}
	}

	for key, val := range d.Attributes {
		add(key, val)
	}

	// standard attributes have higher priority, so are added last.
	add("name", d.Name)
	add("schema", d.Schema)
	add("database", d.Database)
	add("service_name", d.ServiceName)
	add("protocol", d.Protocol)
	add("object_kind", d.ObjectKind)

	return m
}
