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

package automaticupgrades

import (
	"context"
	"net/url"

	"github.com/gravitational/trace"

	"github.com/gravitational/teleport/api/client/proto"
	"github.com/gravitational/teleport/integrations/kube-agent-updater/pkg/maintenance"
	"github.com/gravitational/teleport/integrations/kube-agent-updater/pkg/version"
)

const (
	defaultChannelName        = "default"
	defaultCloudChannelName   = "stable/cloud"
	stableCloudVersionBaseURL = "https://updates.releases.teleport.dev/v1/stable/cloud"
)

// Channels is a map of Channel objects.
type Channels map[string]*Channel

// CheckAndSetDefaults checks that every Channel is valid and initializes them.
// It also creates default channels if they are not already present.
// Cloud must have the `default` and `stable/cloud` channels.
// Self-hosted with automatic upgrades must have the `default` channel.
func (c Channels) CheckAndSetDefaults(features proto.Features) error {
	defaultChannel, err := NewDefaultChannel()
	if err != nil {
		return trace.Wrap(err)
	}

	// If we're on cloud, we need at least "cloud/stable" and "default"
	if features.GetCloud() {
		if _, ok := c[defaultCloudChannelName]; !ok {
			c[defaultCloudChannelName] = defaultChannel
		}
		if _, ok := c[defaultChannelName]; !ok {
			c[defaultChannelName] = c[defaultCloudChannelName]
		}
	}

	// If we're on self-hosted with automatic upgrades, we need a "default" channel
	// We don't want to break existing setups so we'll automatically point to the
	// `cloud/stable` channel.
	// TODO: in v15 make this a hard requirement and error if `default` is not set
	// and automatic upgrades are enabled
	if features.GetAutomaticUpgrades() {
		if _, ok := c[defaultChannelName]; !ok {
			c[defaultChannelName] = defaultChannel
		}
	}

	var errs []error
	for name, channel := range c {
		// Wrapping is not mandatory here, but it adds the channel name in the
		// error, which makes troubleshooting easier.
		err = trace.Wrap(channel.CheckAndSetDefaults(), "failed to create channel %s", name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return trace.NewAggregate(errs...)
}

// DefaultVersion returns the version served by the default upgrade channel.
func (c Channels) DefaultVersion(ctx context.Context) (string, error) {
	channel, ok := c[defaultChannelName]
	if !ok {
		return "", trace.NotFound("default version channel not found")
	}
	targetVersion, err := channel.GetVersion(ctx)
	return targetVersion, trace.Wrap(err)
}

// Channel describes an automatic update channel configuration.
// It can be configured to serve a static version, or forward version requests
// to an upstream version server. Forwarded results are cached for 1 minute.
type Channel struct {
	// ForwardURL is the URL of the upstream version server providing the channel version/criticality.
	ForwardURL string `yaml:"forward_url,omitempty"`
	// StaticVersion is a static version the channel must serve. With or without the leading "v".
	StaticVersion string `yaml:"static_version,omitempty"`
	// Critical is whether the static version channel should be marked as a critical update.
	Critical bool `yaml:"critical"`

	// versionGetter gets the version of the channel. It is populated by CheckAndSetDefaults.
	versionGetter version.Getter
	// criticalTrigger gets the criticality of the channel. It is populated by CheckAndSetDefaults.
	criticalTrigger maintenance.Trigger
}

// CheckAndSetDefaults checks that the Channel configuration is valid and inits
// the version getter and maintenance trigger of the Channel based on its
// configuration. This function must be called before using the channel.
func (c *Channel) CheckAndSetDefaults() error {
	switch {
	case c.ForwardURL != "" && (c.StaticVersion != "" || c.Critical):
		return trace.BadParameter("cannot set both ForwardURL and (StaticVersion or Critical)")
	case c.ForwardURL != "":
		baseURL, err := url.Parse(c.ForwardURL)
		if err != nil {
			return trace.Wrap(err)
		}
		c.versionGetter = version.NewBasicHTTPVersionGetter(baseURL)
		c.criticalTrigger = maintenance.NewBasicHTTPMaintenanceTrigger("remote", baseURL)
	case c.StaticVersion != "":
		c.versionGetter = version.NewStaticGetter(c.StaticVersion, nil)
		c.criticalTrigger = maintenance.NewMaintenanceStaticTrigger("remote", c.Critical)
	default:
		return trace.BadParameter("either ForwardURL or StaticVersion must be set")
	}
	return nil
}

// GetVersion returns the current version of the channel. If io is involved,
// this function implements cache and is safe to call frequently.
func (c *Channel) GetVersion(ctx context.Context) (string, error) {
	return c.versionGetter.GetVersion(ctx)
}

// GetCritical returns the current criticality of the channel. If io is involved,
// this function implements cache and is safe to call frequently.
func (c *Channel) GetCritical(ctx context.Context) (bool, error) {
	return c.criticalTrigger.CanStart(ctx, nil)
}

// NewDefaultChannel creates a default automatic upgrade channel
// It looks up the environment variable, and if not found uses the default
// base URL. This default channel can be used in the proxy (to back its own version server)
// or in other Teleport process such as integration services deploying and
// updating teleport agents.
func NewDefaultChannel() (*Channel, error) {
	forwardURL := GetChannel()
	if forwardURL == "" {
		forwardURL = stableCloudVersionBaseURL
	}
	defaultChannel := &Channel{
		ForwardURL: forwardURL,
	}
	if err := defaultChannel.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}
	return defaultChannel, nil
}
