// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clickhouse

import (
	"reflect"
	"testing"

	metricsapi "github.com/altinity/clickhouse-operator/pkg/apis/metrics"
)

func TestMetricsFilter_DropByRules(t *testing.T) {
	filter := newMetricsFilter(&metricsapi.MetricsFilters{
		DropMetrics:        []string{"exact_match", "^regex_.*"},
		DropMetricPrefixes: []string{"prefix_"},
	})

	cases := []struct {
		name string
		drop bool
	}{
		{name: "exact_match", drop: true},
		{name: "regex_value", drop: true},
		{name: "prefix_metric", drop: true},
		{name: "other_metric", drop: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := filter.shouldDropMetric(tc.name); got != tc.drop {
				t.Fatalf("expected drop=%t for %s, got %t", tc.drop, tc.name, got)
			}
		})
	}
}

func TestMetricsFilter_KeepOverridesDrop(t *testing.T) {
	filter := newMetricsFilter(&metricsapi.MetricsFilters{
		KeepMetrics:        []string{"keep_me", "regex_keep_.*"},
		DropMetrics:        []string{"keep_me", "regex_keep_drop"},
		DropMetricPrefixes: []string{"regex_"},
	})

	if filter.shouldDropMetric("keep_me") {
		t.Fatal("keep list should override drop rules for exact match")
	}
	if filter.shouldDropMetric("regex_keep_value") {
		t.Fatal("keep list should override drop rules for regex match")
	}
	if !filter.shouldDropMetric("other_metric") {
		t.Fatal("metrics not present in keep list must be dropped when keep list is set")
	}
}

func TestMetricsFilter_LabelFiltering(t *testing.T) {
	filter := newMetricsFilter(&metricsapi.MetricsFilters{
		DropLabels: []string{"drop_me", "^skip_.*"},
		KeepLabels: []string{"keep_me", "important_.*"},
	})

	labels := map[string]string{
		"drop_me":         "1",
		"skip_value":      "2",
		"keep_me":         "3",
		"important_label": "4",
		"other":           "5",
	}

	expected := map[string]string{
		"keep_me":         "3",
		"important_label": "4",
	}

	if got := filter.filterLabels(labels); !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected labels after filtering: %+v, expected %+v", got, expected)
	}
}

func TestMetricsFilter_DefaultsNoop(t *testing.T) {
	filter := newMetricsFilter(nil)

	if filter.shouldDropMetric("anything") {
		t.Fatal("default filter must not drop metrics")
	}

	labels := map[string]string{"a": "1", "b": "2"}
	if got := filter.filterLabels(labels); !reflect.DeepEqual(got, labels) {
		t.Fatalf("default filter must keep labels unchanged, got %+v", got)
	}
}
