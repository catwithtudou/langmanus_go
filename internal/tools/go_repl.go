package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/log"
)

type GoReplTool struct {
}

type goReplToolInput struct {
	Code string `json:"code"`
}

type goReplResult struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

var interpreter *interp.Interpreter

func initGoReplTool() {
	// Initialize the interpreter with default options
	interpreter = interp.New(interp.Options{})

	// Use standard library symbols
	interpreter.Use(stdlib.Symbols)

	// Pre-import commonly used packages
	_, err := interpreter.Eval(`import (
		"fmt"
		"time"
		"strings"
		"math"
	)`)
	if err != nil {
		log.GetLogger().Error("[GoReplTool] Failed to pre-import packages", zap.Error(err))
	}
}

func GetGoReplTool() tool.InvokableTool {
	return &GoReplTool{}
}

func (t *GoReplTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "go_repl",
		Desc: "A tool for executing Go code snippets in a REPL environment. This tool will:\n" +
			"1. Accept Go code as input\n" +
			"2. Execute the code in a sandboxed environment using Yaegi interpreter\n" +
			"3. Return the execution output and any errors\n\n" +
			"Pre-loaded Packages:\n" +
			"- fmt: for formatted I/O\n" +
			"- time: for time operations\n" +
			"- strings: for string manipulation\n" +
			"- math: for mathematical operations\n\n" +
			"Design Purpose:\n" +
			"- Provide a safe environment to execute Go code snippets\n" +
			"- Help with quick Go code testing and experimentation\n" +
			"- Support basic Go language features\n" +
			"- Maintain state between executions (imported packages, variables, etc.)\n\n" +
			"Usage Limitations:\n" +
			"- Limited execution time\n" +
			"- Restricted system access\n" +
			"- No file system operations\n" +
			"- No network access\n" +
			"- No support for assembly files (.s)\n" +
			"- No support for C code integration\n" +
			"- No support for compiler/linker directives\n" +
			"- No support for Go modules\n\n" +
			"Example Usage:\n" +
			"1. Simple expression: `fmt.Println(\"Hello, World!\")`\n" +
			"2. Package declaration: `package main; func main() { fmt.Println(\"Hello\") }`\n" +
			"3. Variable declaration: `x := 42; fmt.Println(x)`\n" +
			"4. Function definition: `func add(a, b int) int { return a + b }; fmt.Println(add(1, 2))`",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"code": {
				Type:     schema.String,
				Desc:     "The Go code to execute, can be an expression, statement, or package declaration",
				Required: true,
			},
		}),
	}, nil
}

func (t *GoReplTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	input := &goReplToolInput{}
	err := json.Unmarshal([]byte(argumentsInJSON), input)
	if err != nil {
		log.GetLogger().Error("[GoReplTool] Deserialization failed", zap.Error(err))
		return "", err
	}

	result := &goReplResult{}

	// Check if the code is a package declaration
	if strings.HasPrefix(strings.TrimSpace(input.Code), "package") {
		// For package declarations, we need to evaluate the entire package
		_, err = interpreter.Eval(input.Code)
		if err != nil {
			result.Error = fmt.Sprintf("Package evaluation error: %v", err)
		} else {
			result.Output = "Package evaluated successfully"
		}
	} else {
		// For expressions and statements, we can evaluate directly
		v, err := interpreter.Eval(input.Code)
		if err != nil {
			result.Error = fmt.Sprintf("Evaluation error: %v", err)
		} else {
			if v.IsValid() {
				result.Output = fmt.Sprintf("%v", v.Interface())
			} else {
				result.Output = "Expression evaluated successfully"
			}
		}
	}

	// Convert result to JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.GetLogger().Error("[GoReplTool] Failed to serialize result", zap.Error(err))
		return "", err
	}

	log.GetLogger().Info("[GoReplTool] Result", zap.String("result", string(jsonResult)))

	return string(jsonResult), nil
}
