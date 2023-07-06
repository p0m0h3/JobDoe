package state

import (
	b64 "encoding/base64"
	"errors"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
	"fuzz.codes/fuzzercloud/workerengine/podman"
	"fuzz.codes/fuzzercloud/workerengine/schemas"
)

var Tasks map[string]*schemas.Task

func injectVariables(c []string, p string, v string) {
	for i, slice := range c {
		if strings.Contains(slice, p) {
			c[i] = strings.ReplaceAll(c[i], p, v)
		}
	}
}

func copyTaskFilesToContainer(t *schemas.Task, cid string) error {
	for name, data := range t.Files {
		for _, v := range t.Tool.Spec.Inputs {
			if name == v.Name && v.Type == tsf.FILE {
				reader := b64.NewDecoder(b64.StdEncoding, strings.NewReader(data))
				err := podman.CopyFileIntoContainer(t.ID, reader, "/"+name)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
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
	tool, ok := Tools[req.ToolID]
	if !ok {
		return nil, errors.New("could not find tool")
	}
	t := &schemas.Task{
		Command: make([]string, 0),
		Stdin:   req.Stdin,
		Env:     req.Env,
		Tool: schemas.Tool{
			ID:   req.ToolID,
			Spec: tool,
		},
		Files: req.Files,
	}

	modifier, ok := tool.Exe.Modifiers[req.Modifier]
	if !ok {
		return t, errors.New("could not find modifier")
	}
	t.Command = append(t.Command, tool.Exe.Command)

	modifierTokens := strings.Split(modifier.String, " ")
	t.Command = append(t.Command, modifierTokens...)

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
	Tasks[t.ID] = t

	err = copyTaskFilesToContainer(t, c.ID)
	if err != nil {
		return "", err
	}

	UpdateTask(t)

	err = podman.StartContainer(t.ID)
	if err != nil {
		return t.ID, err
	}

	return t.ID, nil
}

func UpdateTask(t *schemas.Task) error {
	c, err := podman.InspectContainer(t.ID)
	if err != nil {
		return err
	}
	t.Status = c.State.Status
	return nil
}
