package cmd

import (
	"github.com/urfave/cli/v2"
)

const (
	flagLunchBreakInMinutesLong  = "break"
	flagLunchBreakInMinutesShort = "b"
)

// Run executes the merging and splitting values from a CSV file and prints the output in Sage-relatable manner.
func Run() *cli.Command {
	return &cli.Command{
		Name:   "backup",
		Usage:  "backup parts of CES",
		Action: doRun,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    flagLunchBreakInMinutesLong,
				Aliases: []string{flagLunchBreakInMinutesShort},
				Usage:   "lunch break time in minutes",
				Value:   60,
			},
		},
	}
}

// ListBackups lists all backups which can be found for the configuration given in the registry
func doRun(*cli.Context) error {

	return nil
}
