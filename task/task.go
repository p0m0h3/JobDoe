package task

import (
	"errors"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
	"fuzz.codes/fuzzercloud/workerengine/container"
	"fuzz.codes/fuzzercloud/workerengine/tool"
)

type Task struct {
	ID         string
	Image      string
	Command    []string
	EnvVarList map[string]string
	Stdin      string
	Tool       *tsf.Tool
}

var Tasks map[string]*Task = make(map[string]*Task)

func injectVariables(c []string, p string, v string) {
	for i, slice := range c {
		if slice == p {
			c[i] = v
			return
		}
	}
}

func NewTask(req CreateTaskRequest) (*Task, error) {
	tool := tool.Tools[req.ToolName]
	t := &Task{
		Image:      tool.Name,
		Command:    make([]string, 0),
		Stdin:      req.Stdin,
		EnvVarList: req.EnvVarList,
		Tool:       tool,
	}

	modifier, ok := tool.Exe.Modifiers[req.Modifier]
	if !ok {
		return t, errors.New("could not find modifier")
	}

	t.Command = append(t.Command, tool.Exe.Command, modifier.String)

	for _, varPlaceholder := range modifier.Variables {
		variableName := strings.Trim(varPlaceholder, "{}")
		found := false
		for k, v := range req.InputList {
			if variableName == k {
				injectVariables(t.Command, varPlaceholder, v)
				found = true
				break
			}
		}
		if !found {
			return t, errors.New("could not satisfy command variables")
		}
	}

	Tasks[t.ID] = t
	return t, nil
}

func (t *Task) Start() (string, error) {
	return container.CreateContainer(t.Image, t.Command, t.EnvVarList)
}

func (t *Task) Refresh() error {
	return nil
}
