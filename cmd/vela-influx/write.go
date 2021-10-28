// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

type (
	// Write represents the plugin configuration for write information.
	Write struct {
		// name of the metrics to be created
		Name string
		// set of fields to be created with the data point
		Fields *Field
		// raw input of fields provided for plugin
		RawFields string
		// raw input of tags provided for plugin
		RawTags string
		// set of tags to be created with the data point
		Tags *Tag
	}

	// Field represents a set of tags to be created with the data point.
	Field struct {
		Data map[string]interface{}
	}

	// Tag represents a set of fields to be created with the data point.
	Tag struct {
		Data map[string]string
	}
)

// Exec formats and runs the commands for sending data to Influx.
func (w *Write) Exec(client influx.Client, database string) error {
	logrus.Trace("running delete with provided configuration")

	// create batch of data points to send to the database
	//
	// https://pkg.go.dev/github.com/influxdata/influxdb1-client/v2#NewBatchPoints
	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: database,
	})

	// create new data point to add to the batch
	//
	// https://pkg.go.dev/github.com/influxdata/influxdb1-client/v2#NewPoint
	pt, err := influx.NewPoint(w.Name, w.Tags.Data, w.Fields.Data, time.Now())
	if err != nil {
		return fmt.Errorf("unable to create point: %w", err)
	}

	// add the data point to the batch to send to the database
	bp.AddPoint(pt)

	// send the batch point to the database
	return client.Write(bp)
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (w *Write) Unmarshal() error {
	logrus.Trace("unmarshaling fields and tags")

	// cast raw properties into bytes
	bytesFields := []byte(w.RawFields)

	jsonFields := make(map[string]interface{})

	// serialize raw properties into expected Props type
	err := json.Unmarshal(bytesFields, &jsonFields)
	if err != nil {
		return err
	}

	w.Fields.Data = jsonFields

	// when tags are provided unmarshal them
	if len(w.RawTags) != 0 {
		// cast raw properties into bytes
		bytesTags := []byte(w.RawTags)

		jsonTags := make(map[string]string)

		// serialize raw properties into expected Props type
		err = json.Unmarshal(bytesTags, &jsonTags)
		if err != nil {
			return err
		}

		w.Tags.Data = jsonTags
	}

	return nil
}

// Validate verifies the Write is properly configured.
func (w *Write) Validate() error {
	logrus.Trace("validating write configuration")

	// verify Fields is provided
	if len(w.RawFields) == 0 {
		return fmt.Errorf("no write fields provided")
	}

	// verify Name is provided
	if len(w.Name) == 0 {
		return fmt.Errorf("no write name provided")
	}

	// default maps
	w.Fields = &Field{
		Data: make(map[string]interface{}),
	}

	w.Tags = &Tag{
		Data: make(map[string]string),
	}

	return nil
}
