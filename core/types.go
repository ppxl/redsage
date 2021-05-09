package core

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
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

	pipeline := newRedmineWorkPerDay()
	pd.NamedDayRedmineValues[(PipelineName)(pipelineName)] = pipeline

	return pipeline, nil
}

// RedmineWorkPerDay maps a date string to the accumulated amount of time spent, f. i. 2021-05-05 -> 5.75
type RedmineWorkPerDay struct {
	WorkPerDay map[string]float64
}

func newRedmineWorkPerDay() *RedmineWorkPerDay {
	values := make(map[string]float64, 0)
	return &RedmineWorkPerDay{WorkPerDay: values}
}

func (rwpd *RedmineWorkPerDay) Days() int {
	return len(rwpd.WorkPerDay)
}

func (rwpd *RedmineWorkPerDay) PutWorkTime(date string, workTime float64) {
	currentWorkTime := rwpd.WorkTime(date)
	currentWorkTime += workTime

	rwpd.WorkPerDay[date] = currentWorkTime
}

func (rwpd *RedmineWorkPerDay) WorkTime(date string) float64 {
	return rwpd.WorkPerDay[date]
}

// CrunchedOutput contains mappings from pipeline name to Sage compatible work time
type CrunchedOutput struct {
	NamedDaySageValues map[PipelineName]*SageWorkPerDay
}

func (co *CrunchedOutput) String() string {
	result := ""

	sortedKeys := co.SortedKeys()
	for _, dayValue := range sortedKeys {
		logrus.Info("String: dayValue: " + dayValue)
		result += fmt.Sprintf("p: %s", co.NamedDaySageValues[(PipelineName)(dayValue)])
	}

	return result
}

func (co *CrunchedOutput) SortedKeys() []string {
	keys := make([]string, 0, len(co.NamedDaySageValues))
	for k := range co.NamedDaySageValues {
		keys = append(keys, (string)(k))
	}
	sort.Strings(keys)

	return keys
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
type SageWorkPerDay map[string][]TimeSlot

func (swpd *SageWorkPerDay) Days() int {
	return len(*swpd)
}

func (swpd *SageWorkPerDay) PutTimeSlot(date string, slotStart, slotEnd string) {
	timeSlots := swpd.TimeSlots(date)

	timeSlot := TimeSlot{}
	timeSlot.Start = slotStart
	timeSlot.End = slotEnd

	timeSlots = append(timeSlots, timeSlot)

	(*swpd)[date] = timeSlots
}

func (swpd *SageWorkPerDay) String() string {
	result := ""
	sortedKeys := swpd.SortedKeys()
	for _, timeSlot := range sortedKeys {
		result += fmt.Sprintf("d: ts: %s\t", timeSlot)
	}

	return result
}

func (swpd *SageWorkPerDay) SortedKeys() []string {
	keys := make([]string, 0, len(*swpd))
	for k := range *swpd {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (swpd *SageWorkPerDay) TimeSlots(day string) []TimeSlot {
	return (*swpd)[day]
}

func (swpd *SageWorkPerDay) PutEmptyTimeSlot(day string) {
	swpd.PutTimeSlot(day, "-", "-")
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

// DayTimeCounter holds the state of timeslots per day regardless of project.
type DayTimeCounter struct {
	counters      map[string]time.Time
	workStartTime string
}

func NewDayTimeCounter(workStartTime string) *DayTimeCounter {
	counters := make(map[string]time.Time, 0)
	return &DayTimeCounter{counters: counters, workStartTime: workStartTime}
}

func (dtc *DayTimeCounter) GetEndTimeOrDefault(date, defaultStartTime string) time.Time {
	result, ok := dtc.counters[date]
	if !ok {
		goodMorning, err := ParseDateWithTime(date, defaultStartTime)
		if err != nil {
			panic("could not get end time for day: " + err.Error())
		}
		dtc.counters[date] = goodMorning
		return goodMorning
	}

	return result
}
func (dtc *DayTimeCounter) EndTime(date string, endTime time.Time) {
	dtc.counters[date] = endTime
}
