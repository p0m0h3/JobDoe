package task

import (
	"errors"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
	"fuzz.codes/fuzzercloud/workerengine/podman"
	"fuzz.codes/fuzzercloud/workerengine/schemas"
)

type Task struct {
	ID      string            `json:"id"`
	Command []string          `json:"cmd"`
	Env     map[string]string `json:"env"`
	Stdin   string            `json:"stdin"`
	Status  string            `json:"status"`
	ToolID  string            `json:"toolId"`
	Spec    *tsf.Tool         `json:"tool"`
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

func NewTask(req schemas.CreateTaskRequest) (*Task, error) {
	tool := Tools[req.ToolID]
	t := &Task{
		Command: make([]string, 0),
		Stdin:   req.Stdin,
		Env:     req.Env,
		ToolID:  req.ToolID,
		Spec:    tool,
	}

	modifier, ok := tool.Exe.Modifiers[req.Modifier]
	if !ok {
		return t, errors.New("could not find modifier")
	}

	t.Command = append(t.Command, tool.Exe.Command, modifier.String)

	for _, varPlaceholder := range modifier.Variables {
		variableName := strings.Trim(varPlaceholder, "{}")
		found := false
		for k, v := range req.Inputs {
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

	return t, nil
}

func (t *Task) Start() (string, error) {
	id, err := podman.CreateContainer(t.Spec.Name, t.Command, t.Env)
	if err != nil {
		return "", err
	}
	t.ID = id

	Tasks[t.ID] = t
	return id, nil
}
