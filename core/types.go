package core

import (
	"errors"
	"fmt"
)

const timeSlotFormat = "%s - %s"

func NewPipelineData() *PipelineData {
	values := make(map[PipelineName]*RedmineWorkPerDay, 0)
	return &PipelineData{NamedDayRedmineValues: values}
}

// PipelineData contain parsed Redmine Pipeline data. Redmine divides hour in a decimal way, i. e. 0.5 means 30 minutes.
//
// Example:
//  data := NewPipelineData()
//  pipeline, err := data.AddPipeline("My Pipeline")
//  pipeline.PutWorkTime("2021-05-05", 7.5)
//  pipeline.PutWorkTime("2021-05-06", 4)
type PipelineData struct {
	// NamedDayValues maps the pipeline name to actual values per day, f. i. Pipeline 1 -> 2021-05-05 -> 5.75
	NamedDayRedmineValues map[PipelineName]*RedmineWorkPerDay
}

func (pd *PipelineData) Entries() int {
	return len(pd.NamedDayRedmineValues)
}

func (pd *PipelineData) AddPipeline(pipelineName string) (*RedmineWorkPerDay, error) {
	if pipelineName == "" {
		return nil, errors.New("pipeline name must not be empty")
	}

	pipeline := &RedmineWorkPerDay{}
	pd.NamedDayRedmineValues[(PipelineName)(pipelineName)] = pipeline

	return pipeline, nil
}

// RedmineWorkPerDay maps a date string to the accumulated amount of time spent, f. i. 2021-05-05 -> 5.75
type RedmineWorkPerDay map[string]float64

func (rwpd *RedmineWorkPerDay) Days() int {
	return len(*rwpd)
}

func (rwpd *RedmineWorkPerDay) PutWorkTime(date string, workTime float64) {
	(*rwpd)[date] = workTime
}

// CrunchedOutput contains mappings from pipeline name to Sage compatible work time
type CrunchedOutput struct {
	NamedDaySageValues map[PipelineName]*SageWorkPerDay
}

func (co *CrunchedOutput) String() string {
	result := ""

	for pipeline, dayValue := range co.NamedDaySageValues {
		result += fmt.Sprintf("p: %s, %s", pipeline, dayValue)
	}

	return result
}

func NewCrunchedOutput() *CrunchedOutput {
	values := make(map[PipelineName]*SageWorkPerDay, 0)
	return &CrunchedOutput{NamedDaySageValues: values}
}

func (co *CrunchedOutput) AddPipeline(pipelineName string) (*SageWorkPerDay, error) {
	if pipelineName == "" {
		return nil, errors.New("pipeline name must not be empty")
	}

	pipeline := &SageWorkPerDay{}
	co.NamedDaySageValues[(PipelineName)(pipelineName)] = pipeline

	return pipeline, nil
}

// SageWorkPerDay maps a date string to a simplified Sage time slow, f. i. 2021-05-05 -> 13:00 - 14:00
type SageWorkPerDay map[string]*TimeSlot

func (swpd *SageWorkPerDay) Days() int {
	return len(*swpd)
}

func (swpd *SageWorkPerDay) PutTimeSlot(date string, slotStart, slotEnd string) {
	timeSlot := swpd.timeSlot(date)
	if timeSlot == nil {
		timeSlot = &TimeSlot{}
	}

	timeSlot.Start = slotStart
	timeSlot.End = slotEnd

	(*swpd)[date] = timeSlot
}

func (swpd *SageWorkPerDay) String() string {
	result := ""
	for pipeline, timeSlot := range *swpd {
		result += fmt.Sprintf("d: %s, ts: %s", pipeline, timeSlot)
	}

	return result
}

func (swpd *SageWorkPerDay) timeSlot(day string) *TimeSlot {
	return (*swpd)[day]
}

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
