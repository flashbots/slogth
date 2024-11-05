package config

import (
	"time"
)

type Config struct {
	Delay         time.Duration
	DropThreshold int
	Stderr        bool
}

func (c *Config) Preprocess() error {
	return nil
}
