package reader

import (
	"bufio"
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const pipelineA = "Pipeline A"

func Test_csvReader_Read(t *testing.T) {
	t.Run("should return single line for german Remine CSV", func(t *testing.T) {
		file, _ := ioutil.TempFile(os.TempDir(), "redmineCSV-")
		path := file.Name()
		defer os.Remove(path)
		csvWriter := bufio.NewWriter(file)
		_, err := csvWriter.Write([]byte(`Anforderungspipeline;2021-05-03;2021-05-04;2021-05-05;2021-05-06;Gesamtzeit
Pipeline A;7,50;6,00;"";4,50;18,00
Gesamtzeit;7,50;6,00;"";4,50;18,00
`))
		assert.NoError(t, err)
		err = csvWriter.Flush()
		assert.NoError(t, err)

		sut := newCSVReader(CSVOptions{
			Filename:         path,
			CSVDelimiter:     ";",
			DecimalDelimiter: ",",
		})

		//when
		actual, err := sut.Read()

		// then
		require.NoError(t, err)
		require.Equal(t, 2, actual.Entries())

		expected := core.NewPipelineData()

		expectedEntry, err := expected.AddPipeline(pipelineA)
		require.NoError(t, err)
		expectedEntry.PutWorkTime("2021-05-03", 7.50)
		expectedEntry.PutWorkTime("2021-05-04", 6)
		expectedEntry.PutWorkTime("2021-05-05", 0)
		expectedEntry.PutWorkTime("2021-05-06", 4.50)
		expectedEntry.PutWorkTime("Gesamtzeit", 18)

		expectedSums, err := expected.AddPipeline("Gesamtzeit")
		require.NoError(t, err)
		expectedSums.PutWorkTime("2021-05-03", 7.50)
		expectedSums.PutWorkTime("2021-05-04", 6)
		expectedSums.PutWorkTime("2021-05-05", 0)
		expectedSums.PutWorkTime("2021-05-06", 4.50)
		expectedSums.PutWorkTime("Gesamtzeit", 18)
		assert.Equal(t, expected, actual)
	})

	t.Run("should cut away selected columns from german Remine CSV", func(t *testing.T) {
		file, _ := ioutil.TempFile(os.TempDir(), "redmineCSV-")
		path := file.Name()
		defer os.Remove(path)
		csvWriter := bufio.NewWriter(file)
		_, err := csvWriter.Write([]byte(`Anforderungspipeline;2021-05-03;2021-05-04;2021-05-05;2021-05-06;Gesamtzeit
Pipeline A;7,50;6,00;"";4,50;18,00
Gesamtzeit;7,50;6,00;"";4,50;18,00
`))
		assert.NoError(t, err)
		err = csvWriter.Flush()
		assert.NoError(t, err)

		sut := newCSVReader(CSVOptions{
			Filename:         path,
			CSVDelimiter:     ";",
			DecimalDelimiter: ",",
			SkipColumnNames:  []string{"Gesamtzeit"},
		})

		//when
		actual, err := sut.Read()

		// then
		require.NoError(t, err)
		require.Equal(t, 2, actual.Entries())
		require.Equal(t, actual.NamedDayRedmineValues[pipelineA].Days(), 4)

		expected := core.NewPipelineData()
		expectedEntry, _ := expected.AddPipeline(pipelineA)
		expectedEntry.PutWorkTime("2021-05-03", 7.50)
		expectedEntry.PutWorkTime("2021-05-04", 6)
		expectedEntry.PutWorkTime("2021-05-05", 0)
		expectedEntry.PutWorkTime("2021-05-06", 4.50)

		expectedSums, _ := expected.AddPipeline("Gesamtzeit")
		expectedSums.PutWorkTime("2021-05-03", 7.50)
		expectedSums.PutWorkTime("2021-05-04", 6)
		expectedSums.PutWorkTime("2021-05-05", 0)
		expectedSums.PutWorkTime("2021-05-06", 4.50)
		assert.Equal(t, expected, actual)
	})

	t.Run("should cut away selected lines from german Remine CSV", func(t *testing.T) {
		file, _ := ioutil.TempFile(os.TempDir(), "redmineCSV-")
		path := file.Name()
		defer os.Remove(path)
		csvWriter := bufio.NewWriter(file)
		_, err := csvWriter.Write([]byte(`Anforderungspipeline;2021-05-03;2021-05-04;2021-05-05;2021-05-06;Gesamtzeit
Pipeline A;7,50;6,00;"";4,50;18,00
Gesamtzeit;7,50;6,00;"";4,50;18,00
`))
		assert.NoError(t, err)
		err = csvWriter.Flush()
		assert.NoError(t, err)

		sut := newCSVReader(CSVOptions{
			Filename:         path,
			CSVDelimiter:     ";",
			DecimalDelimiter: ",",
			SkipSummaryLine:  true,
		})

		//when
		actual, err := sut.Read()

		// then
		require.NoError(t, err)
		require.Equal(t, 1, actual.Entries())
		require.Equal(t, actual.NamedDayRedmineValues[pipelineA].Days(), 5)

		expected := core.NewPipelineData()
		expectedEntry, _ := expected.AddPipeline(pipelineA)
		expectedEntry.PutWorkTime("2021-05-03", 7.50)
		expectedEntry.PutWorkTime("2021-05-04", 6)
		expectedEntry.PutWorkTime("2021-05-05", 0)
		expectedEntry.PutWorkTime("2021-05-06", 4.50)
		expectedEntry.PutWorkTime("Gesamtzeit", 18)

		assert.Equal(t, expected, actual)
	})
}

func Test_formatDecimal(t *testing.T) {
	t.Run("should replace german decimal", func(t *testing.T) {
		actual := formatDecimal("123,45", ",")
		assert.Equal(t, "123.45", actual)
	})
	t.Run("should replace empty decimal with 0", func(t *testing.T) {
		actual := formatDecimal("", ",")
		assert.Equal(t, "0.0", actual)
	})
}

func Test_buildSkipColumns(t *testing.T) {
	type args struct {
		headers           []string
		headerToBeSkipped []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "should skip no columns", args: args{headers: []string{"header 1", "header 2"}, headerToBeSkipped: []string{}}, want: []int{}},
		{name: "should skip column 0", args: args{headers: []string{"header 1", "header 2"}, headerToBeSkipped: []string{"header 1"}}, want: []int{0}},
		{name: "should skip column 1", args: args{headers: []string{"header 1", "header 2"}, headerToBeSkipped: []string{"header 2"}}, want: []int{1}},
		{name: "should skip both columns", args: args{headers: []string{"header 1", "header 2"}, headerToBeSkipped: []string{"header 1", "header 2"}}, want: []int{0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSkipColumns(tt.args.headers, tt.args.headerToBeSkipped); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSkipColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLastLine(t *testing.T) {
	data := [][]string{
		{"col1", "col2"},
		{"col1", "col2"},
	}
	require.False(t, isLastLine(0, data))
	require.True(t, isLastLine(1, data))
}
