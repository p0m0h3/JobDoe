package state

import (
	"errors"
	"strings"

	"fuzz.codes/fuzzercloud/workerengine/podman"
	"fuzz.codes/fuzzercloud/workerengine/schemas"
)

var Tasks map[string]*schemas.Task

func injectVariables(c []string, p string, v string) {
	for i, slice := range c {
		if slice == p {
			c[i] = v
			return
		}
	}
}

func ReadTasks() error {
	containers, err := podman.GetAllContainers()
	if err != nil {
		return err
	}

	for _, c := range containers {
		task := &schemas.Task{
			ID:      c.ID,
			Command: c.Command,
			Status:  c.Status,
		}
		Tasks[task.ID] = task
	}
	return nil
}

func NewTask(req schemas.CreateTaskRequest) (*schemas.Task, error) {
	tool := Tools[req.ToolID]
	t := &schemas.Task{
		Command: make([]string, 0),
		Stdin:   req.Stdin,
		Env:     req.Env,
		Tool: schemas.Tool{
			ID:   req.ToolID,
			Spec: tool,
		},
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

func StartTask(t *schemas.Task) (string, error) {
	c, err := podman.CreateContainer(t.Tool.Spec.Name, t.Command, t.Env)
	if err != nil {
		return "", err
	}

	t.ID = c.ID
	UpdateTask(t)

	Tasks[t.ID] = t
	return c.ID, nil
}

func UpdateTask(t *schemas.Task) error {
	c, err := podman.InspectContainer(t.ID)
	if err != nil {
		return err
	}
	t.Status = c.State.Status
	return nil
}
