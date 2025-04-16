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

// BuildGraph 构建计算图
func BuildGraph(ctx context.Context) (compose.Runnable[map[string]any, *schema.Message], error) {
	g := compose.NewGraph[map[string]any, *schema.Message]()

	nodes, err := nodes.BuildNodes()
	if err != nil {
		log.GetLogger().Error("[BuildGraph]failed to build nodes", zap.Error(err))
		return nil, err
	}

	_ = g.AddLambdaNode(string(config.CoordinatorAgent), nodes.Coordinator)
	_ = g.AddLambdaNode(string(config.PlannerAgent), nodes.Planner)

	_ = g.AddEdge(compose.START, string(config.CoordinatorAgent))
	_ = g.AddEdge(string(config.CoordinatorAgent), string(config.PlannerAgent))
	_ = g.AddEdge(string(config.PlannerAgent), compose.END)

	return g.Compile(ctx)
}
