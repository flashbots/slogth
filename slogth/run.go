package slogth

import (
	"os"

	"github.com/flashbots/slogth/config"
)

func Run(cfg *config.Config) error {
	s := new()

	s.delay = cfg.Delay
	s.dropThreshold = cfg.DropThreshold
	s.input = os.Stdin

	if cfg.Stderr {
		s.output = os.Stderr
	} else {
		s.output = os.Stdout
	}

	return s.run()
}
