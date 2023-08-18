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

const (
	FILES_PREFIX  = "/files/"
	OUTPUT_PREFIX = "output/"
)

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
	return podman.CopyIntoContainer(t.ID, reader, FILES_PREFIX)
}

func findProfile(t *tsf.Tool, profName string) tsf.Profile {
	for _, prof := range t.Execute.Profiles {
		if prof.Name == profName {
			return prof
		}
	}

	return tsf.Profile{}
}

func findModifiers(t *tsf.Tool, modNames []string) []string {
	result := make([]string, 0)

	for _, modName := range modNames {
		for _, mod := range t.Execute.Modifiers {
			if modName == mod.Name {
				result = append(result, mod.Format)
			}
		}
	}

	return result
}

func formatVariable(input tsf.Input, value string) string {
	if input.Format != "" {
		return strings.ReplaceAll(input.Format, "%s", value)
	}
	return value
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

	profile := findProfile(tool, req.Profile)

	modifiers := findModifiers(tool, req.Modifiers)

	t.Command = append(t.Command, tool.Execute.Command)

	profileTokens := strings.Split(profile.Format, " ")
	t.Command = append(t.Command, profileTokens...)
	t.Command = append(t.Command, modifiers...)

	for _, iovar := range tool.Execute.Inputs {
		for k, v := range req.Inputs {
			if iovar.Name == k {
				if iovar.Type == tsf.FILE {
					injectVariable(t.Command, fmt.Sprint("{", k, "}"), fmt.Sprint(FILES_PREFIX, v))
				} else {
					injectVariable(t.Command, fmt.Sprint("{", k, "}"), formatVariable(iovar, v))
				}
				break
			}
		}
	}

	for _, outputFile := range tool.Execute.Outputs {
		injectVariable(
			t.Command,
			fmt.Sprint("{", outputFile.Name, "}"),
			fmt.Sprint(fmt.Sprint(FILES_PREFIX, OUTPUT_PREFIX), outputFile.Name),
		)
	}

	return t, nil
}

func StartTask(t *schemas.Task) (string, error) {
	c, err := podman.CreateContainer(t.Tool.Spec.Sandbox.Name, t.Command, t.Env)
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
