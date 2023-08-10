package state

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
	"fuzz.codes/fuzzercloud/workerengine/podman"
	"fuzz.codes/fuzzercloud/workerengine/schemas"
)

const FILES_PREFIX = "/files/"

var Tasks map[string]*schemas.Task

func injectVariable(c []string, p string, v string) {
	for i, slice := range c {
		if strings.Contains(slice, p) {
			c[i] = strings.ReplaceAll(c[i], p, v)
		}
	}
}

func copyTaskFilesIntoContainer(t *schemas.Task, cid string) error {
	reader := b64.NewDecoder(b64.StdEncoding, strings.NewReader(t.Files))
	return podman.CopyTarIntoContainer(t.ID, reader, FILES_PREFIX)
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

func ResetTasks() {
	Tasks = make(map[string]*schemas.Task)
}

func NewTask(req schemas.CreateTaskRequest) (*schemas.Task, error) {
	tool, ok := Tools[req.ToolID]
	if !ok {
		return nil, errors.New("could not find tool")
	}
	t := &schemas.Task{
		Command: make([]string, 0),
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

	modifierTokens := strings.Split(modifier, " ")
	t.Command = append(t.Command, modifierTokens...)

	for _, iovar := range tool.Inputs {
		for k, v := range req.Inputs {
			if iovar.Name == k {
				if iovar.Type == tsf.STRING {
					injectVariable(t.Command, fmt.Sprint("{", k, "}"), v)
				} else if iovar.Type == tsf.FILE {
					injectVariable(t.Command, fmt.Sprint("{", k, "}"), fmt.Sprint(FILES_PREFIX, v))
				}
				break
			}
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

	err = copyTaskFilesIntoContainer(t, c.ID)
	if err != nil {
		return c.ID, err
	}

	err = UpdateTask(t)
	if err != nil {
		return t.ID, err
	}

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
