package main

import (
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/core"
	"github.com/ppxl/sagemine/cruncher"
	"github.com/ppxl/sagemine/logging"
	"github.com/ppxl/sagemine/reader"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	flagLunchBreakInMinutesLong  = "break"
	flagLunchBreakInMinutesShort = "b"
	flagSinglePipelinesLong      = "single"
	flagSinglePipelinesShort     = "s"
	flagCSVColumnDelimiterLong   = "csv-column-delimiter"
	flagCSVColumnDelimiterShort  = "c"
	flagDecimalDelimiterLong     = "decimal-delimiter"
	flagDecimalDelimiterShort    = "d"
	flagIgnoreSummaryLineLong    = "ignore-summary-line"
	flagIgnoreSummaryLineShort   = "i"
)

var (
	// Version of the application
	Version string
	log     = logrus.New()
)

type runArgs struct {
	lunchBreakInMin  int
	singlePipelines  []string
	filename         string
	csvDelimiter     string
	decimalDelimiter string
	skipColumnNames  []string
	skipSummaryLine  bool
}

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

func run() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "read Redmine work time data and convert them to Sage-compatible data",
		Action:    doCliRun,
		ArgsUsage: "redmine CSV file",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    flagLunchBreakInMinutesLong,
				Aliases: []string{flagLunchBreakInMinutesShort},
				Usage:   "lunch break time in minutes (optional)",
				Value:   60,
			},
			&cli.StringSliceFlag{
				Name:    flagSinglePipelinesLong,
				Aliases: []string{flagSinglePipelinesShort},
				Usage: "These pipelines will receive their own pipeline and will not be joint into a single pseudo-pipeline (optional). " +
					"All other pipelines will be merged into a single pseudo-pipeline.",
			},
			&cli.StringFlag{
				Name:    flagCSVColumnDelimiterLong,
				Aliases: []string{flagCSVColumnDelimiterShort},
				Usage:   "this delimiter will be used to parse CSV columns (optional)",
				Value:   ";",
			},
			&cli.StringFlag{
				Name:    flagDecimalDelimiterLong,
				Aliases: []string{flagDecimalDelimiterShort},
				Usage:   "Set the decimal delimiter if the decimals in the CSV export uses a different format than '2.75' (optional)",
				Value:   ".",
			},
			&cli.BoolFlag{
				Name:    flagIgnoreSummaryLineLong,
				Aliases: []string{flagIgnoreSummaryLineShort},
				Usage:   "Set if the last line in the CSV export should be included or nto (optional)",
				Value:   true,
			},
		},
	}
}

func doCliRun(cliCtx *cli.Context) error {
	filename := cliCtx.Args().First()
	lunchBreakInMin := cliCtx.Int(flagLunchBreakInMinutesLong)
	csvColumnDelimiter := cliCtx.String(flagCSVColumnDelimiterLong)
	decimalDelimiter := cliCtx.String(flagDecimalDelimiterLong)
	ignoreSummaryLine := cliCtx.Bool(flagIgnoreSummaryLineLong)
	singlePipelines := cliCtx.StringSlice(flagSinglePipelinesLong)

	args := runArgs{
		lunchBreakInMin:  lunchBreakInMin,
		singlePipelines:  singlePipelines,
		filename:         filename,
		csvDelimiter:     csvColumnDelimiter,
		decimalDelimiter: decimalDelimiter,
		skipColumnNames:  nil,
		skipSummaryLine:  ignoreSummaryLine,
	}

	return doRun(args)
}

func doRun(args runArgs) error {
	// read
	options := reader.Options{
		Type: reader.CSV,
		CSVOptions: reader.CSVOptions{
			Filename:              args.filename,
			CSVDelimiter:          ";",
			InputDecimalDelimiter: ",",
			SkipColumnNames:       []string{"Gesamtzeit"},
			SkipSummaryLine:       false,
		},
		APIOptions: reader.APIOptions{},
	}
	redmineReader := reader.New(options)

	data, err := redmineReader.Read()
	if err != nil {
		return errors.Wrapf(err, "error while reading from %s", args.filename)
	}

	// crunch
	crunchConfig := cruncher.Config{
		LunchBreakInMin:     args.lunchBreakInMin,
		SinglePipelineNames: []core.PipelineName{(core.PipelineName)("Pipeline A")},
	}
	crunch := cruncher.New()

	crunched, err := crunch.Crunch(data, crunchConfig)
	if err != nil {
		return errors.Wrapf(err, "error while crunching data")
	}

	//fake
	if crunched.NamedDaySageValues != nil {

	}

	return nil
}
