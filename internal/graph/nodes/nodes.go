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
	Browser         *compose.Lambda
	BrowserNode     *BrowserNode
	Coder           *compose.Lambda
	CoderNode       *CoderNode
	Reporter        *compose.Lambda
	ReporterNode    *ReporterNode
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

	browser, browserNode, err := BrowserCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build browser node: %w", err)
	}

	coder, coderNode, err := CoderCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build coder node: %w", err)
	}

	reporter, reporterNode, err := ReporterCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build reporter node: %w", err)
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
		Browser:         browser,
		BrowserNode:     browserNode,
		Coder:           coder,
		CoderNode:       coderNode,
		Reporter:        reporter,
		ReporterNode:    reporterNode,
	}, nil
}
