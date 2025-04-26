package nodes

import (
	"fmt"

	"github.com/cloudwego/eino/compose"
)

type ComposeNodes struct {
	Coordinator     *compose.Lambda
	CoordinatorNode *CoordinatorNode
	Planner         *compose.Lambda
	PlannerNode     *PlannerNode
	Supervisor      *compose.Lambda
	SupervisorNode  *SupervisorNode
	Researcher      *compose.Lambda
	ResearcherNode  *ResearcherNode
}

func BuildNodes() (*ComposeNodes, error) {
	coordinator, coordinatorNode, err := CoordinatorCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build coordinator node: %w", err)
	}

	planner, plannerNode, err := PlannerCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build planner node: %w", err)
	}

	supervisor, supervisorNode, err := SupervisorCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build supervisor node: %w", err)
	}

	researcher, researcherNode, err := ResearcherCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build researcher node: %w", err)
	}

	return &ComposeNodes{
		Coordinator:     coordinator,
		CoordinatorNode: coordinatorNode,
		Planner:         planner,
		PlannerNode:     plannerNode,
		Supervisor:      supervisor,
		SupervisorNode:  supervisorNode,
		Researcher:      researcher,
		ResearcherNode:  researcherNode,
	}, nil
}
