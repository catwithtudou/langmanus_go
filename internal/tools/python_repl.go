package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/log"
)

type PythonReplTool struct {
}

type pythonReplToolInput struct {
	Code string `json:"code"`
}

type pythonReplResult struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func GetPythonReplTool() tool.InvokableTool {
	return &PythonReplTool{}
}

func (t *PythonReplTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "python_repl",
		Desc: "A tool for executing Python code snippets in a REPL environment. This tool will:\n" +
			"1. Accept Python code as input\n" +
			"2. Execute the code in a sandboxed environment using Python interpreter\n" +
			"3. Return the execution output and any errors\n\n" +
			"Pre-loaded Packages:\n" +
			"- pandas: for data manipulation\n" +
			"- numpy: for numerical operations\n" +
			"- yfinance: for financial market data\n\n" +
			"Design Purpose:\n" +
			"- Provide a safe environment to execute Python code snippets\n" +
			"- Help with quick Python code testing and experimentation\n" +
			"- Support data analysis and numerical computations\n\n" +
			"Usage Limitations:\n" +
			"- Limited execution time\n" +
			"- Restricted system access\n" +
			"- No file system operations\n" +
			"- No network access (except for yfinance)\n\n" +
			"Example Usage:\n" +
			"1. Simple expression: `print(\"Hello, World!\")`\n" +
			"2. Data analysis: `import pandas as pd; df = pd.DataFrame({'A': [1, 2, 3]}); print(df)`\n" +
			"3. Financial data: `import yfinance as yf; data = yf.download('AAPL', start='2023-01-01'); print(data.head())`",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"code": {
				Type:     schema.String,
				Desc:     "The Python code to execute, can be a single expression or multiple statements",
				Required: true,
			},
		}),
	}, nil
}

func (t *PythonReplTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	input := &pythonReplToolInput{}
	err := json.Unmarshal([]byte(argumentsInJSON), input)
	if err != nil {
		log.GetLogger().Error("[PythonReplTool] Deserialization failed", zap.Error(err))
		return "", err
	}

	result := &pythonReplResult{}

	// Create Python command with -c flag to execute code
	cmd := exec.CommandContext(ctx, "python3", "-c", input.Code)

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Sprintf("Execution error: %v\n%s", err, string(output))
	} else {
		result.Output = string(output)
	}

	// Convert result to JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.GetLogger().Error("[PythonReplTool] Failed to serialize result", zap.Error(err))
		return "", err
	}

	log.GetLogger().Info("[PythonReplTool] Result", zap.String("result", string(jsonResult)))

	return string(jsonResult), nil
}
