package config

type AgentType string

const (
	CoordinatorAgent AgentType = "coordinator"
	PlannerAgent     AgentType = "planner"
	SupervisorAgent  AgentType = "supervisor"
	ResearcherAgent  AgentType = "researcher"
	CoderAgent       AgentType = "coder"
	BrowserAgent     AgentType = "browser"
	ReporterAgent    AgentType = "reporter"
)

var TeamMembers = []AgentType{ResearcherAgent, CoderAgent, BrowserAgent, ReporterAgent}

var AgentLLMap = map[AgentType]LLMType{
	CoordinatorAgent: BasicLLM,
	PlannerAgent:     ReasoningLLM,
	SupervisorAgent:  BasicLLM,
	ResearcherAgent:  BasicLLM,
	CoderAgent:       BasicLLM,
	BrowserAgent:     VisionLLM,
	ReporterAgent:    BasicLLM,
}
