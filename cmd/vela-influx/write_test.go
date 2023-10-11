// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
)

func TestInflux_Write_Exec(t *testing.T) {
	// TODO Write this test
}

func TestInflux_Write_Validate(t *testing.T) {
	// setup types
	w := &Write{
		Name:      "build_report",
		RawFields: `[{"name": "single", "value": "foo"}]`,
		RawTags:   `[{"name": "single", "value": "foo"}]`,
	}

	err := w.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestInflux_Write_Validate_NoFields(t *testing.T) {
	// setup types
	w := &Write{
		Name:    "build_report",
		RawTags: `[{"name": "single", "value": "foo"}]`,
	}

	err := w.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestInflux_Write_Validate_NoName(t *testing.T) {
	// setup types
	w := &Write{
		RawFields: `[{"name": "single", "value": "foo"}]`,
		RawTags:   `[{"name": "single", "value": "foo"}]`,
	}

	err := w.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestInflux_Write_Validate_NoTags(t *testing.T) {
	// setup types
	w := &Write{
		Name:      "build_report",
		RawFields: `[{"name": "single", "value": "foo"}]`,
	}

	err := w.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
