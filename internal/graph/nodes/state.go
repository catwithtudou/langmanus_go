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

// Router 用于定义下一步路由的结构体
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
