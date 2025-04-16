package nodes

import (
	"fmt"

	"github.com/cloudwego/eino/compose"
)

type ComposeNodes struct {
	Coordinator *compose.Lambda
	Planner     *compose.Lambda
}

func BuildNodes() (*ComposeNodes, error) {
	coordinator, err := CoordinatorCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build coordinator node: %w", err)
	}

	planner, err := PlannerCompose()
	if err != nil {
		return nil, fmt.Errorf("failed to build planner node: %w", err)
	}

	return &ComposeNodes{
		Coordinator: coordinator,
		Planner:     planner,
	}, nil
}
