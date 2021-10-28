// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	influx "github.com/influxdata/influxdb1-client/v2"

	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Kubernetes information.
type Config struct {
	// is the address to an influx server
	Addr string
	// is the database for interaction on the server
	Database string
	// is the user password for authenticating to server
	Password string
	// is the user name for authenticating to server
	Username string
}

// New creates an Influx client for reporting metrics.
func (c *Config) New() (influx.Client, error) {
	logrus.Trace("creating new influx client from plugin configuration")

	// create config for InfluxDB client
	conf := influx.HTTPConfig{
		Addr:     c.Addr,
		Password: c.Password,
		Username: c.Username,
	}

	// create new InfluxDB client
	client, err := influx.NewHTTPClient(conf)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	return client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config configuration")

	// verify Addr is provided
	if len(c.Addr) == 0 {
		return fmt.Errorf("no config addr provided")
	}

	// verify Database is provided
	if len(c.Database) == 0 {
		return fmt.Errorf("no config database provided")
	}

	return nil
}
