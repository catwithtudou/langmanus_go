package nodes

import (
	"context"

	"catwithtudou/langmanus_go/config"
)

const (
	Finish config.AgentType = "FINISH"
)

type State struct {
	teamMembers []config.AgentType

	next                 string
	fullPlan             string
	deepThinkingMode     bool
	searchBeforePlanning bool
}

// Router is a struct used to define the next step routing
type Router struct {
	Next string `json:"next"`
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
