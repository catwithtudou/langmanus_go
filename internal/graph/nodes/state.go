package nodes

import (
	"context"

	"catwithtudou/langmanus_go/config"
)

const (
	Finish config.AgentType = "finish"
)

type State struct {
	teamMembers []config.AgentType

	next                 string
	fullPlan             string
	deepThinkingMode     bool
	searchBeforePlanning bool
}

func NewState(ctx context.Context) *State {
	return &State{
		teamMembers:          append(config.TeamMembers, Finish),
		next:                 "",
		fullPlan:             "",
		deepThinkingMode:     false,
		searchBeforePlanning: false,
	}
}
