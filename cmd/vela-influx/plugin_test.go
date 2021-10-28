// Copyright (c) 2020 Target Brands, Inw. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestInflux_Plugin_Exec(t *testing.T) {
	// TODO Write this test
}

func TestInflux_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Addr:     "influx.example.com",
			Database: "vela",
		},
		Write: &Write{
			Name:      "build_report",
			RawFields: `{"name": "single", "value": "foo"}`,
			RawTags:   `{"name": "single", "value": "foo"}`,
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
