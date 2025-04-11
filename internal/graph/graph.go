package graph

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"catwithtudou/langmanus_go/internal/constants"
)

func BuildGraph[I, O any](ctx context.Context) (compose.Runnable[I, O], error) {
	g := compose.NewGraph[map[string]any, *schema.Message]()

	_ = g.AddLambdaNode(constants.CoordinatorAgent, nil)

	_ = g.AddEdge(compose.START, constants.CoordinatorAgent)

	_ = g.AddBranch(constants.CoordinatorAgent)

	return nil, nil
}
