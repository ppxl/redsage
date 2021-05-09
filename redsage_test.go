package main

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func Test_doRun(t *testing.T) {
	t.Run("should join 2 pipelines and break them up with lunch breaks", func(t *testing.T) {
		file, _ := ioutil.TempFile(os.TempDir(), "redmineCSV-")
		path := file.Name()
		defer os.Remove(path)
		csvWriter := bufio.NewWriter(file)
		_, _ = csvWriter.Write([]byte(`Anforderungspipeline;2021-05-03;2021-05-04;2021-05-05;2021-05-06;Gesamtzeit
Pipeline A;7,50;6,00;"";4,50;18,00
Pipeline B/2;1,50;2,00;"";3,50;6,00
Gesamtzeit;7,50;6,00;"";4,50;18,00
`))
		args := runArgs{
			lunchBreakInMin:  60,
			singlePipelines:  []string{},
			filename:         path,
			csvDelimiter:     ";",
			decimalDelimiter: ",",
			skipColumnNames:  []string{"Gesamtzeit"},
			skipSummaryLine:  true,
		}

		// when
		actual := doRun(args)

		// then
		require.NoError(t, actual)
	})
}
