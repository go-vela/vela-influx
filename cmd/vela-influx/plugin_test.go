// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
)

func TestInflux_Plugin_Exec(_ *testing.T) {
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
