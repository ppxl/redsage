package transformer

import (
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const pipelineAName = "Pipeline A"
const pipelineBName = "Pipeline 2/B"

func Test_joinTransformer_Transform(t *testing.T) {
	t.Run("should join pipelines", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA, _ := input.AddPipeline(pipelineAName)
		pipelineA.PutWorkTime("2021-05-03", 0)
		pipelineA.PutWorkTime("2021-05-04", 0)
		pipelineA.PutWorkTime("2021-05-05", 1)
		pipelineA.PutWorkTime("2021-05-06", 2)
		pipelineA.PutWorkTime("2021-05-07", 3.5)

		pipelineB, _ := input.AddPipeline(pipelineBName)
		pipelineB.PutWorkTime("2021-05-03", 0)
		pipelineB.PutWorkTime("2021-05-04", 1)
		pipelineB.PutWorkTime("2021-05-05", 2)
		pipelineB.PutWorkTime("2021-05-06", 1.5)
		pipelineB.PutWorkTime("2021-05-07", 0.5)
		sut := &joinTransformer{}

		// when
		actual, err := sut.Transform(input)

		// then
		require.NoError(t, err)

		expected := core.NewPipelineData()
		pipelineAJoined, _ := expected.AddPipeline("Pipeline A-joined")
		pipelineAJoined.PutWorkTime("2021-05-03", 0)
		pipelineAJoined.PutWorkTime("2021-05-04", 1)
		pipelineAJoined.PutWorkTime("2021-05-05", 3)
		pipelineAJoined.PutWorkTime("2021-05-06", 3.5)
		pipelineAJoined.PutWorkTime("2021-05-07", 4)
		assert.Equal(t, expected, actual)
	})
}
