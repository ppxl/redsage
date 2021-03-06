package transformer

import (
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/core"
)

// Config contains configuration values that modify the joining behaviour.
type Config struct {
	SinglePipelineNames []string
}

type Transformer interface {
	Transform(pdata *core.PipelineData, config Config) (*core.PipelineData, error)
}

func New() *joinTransformer {
	return &joinTransformer{}
}

type joinTransformer struct {
}

func (j *joinTransformer) Transform(pdata *core.PipelineData, config Config) (*core.PipelineData, error) {
	result := core.NewPipelineData()
	var err error
	var joinedPipeline *core.RedmineWorkPerDay

	for redminePipeline, workPerDay := range pdata.NamedDayRedmineValues {
		if joinedPipeline == nil {
			pipelineName := string(redminePipeline) + "-joined"
			joinedPipeline, err = result.AddPipeline(pipelineName)

			if err != nil {
				return nil, errors.Wrap(err, "error while crunching time data")
			}
		}

		for date, worktime := range workPerDay.WorkPerDay {
			joinedPipeline.PutWorkTime(date, worktime)
		}
	}

	return result, nil
}
