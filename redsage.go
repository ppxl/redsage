package main

import (
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/logging"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	// Version of the application
	Version string

	log = logrus.New()
)

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "define log level",
			Value: "warning",
		},
	}
}

func configureLogging(cliCtx *cli.Context) error {
	logLevel := cliCtx.String("log-level")
	logLevelParsed, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return errors.Wrapf(err, "could not parse log level %s to logrus level", logLevel)
	}
	err = logging.Init(logLevelParsed)
	log = logging.Logger()
	if err != nil {
		return errors.Wrap(err, "could not initialize logging")
	}
	return nil
}

func configureApplication(cliCtx *cli.Context) error {
	err := configureLogging(cliCtx)
	if err != nil {
		return err
	}

	return nil
}

func checkMainError(err error) {
	if err != nil {
		println("%+s\n", err)
		os.Exit(1)
	}
}

// projects main function
func main() {
	app := cli.NewApp()
	app.Name = "redsage"
	app.Usage = "Maintain sanity while combining Redmine activity times and Sage project times"
	app.Version = Version
	app.Commands = []*cli.Command{run()}

	app.Flags = createGlobalFlags()
	app.Before = configureApplication
	err := app.Run(os.Args)
	checkMainError(err)
}

const (
	flagLunchBreakInMinutesLong  = "break"
	flagLunchBreakInMinutesShort = "b"
	flagSinglePipelinesLong      = "single"
	flagSinglePipelinesShort     = "s"
)

func run() *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  "read Redmine work time data and convert them to Sage-compatible data",
		Action: doRun,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    flagLunchBreakInMinutesLong,
				Aliases: []string{flagLunchBreakInMinutesShort},
				Usage:   "lunch break time in minutes",
				Value:   60,
			},
			&cli.StringSliceFlag{
				Name:    flagSinglePipelinesLong,
				Aliases: []string{flagSinglePipelinesShort},
				Usage:   "these pipelines will receive their own pipeline and will not be joint into a single pseudo-pipeline",
			},
		},
	}
}

func doRun(*cli.Context) error {

	return nil
}
