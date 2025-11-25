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
	"strings"

	"github.com/altinity/clickhouse-operator/pkg/apis/metrics"
	"github.com/altinity/clickhouse-operator/pkg/util"
)

type metricsFilter struct {
	dropMetrics        []string
	dropMetricPrefixes []string
	keepMetrics        []string
	dropLabels         []string
	keepLabels         []string
}

func newMetricsFilter(filters *metrics.MetricsFilters) *metricsFilter {
	if filters == nil {
		return &metricsFilter{}
	}

	return &metricsFilter{
		dropMetrics:        util.NonEmpty(filters.DropMetrics),
		dropMetricPrefixes: util.NonEmpty(filters.DropMetricPrefixes),
		keepMetrics:        util.NonEmpty(filters.KeepMetrics),
		dropLabels:         util.NonEmpty(filters.DropLabels),
		keepLabels:         util.NonEmpty(filters.KeepLabels),
	}
}

func (f *metricsFilter) shouldDropMetric(name string) bool {
	if f == nil {
		return false
	}

	if len(f.keepMetrics) > 0 {
		if f.match(name, f.keepMetrics) {
			return false
		}
		return true
	}

	if f.match(name, f.dropMetrics) {
		return true
	}

	for _, prefix := range f.dropMetricPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

func (f *metricsFilter) filterLabels(labels map[string]string) map[string]string {
	if f == nil {
		return labels
	}

	// Do not change original map
	filtered := make(map[string]string, len(labels))
	for name, value := range labels {
		if f.shouldDropLabel(name) {
			continue
		}
		filtered[name] = value
	}

	return filtered
}

func (f *metricsFilter) shouldDropLabel(name string) bool {
	if f == nil {
		return false
	}

	if len(f.keepLabels) > 0 {
		if f.match(name, f.keepLabels) {
			return false
		}
		return true
	}

	return f.match(name, f.dropLabels)
}

func (f *metricsFilter) match(needle string, haystack []string) bool {
	if len(haystack) == 0 {
		return false
	}
	return util.MatchArrayOfRegexps(needle, haystack)
}
