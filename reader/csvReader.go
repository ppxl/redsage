package reader

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/core"
	"github.com/ppxl/sagemine/logging"
	"os"
	"strconv"
	"strings"
)

const (
	CSV = iota
	RestAPI
)

var log = logging.Logger()

type CSVOptions struct {
	Filename         string
	CSVDelimiter     string
	DecimalDelimiter string
	SkipColumnNames  []string
	SkipSummaryLine  bool
}

type APIOptions struct {
	RedmineURL      string
	RedmineUser     string
	RedminePassword string
}

type Options struct {
	Type       int
	CSVOptions CSVOptions
	APIOptions APIOptions
}

type RedmineDataReader interface {
	Read() (*core.PipelineData, error)
}

type csvReader struct {
	options CSVOptions
}

func New(options Options) *csvReader {
	switch options.Type {
	case CSV:
		return newCSVReader(options.CSVOptions)
	case RestAPI:
		fallthrough
	default:
		log.Panicf("unsupported Redmine reader type %d", options.Type)
	}
	return nil
}

func newCSVReader(options CSVOptions) *csvReader {
	return &csvReader{options: options}
}

func (cr *csvReader) Read() (*core.PipelineData, error) {
	fileReader, err := os.OpenFile(cr.options.Filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fileReader.Close()
	commaRunes := []rune(cr.options.CSVDelimiter)

	r := csv.NewReader(fileReader)
	r.Comma = commaRunes[0]
	r.Comment = '#'

	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	result := core.NewPipelineData()
	columnHeaders := []string{}
	columnsToSkip := []int{}

	for currentLine, line := range data {
		var pipeline core.PipelineName
		entry := core.RedmineWorkPerDay{}
		if currentLine == 0 {
			columnHeaders = line
			columnsToSkip = buildSkipColumns(columnHeaders, cr.options.SkipColumnNames)
			continue
		}

		if cr.options.SkipSummaryLine && isLastLine(currentLine, data) {
			break
		}

		for currentColumn, cell := range line {
			fmt.Printf("%s\t", cell)

			if currentColumn == 0 {
				pipeline = (core.PipelineName)(cell)
				continue
			}

			if skipColumn(currentColumn, columnsToSkip) {
				continue
			}

			workTimeRaw := formatDecimal(cell, cr.options.DecimalDelimiter)

			workTime, err := strconv.ParseFloat(workTimeRaw, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "could not cast value '%s' to float (line %d, column %d)", cell, currentLine, currentColumn)
			}

			currentDay := columnHeaders[currentColumn]

			entry[currentDay] = workTime

		}

		result.NamedDayRedmineValues[pipeline] = entry
		fmt.Println()
	}

	return result, nil
}

func isLastLine(currentLine int, data [][]string) bool {
	return currentLine == len(data)-1
}

func skipColumn(currentColumn int, skip []int) bool {
	for _, skipColumn := range skip {
		if currentColumn == skipColumn {
			return true
		}
	}

	return false
}

func buildSkipColumns(headers []string, headerToBeSkipped []string) []int {
	result := []int{}
	for index, header := range headers {
		for _, skip := range headerToBeSkipped {
			if header == skip {
				result = append(result, index)
			}
		}
	}

	return result
}

func formatDecimal(cell string, delimiter string) string {
	if cell == "" {
		return "0.0"
	}

	floatDecimalDelimiter := "."
	return strings.Replace(cell, delimiter, floatDecimalDelimiter, 1)
}
