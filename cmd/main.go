package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/flashbots/slogth/config"
	"github.com/flashbots/slogth/slogth"
)

var (
	version = "development"
)

const (
	envPrefix       = "SLOGTH_"
	categoryGeneral = "GENERAL"
)

func main() {
	cfg := &config.Config{}

	flags := []cli.Flag{
		&cli.DurationFlag{
			Aliases:     []string{"d"},
			Category:    categoryGeneral,
			Destination: &cfg.Delay,
			EnvVars:     []string{envPrefix + "DELAY"},
			Name:        "delay",
			Usage:       "delay ingested logs by specified `duration`",
		},

		&cli.IntFlag{
			Category:    categoryGeneral,
			Destination: &cfg.DropThreshold,
			EnvVars:     []string{envPrefix + "DROP_THRESHOLD"},
			Name:        "drop-threshold",
			Usage:       "`count` of in-flight messages at which slogth should start dropping them (rate-limit)",
			Value:       0,
		},

		&cli.BoolFlag{
			Aliases:     []string{"e"},
			Category:    categoryGeneral,
			Destination: &cfg.Stderr,
			EnvVars:     []string{envPrefix + "STDERR"},
			Name:        "stderr",
			Usage:       "use stderr for output (stdout is used by default)",
			Value:       false,
		},
	}

	app := &cli.App{
		Name:    "slogth",
		Usage:   "delayed logs emission",
		Version: version,

		Flags: flags,

		HideHelpCommand: true,

		Commands: []*cli.Command{CommandHelp()},

		Before: func(_ *cli.Context) error {
			return cfg.Preprocess()
		},

		Action: func(_ *cli.Context) error {
			return slogth.Run(cfg)
		},
	}

	defer func() {
		zap.L().Sync() //nolint:errcheck
	}()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\nFailed with error:\n\n%s\n\n", err.Error())
		os.Exit(1)
	}
}
