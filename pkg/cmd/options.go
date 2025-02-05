// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package cmd

import (
	"time"

	"github.com/gardener/gardener-extension-shoot-dns-service/pkg/controller/config"
	"github.com/gardener/gardener-extension-shoot-dns-service/pkg/controller/healthcheck"
	"github.com/gardener/gardener-extension-shoot-dns-service/pkg/controller/lifecycle"
	"github.com/gardener/gardener-extension-shoot-dns-service/pkg/controller/replication"

	"github.com/gardener/gardener/extensions/pkg/controller/cmd"
	extensionshealthcheckcontroller "github.com/gardener/gardener/extensions/pkg/controller/healthcheck"
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSServiceOptions holds options related to the dns service.
type DNSServiceOptions struct {
	SeedID                string
	DNSClass              string
	ReplicateDNSProviders bool
	OwnerDNSActivation    bool
	config                *DNSServiceConfig
}

// HealthOptions holds options for health checks.
type HealthOptions struct {
	HealthCheckSyncPeriod time.Duration
	config                *HealthConfig
}

// AddFlags implements Flagger.AddFlags.
func (o *DNSServiceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.SeedID, "seed-id", "", "ID of the current cluster")
	fs.StringVar(&o.DNSClass, "dns-class", "garden", "DNS class used to filter DNS source resources in shoot clusters")
	fs.BoolVar(&o.ReplicateDNSProviders, "replicate-dns-providers", false, "enables replication of DNSProviders from shoot cluster to seed cluster")
	fs.BoolVar(&o.OwnerDNSActivation, "enable-owner-dns-activation", false, "enables DNS activation of the shootdns DNSOwner")
}

// AddFlags implements Flagger.AddFlags.
func (o *HealthOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.HealthCheckSyncPeriod, "healthcheck-sync-period", time.Second*30, "sync period for the health check controller")
}

// Complete implements Completer.Complete.
func (o *DNSServiceOptions) Complete() error {
	o.config = &DNSServiceConfig{
		SeedID:                o.SeedID,
		DNSClass:              o.DNSClass,
		ReplicateDNSProviders: o.ReplicateDNSProviders,
		OwnerDNSActivation:    o.OwnerDNSActivation,
	}
	return nil
}

// Complete implements Completer.Complete.
func (o *HealthOptions) Complete() error {
	o.config = &HealthConfig{HealthCheckSyncPeriod: metav1.Duration{Duration: o.HealthCheckSyncPeriod}}
	return nil
}

// Completed returns the decoded CertificatesServiceConfiguration instance. Only call this if `Complete` was successful.
func (o *DNSServiceOptions) Completed() *DNSServiceConfig {
	return o.config
}

// Completed returns the completed HealthOptions. Only call this if `Complete` was successful.
func (o *HealthOptions) Completed() *HealthConfig {
	return o.config
}

// DNSServiceConfig contains configuration information about the dns service.
type DNSServiceConfig struct {
	SeedID                string
	DNSClass              string
	ReplicateDNSProviders bool
	OwnerDNSActivation    bool
}

// Apply applies the DNSServiceOptions to the passed ControllerOptions instance.
func (c *DNSServiceConfig) Apply(cfg *config.DNSServiceConfig) {
	cfg.SeedID = c.SeedID
	cfg.DNSClass = c.DNSClass
	cfg.ReplicateDNSProviders = c.ReplicateDNSProviders
	cfg.OwnerDNSActivation = c.OwnerDNSActivation
}

// HealthConfig contains configuration information about the health check controller.
type HealthConfig struct {
	HealthCheckSyncPeriod metav1.Duration
}

// ApplyHealthCheckConfig applies the `HealthConfig` to the passed health configurtaion.
func (c *HealthConfig) ApplyHealthCheckConfig(config *healthcheckconfig.HealthCheckConfig) {
	config.SyncPeriod = c.HealthCheckSyncPeriod
}

// ControllerSwitches are the cmd.ControllerSwitches for the provider controllers.
func ControllerSwitches() *cmd.SwitchOptions {
	return cmd.NewSwitchOptions(
		cmd.Switch(lifecycle.Name, lifecycle.AddToManager),
		cmd.Switch(replication.Name, replication.AddToManager),
		cmd.Switch(extensionshealthcheckcontroller.ControllerName, healthcheck.RegisterHealthChecks),
	)
}
