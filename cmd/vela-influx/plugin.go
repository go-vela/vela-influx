// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

type (
	// Plugin struct represents fields user can present to plugin
	Plugin struct {
		// config arguments loaded for the plugin
		Config *Config
		// write arguments loaded for the plugin
		Write *Write
	}
)

// Exec formats and runs the commands for sending metrics to InfluxDB.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// create new Influx client from config configuration
	client, err := p.Config.New()
	if err != nil {
		return err
	}

	logrus.Info("writing metric")

	return p.Write.Exec(client, p.Config.Database)
}

// Validate function to validate plugin configuration
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// validate write configuration
	err = p.Write.Validate()
	if err != nil {
		return err
	}

	// convert raw data from plugin
	err = p.Write.Unmarshal()
	if err != nil {
		return err
	}

	return nil
}
