package internal

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/internal/graph"
	"catwithtudou/langmanus_go/log"
)

func RunAgentWorkflow(ctx context.Context, input string) (output *schema.Message, err error) {

	graph, err := graph.BuildGraph(ctx)
	if err != nil {
		log.GetLogger().Error("[RunAgentWorkflow]failed to build graph", zap.Error(err))
		return nil, err
	}

	return graph.Invoke(ctx, map[string]any{
		"input": input,
	})
}
