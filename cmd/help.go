package main

import (
	"github.com/urfave/cli/v2"
)

func CommandHelp() *cli.Command {
	return &cli.Command{
		Usage: "show the list of commands or help for one command",
		Name:  "help",

		Action: func(clictx *cli.Context) error {
			return cli.ShowAppHelp(clictx)
		},
	}
}
