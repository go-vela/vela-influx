// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-influx/version"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	// Plugin Information
	cmd := cli.Command{
		Name:      "vela-influx",
		Usage:     "Vela Influx plugin for for sending data to an InfluxDB",
		Copyright: "Copyright 2021 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Version: v.Semantic(),
		Action:  run,
	}

	// Plugin Flags

	cmd.Flags = []cli.Flag{

		&cli.StringFlag{
			Name:  "log.level",
			Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value: "info",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("INFLUX_LOG_LEVEL"),
				cli.File("/vela/parameters/influx/log_level"),
				cli.File("/vela/secrets/influx/log_level"),
			),
		},

		// Config Flags

		&cli.StringFlag{
			Name:  "config.addr",
			Usage: "Influx instance to communicate with",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_ADDR"),
				cli.EnvVar("INFLUX_ADDR"),
				cli.File("/vela/parameters/influx/addr"),
				cli.File("/vela/secrets/influx/addr"),
			),
		},
		&cli.StringFlag{
			Name:  "config.database",
			Usage: "name of database within Influx instance",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_DATABASE"),
				cli.EnvVar("INFLUX_DATABASE"),
				cli.File("/vela/parameters/influx/database"),
				cli.File("/vela/secrets/influx/database"),
			),
		},
		&cli.StringFlag{
			Name:  "config.password",
			Usage: "user password for communication with the Influx instance",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_PASSWORD"),
				cli.EnvVar("INFLUX_PASSWORD"),
				cli.File("/vela/parameters/influx/password"),
				cli.File("/vela/secrets/influx/password"),
			),
		},
		&cli.StringFlag{
			Name:  "config.username",
			Usage: "user name for communication with the Influx instance",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_USERNAME"),
				cli.EnvVar("INFLUX_USERNAME"),
				cli.File("/vela/parameters/influx/username"),
				cli.File("/vela/secrets/influx/username"),
			),
		},

		// Write Flags

		&cli.StringFlag{
			Name:  "write.name",
			Usage: "name of the metrics to be created",
			Value: "build_metrics",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_NAME"),
				cli.EnvVar("INFLUX_NAME"),
				cli.File("/vela/parameters/influx/name"),
				cli.File("/vela/secrets/influx/name"),
			),
		},
		&cli.StringFlag{
			Name:  "write.fields",
			Usage: "set of fields to be created with the data point",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_FIELDS"),
				cli.EnvVar("INFLUX_FIELDS"),
				cli.File("/vela/parameters/influx/fields"),
				cli.File("/vela/secrets/influx/fields"),
			),
		},
		&cli.StringFlag{
			Name:  "write.tags",
			Usage: "set of tags to be created with the data point",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TAGS"),
				cli.EnvVar("INFLUX_TAGS"),
				cli.File("/vela/parameters/influx/tags"),
				cli.File("/vela/secrets/influx/tags"),
			),
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(_ context.Context, c *cli.Command) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-influx",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/influx",
		"registry": "https://hub.docker.com/r/target/vela-influx",
	}).Info("Vela Influx Plugin")

	// create the plugin
	p := &Plugin{
		Config: &Config{
			Addr:     c.String("config.addr"),
			Database: c.String("config.database"),
			Password: c.String("config.password"),
			Username: c.String("config.username"),
		},
		Write: &Write{
			Name:      c.String("write.name"),
			RawFields: c.String("write.fields"),
			RawTags:   c.String("write.tags"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
