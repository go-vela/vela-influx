// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for InfluxDB information.
type Config struct {
	// the address to the Influx instance
	Addr string
	// the database to send data on Influx instance
	Database string
	// the user password for authenticating to Influx instance
	Password string
	// the user name for authenticating to Influx instance
	Username string
}

// New creates an Influx client for reporting metrics.
func (c *Config) New() (influx.Client, error) {
	logrus.Trace("creating new influx client from plugin configuration")

	// create config for InfluxDB client
	//
	// https://pkg.go.dev/github.com/influxdata/influxdb1-client/v2#HTTPConfig
	conf := influx.HTTPConfig{
		Addr:     c.Addr,
		Password: c.Password,
		Username: c.Username,
	}

	// create new InfluxDB client
	//
	// https://pkg.go.dev/github.com/influxdata/influxdb1-client/v2#NewHTTPClient
	client, err := influx.NewHTTPClient(conf)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config configuration")

	// verify influx address is provided
	if len(c.Addr) == 0 {
		return fmt.Errorf("no config addr provided")
	}

	// verify influx database is provided
	if len(c.Database) == 0 {
		return fmt.Errorf("no config database provided")
	}

	return nil
}
