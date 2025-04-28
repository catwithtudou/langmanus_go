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

// BuildGraph constructs the computation graph
func BuildGraph(ctx context.Context) (compose.Runnable[map[string]any, *schema.Message], error) {
	g := compose.NewGraph[map[string]any, *schema.Message](compose.WithGenLocalState(nodes.NewState))

	nodes, err := nodes.BuildNodes()
	if err != nil {
		log.GetLogger().Error("[BuildGraph]failed to build nodes", zap.Error(err))
		return nil, err
	}

	_ = g.AddLambdaNode(string(config.CoordinatorAgent), nodes.Coordinator)
	_ = g.AddLambdaNode(string(config.PlannerAgent), nodes.Planner)
	_ = g.AddLambdaNode(string(config.SupervisorAgent), nodes.Supervisor)
	_ = g.AddLambdaNode(string(config.ResearcherAgent), nodes.Researcher)

	_ = g.AddEdge(compose.START, string(config.CoordinatorAgent))

	_ = g.AddBranch(string(config.CoordinatorAgent), compose.NewGraphBranch(
		nodes.CoordinatorNode.Branch,
		nodes.CoordinatorNode.BranchNodes()))

	_ = g.AddBranch(string(config.PlannerAgent), compose.NewGraphBranch(
		nodes.PlannerNode.Branch,
		nodes.PlannerNode.BranchNodes()))

	_ = g.AddBranch(string(config.SupervisorAgent), compose.NewGraphBranch(
		nodes.SupervisorNode.Branch,
		nodes.SupervisorNode.BranchNodes()))

	_ = g.AddEdge(string(config.ResearcherAgent), string(config.SupervisorAgent))

	return g.Compile(ctx)
}
