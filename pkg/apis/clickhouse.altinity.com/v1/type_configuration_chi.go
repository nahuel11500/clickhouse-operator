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

package v1

import "github.com/altinity/clickhouse-operator/pkg/apis/common/types"

const (
	// CommonConfigDirClickHouse specifies folder's name, where generated common XML files for ClickHouse would be placed
	CommonConfigDirClickHouse = "config.d"

	// UsersConfigDirClickHouse specifies folder's name, where generated users XML files for ClickHouse would be placed
	UsersConfigDirClickHouse = "users.d"

	// HostConfigDirClickHouse specifies folder's name, where generated host XML files for ClickHouse would be placed
	HostConfigDirClickHouse = "conf.d"

	// TemplatesDirClickHouse specifies folder's name where ClickHouseInstallationTemplates are located
	TemplatesDirClickHouse = "templates.d"
)

const (
	// CommonConfigDirKeeper specifies folder's name, where generated common XML files for ClickHouse would be placed
	CommonConfigDirKeeper = "keeper_config.d"

	// UsersConfigDirKeeper specifies folder's name, where generated users XML files for ClickHouse would be placed
	UsersConfigDirKeeper = "users.d"

	// HostConfigDirKeeper specifies folder's name, where generated host XML files for ClickHouse would be placed
	HostConfigDirKeeper = "conf.d"

	// TemplatesDirKeeper specifies folder's name where ClickHouseInstallationTemplates are located
	TemplatesDirKeeper = "templates.d"
)

// Configuration defines configuration section of .spec
type Configuration struct {
	Zookeeper *ZookeeperConfig      `json:"zookeeper,omitempty" yaml:"zookeeper,omitempty"`
	Users     *Settings             `json:"users,omitempty"     yaml:"users,omitempty"`
	Profiles  *Settings             `json:"profiles,omitempty"  yaml:"profiles,omitempty"`
	Quotas    *Settings             `json:"quotas,omitempty"    yaml:"quotas,omitempty"`
	Settings  *Settings             `json:"settings,omitempty"  yaml:"settings,omitempty"`
	Files     *Settings             `json:"files,omitempty"     yaml:"files,omitempty"`
	Metrics   *MetricsConfiguration `json:"metrics,omitempty"   yaml:"metrics,omitempty"`
	Clusters  []*Cluster            `json:"clusters,omitempty"  yaml:"clusters,omitempty"`
}

// NewConfiguration creates new Configuration objects
func NewConfiguration() *Configuration {
	return new(Configuration)
}

func (c *Configuration) Ensure() *Configuration {
	if c == nil {
		c = NewConfiguration()
	}
	return c
}

func (c *Configuration) GetUsers() *Settings {
	if c == nil {
		return nil
	}
	return c.Users
}

func (c *Configuration) GetProfiles() *Settings {
	if c == nil {
		return nil
	}
	return c.Profiles
}

func (c *Configuration) GetQuotas() *Settings {
	if c == nil {
		return nil
	}
	return c.Quotas
}

func (c *Configuration) GetSettings() *Settings {
	if c == nil {
		return nil
	}
	return c.Settings
}

func (c *Configuration) GetFiles() *Settings {
	if c == nil {
		return nil
	}
	return c.Files
}

func (c *Configuration) GetMetrics() *MetricsConfiguration {
	if c == nil {
		return nil
	}
	return c.Metrics
}

// MergeFrom merges from specified source
func (c *Configuration) MergeFrom(from *Configuration, _type MergeType) *Configuration {
	if from == nil {
		return c
	}

	c = c.Ensure()

	c.Zookeeper = c.Zookeeper.MergeFrom(from.Zookeeper, _type)
	c.Users = c.Users.MergeFrom(from.Users)
	c.Profiles = c.Profiles.MergeFrom(from.Profiles)
	c.Quotas = c.Quotas.MergeFrom(from.Quotas)
	c.Settings = c.Settings.MergeFrom(from.Settings)
	c.Files = c.Files.MergeFrom(from.Files)
	c.Metrics = c.Metrics.MergeFrom(from.Metrics, _type)

	// TODO merge clusters
	// Copy Clusters for now
	c.Clusters = from.Clusters

	return c
}

type MetricsConfiguration struct {
	Filters *MetricsFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
}

func (m *MetricsConfiguration) MergeFrom(from *MetricsConfiguration, _type MergeType) *MetricsConfiguration {
	if from == nil {
		return m
	}

	if m == nil {
		return from
	}

	m.Filters = m.Filters.MergeFrom(from.Filters, _type)

	return m
}

type MetricsFilters struct {
	DropMetrics        *types.Strings `json:"dropMetrics,omitempty"        yaml:"dropMetrics,omitempty"`
	DropMetricPrefixes *types.Strings `json:"dropMetricPrefixes,omitempty" yaml:"dropMetricPrefixes,omitempty"`
	KeepMetrics        *types.Strings `json:"keepMetrics,omitempty"        yaml:"keepMetrics,omitempty"`
	DropLabels         *types.Strings `json:"dropLabels,omitempty"         yaml:"dropLabels,omitempty"`
	KeepLabels         *types.Strings `json:"keepLabels,omitempty"         yaml:"keepLabels,omitempty"`
}

func (m *MetricsFilters) MergeFrom(from *MetricsFilters, _type MergeType) *MetricsFilters {
	if from == nil {
		return m
	}

	if m == nil {
		return from
	}

	switch _type {
	case MergeTypeFillEmptyValues:
		if !m.DropMetrics.HasValue() {
			m.DropMetrics = m.DropMetrics.MergeFrom(from.DropMetrics)
		}
		if !m.DropMetricPrefixes.HasValue() {
			m.DropMetricPrefixes = m.DropMetricPrefixes.MergeFrom(from.DropMetricPrefixes)
		}
		if !m.KeepMetrics.HasValue() {
			m.KeepMetrics = m.KeepMetrics.MergeFrom(from.KeepMetrics)
		}
		if !m.DropLabels.HasValue() {
			m.DropLabels = m.DropLabels.MergeFrom(from.DropLabels)
		}
		if !m.KeepLabels.HasValue() {
			m.KeepLabels = m.KeepLabels.MergeFrom(from.KeepLabels)
		}
	case MergeTypeOverrideByNonEmptyValues:
		if from.DropMetrics.HasValue() {
			m.DropMetrics = m.DropMetrics.MergeFrom(from.DropMetrics)
		}
		if from.DropMetricPrefixes.HasValue() {
			m.DropMetricPrefixes = m.DropMetricPrefixes.MergeFrom(from.DropMetricPrefixes)
		}
		if from.KeepMetrics.HasValue() {
			m.KeepMetrics = m.KeepMetrics.MergeFrom(from.KeepMetrics)
		}
		if from.DropLabels.HasValue() {
			m.DropLabels = m.DropLabels.MergeFrom(from.DropLabels)
		}
		if from.KeepLabels.HasValue() {
			m.KeepLabels = m.KeepLabels.MergeFrom(from.KeepLabels)
		}
	}

	return m
}
