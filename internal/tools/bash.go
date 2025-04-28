package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/log"
)

type BashTool struct {
}

type bashToolInput struct {
	Command string `json:"command"`
}

func GetBashTool() tool.InvokableTool {
	return &BashTool{}
}

func (t *BashTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "bash_tool",
		Desc: "A tool for executing bash commands",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"command": {
				Type:     schema.String,
				Desc:     "The bash command to execute",
				Required: true,
			},
		}),
	}, nil
}

func (t *BashTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	input := &bashToolInput{}
	err := json.Unmarshal([]byte(argumentsInJSON), input)
	if err != nil {
		log.GetLogger().Error("[BashTool] Deserialization failed", zap.Error(err))
		return "", err
	}

	log.GetLogger().Info("[BashTool] Executing command", zap.String("command", input.Command))

	// Create a new command with the input command
	cmd := exec.CommandContext(ctx, "bash", "-c", input.Command)

	// Capture both stdout and stderr
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		// If the command failed, return both stdout and stderr
		errorMsg := fmt.Sprintf("Command execution failed with exit code %d.\nStdout: %s\nStderr: %s",
			cmd.ProcessState.ExitCode(),
			stdout.String(),
			stderr.String())
		log.GetLogger().Error("[BashTool] Command execution failed", zap.Error(err), zap.String("stdout", stdout.String()), zap.String("stderr", stderr.String()))
		return errorMsg, nil
	}

	// Return the command output
	return stdout.String(), nil
}
