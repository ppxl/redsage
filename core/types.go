package core

import "fmt"

const timeSlotFormat = "%s - %s"

// Config contains configuration values that modify the number crunching behaviour.
type Config struct {
	LunchBreakInMin     int
	SinglePipelineNames []PipelineName
}

func NewPipelineData() *PipelineData {
	values := make(map[PipelineName]RedmineWorkPerDay, 0)
	return &PipelineData{NamedDayRedmineValues: values}
}

// PipelineData contain parsed Redmine Pipeline data. Redmine divides hour in a decimal way, i. e. 0.5 means 30 minutes.
type PipelineData struct {
	// NamedDayValues maps the pipeline name to actual values per day, f. i. Pipeline 1 -> 2021-05-05 -> 5.75
	NamedDayRedmineValues map[PipelineName]RedmineWorkPerDay
}

func (pd PipelineData) Entries() int {
	return len(pd.NamedDayRedmineValues)
}

// RedmineWorkPerDay maps a date string to the accumulated amount of time spent, f. i. 2021-05-05 -> 5.75
type RedmineWorkPerDay map[string]float64

func (rwpd RedmineWorkPerDay) Days() int {
	return len(rwpd)
}

// CrunchedOutput contains mappings from pipeline name to Sage compatible work time
type CrunchedOutput struct {
	NamedDaySageValues map[PipelineName]SageWorkPerDay
}

// SageWorkPerDay maps a date string to a simplified Sage time slow, f. i. 2021-05-05 -> 13:00 - 14:00
type SageWorkPerDay map[string]TimeSlot

// TimeSlot represents a dateless wall clock interval of work, f. i. from 13:00 till 14:15
type TimeSlot struct {
	// Start contains the time slot's ending time in 24-hour format, f. i. "13:00" for 1 pm
	Start string
	// End contains the time slot's ending time in 24-hour format, f. i. "13:00" for 1 pm
	End string
}

func (t *TimeSlot) String() string {
	return fmt.Sprintf(timeSlotFormat, t.Start, t.End)
}

// PipelineName contains the name of a pipeline.
type PipelineName string
