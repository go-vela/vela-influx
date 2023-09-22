// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
)

func TestInflux_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Addr:     "influx.example.com",
		Database: "vela",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestInflux_Config_Validate_NoAddr(t *testing.T) {
	// setup types
	c := &Config{
		Database: "vela",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestInflux_Config_Validate_NoDatabase(t *testing.T) {
	// setup types
	c := &Config{
		Addr: "influx.example.com",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
