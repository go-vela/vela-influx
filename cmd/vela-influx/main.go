// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"time"

	"github.com/go-vela/vela-influx/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
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
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-influx"
	app.HelpName = "vela-influx"
	app.Usage = "Vela Influx plugin for for sending data to an InfluxDB"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "INFLUX_LOG_LEVEL"},
			FilePath: "/vela/parameters/influx/log_level,/vela/secrets/influx/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},

		// Config Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_ADDR", "INFLUX_ADDR"},
			FilePath: string("/vela/parameters/influx/addr,/vela/secrets/influx/addr"),
			Name:     "config.addr",
			Usage:    "Influx instance to communicate with",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_DATABASE", "INFLUX_DATABASE"},
			FilePath: string("/vela/parameters/influx/database,/vela/secrets/influx/database"),
			Name:     "config.database",
			Usage:    "name of database within Influx instance",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "INFLUX_PASSWORD"},
			FilePath: string("/vela/parameters/influx/password,/vela/secrets/influx/password"),
			Name:     "config.password",
			Usage:    "user password for communication with the Influx instance",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "INFLUX_USERNAME"},
			FilePath: string("/vela/parameters/influx/username,/vela/secrets/influx/username"),
			Name:     "config.username",
			Usage:    "user name for communication with the Influx instance",
		},

		// Write Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_NAME", "INFLUX_NAME"},
			FilePath: string("/vela/parameters/influx/name,/vela/secrets/influx/name"),
			Name:     "write.name",
			Usage:    "name of the metrics to be created",
			Value:    "build_metrics",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_FIELDS", "INFLUX_FIELDS"},
			FilePath: string("/vela/parameters/influx/name,/vela/secrets/influx/name"),
			Name:     "write.fields",
			Usage:    "set of fields to be created with the data point",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TAGS", "INFLUX_TAGS"},
			FilePath: string("/vela/parameters/influx/tags,/vela/secrets/influx/tags"),
			Name:     "write.tags",
			Usage:    "set of tags to be created with the data point",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/influx",
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
