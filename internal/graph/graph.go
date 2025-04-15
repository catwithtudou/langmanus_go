package graph

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/graph/nodes"
	"catwithtudou/langmanus_go/log"
)

func BuildGraph[I, O any](ctx context.Context) (compose.Runnable[I, O], error) {
	g := compose.NewGraph[map[string]any, *schema.Message]()

	coordinator, err := nodes.CoordinatorCompose()
	if err != nil {
		log.GetLogger().Error("[BuildGraph]failed to build coordinator node", zap.Error(err))
		return nil, err
	}

	_ = g.AddLambdaNode(string(config.CoordinatorAgent), coordinator)

	_ = g.AddEdge(compose.START, string(config.CoordinatorAgent))

	// ...TODO: add other nodes

	return nil, nil
}
