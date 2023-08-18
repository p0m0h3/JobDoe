package state

import (
	"archive/tar"
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
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

func makeArchiveFromFiles(files map[string]string) io.Reader {
	buffer := &bytes.Buffer{}
	ar := tar.NewWriter(buffer)
	defer ar.Close()

	for name, content := range files {
		// decode the base64 content and write it onto a buffer to get the size
		decoder := b64.NewDecoder(b64.StdEncoding, strings.NewReader(content))
		data := &bytes.Buffer{}
		io.Copy(data, decoder)

		// write file header to archive buffer
		header := &tar.Header{
			Name:  name,
			Mode:  int64(os.FileMode(0660)),
			Uname: "root",
			Gname: "root",
			Size:  int64(data.Len()),
		}
		ar.WriteHeader(header)
		io.Copy(ar, data)
	}

	return buffer
}

func copyTaskFilesIntoContainer(t *schemas.Task) error {
	ar := makeArchiveFromFiles(t.Files)
	return podman.CopyIntoContainer(t.ID, ar, FILES_PREFIX)
}

func findProfile(t *tsf.Tool, profName string) (tsf.Profile, bool) {
	for _, prof := range t.Execute.Profiles {
		if prof.Name == profName {
			return prof, true
		}
	}

	return tsf.Profile{}, false
}

func findModifiers(t *tsf.Tool, modNames []string) ([]string, bool) {
	result := make([]string, 0)
	found := false

	for _, modName := range modNames {
		for _, mod := range t.Execute.Modifiers {
			if modName == mod.Name {
				result = append(result, mod.Format)
				found = true
			}
		}
	}

	return result, found
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
		Files: make(map[string]string),
	}

	t.Command = append(t.Command, tool.Execute.Command)

	if profile, found := findProfile(tool, req.Profile); found {
		profileTokens := strings.Split(profile.Format, " ")
		t.Command = append(t.Command, profileTokens...)
	}

	if modifiers, found := findModifiers(tool, req.Modifiers); found {
		t.Command = append(t.Command, modifiers...)
	}

	for _, iovar := range tool.Execute.Inputs {
		for k, v := range req.Inputs {
			if iovar.Name == k {
				if iovar.Type == tsf.FILE {
					injectVariable(t.Command, fmt.Sprint("{", k, "}"), fmt.Sprint(FILES_PREFIX, k))
					t.Files[k] = v

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

	err = copyTaskFilesIntoContainer(t)
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
