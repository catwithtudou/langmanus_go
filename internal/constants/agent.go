package constants

const (
	CoordinatorAgent = "coordinator"
	PlannerAgent     = "planner"
	SupervisorAgent  = "supervisor"
	ResearcherAgent  = "researcher"
	CoderAgent       = "coder"
	BrowserAgent     = "browser"
	ReporterAgent    = "reporter"
)

var AgentLLMap = map[string]LLMType{
	CoordinatorAgent: BasicLLM,
	PlannerAgent:     ReasoningLLM,
	SupervisorAgent:  BasicLLM,
	ResearcherAgent:  BasicLLM,
	CoderAgent:       BasicLLM,
	BrowserAgent:     VisionLLM,
	ReporterAgent:    BasicLLM,
}
