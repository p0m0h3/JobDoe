package state

import (
	"archive/tar"
	"bytes"
	b64 "encoding/base64"
	"errors"
	"io"
	"os"
	"strings"

	"git.fuzz.codes/fuzzercloud/workerengine/podman"
	"git.fuzz.codes/fuzzercloud/workerengine/schemas"
)

const (
	FILES_PREFIX = "/files/"
)

var Tasks map[string]*schemas.Task

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
	tool, found := Tools[req.Tool]
	if !found {
		return nil, errors.New("could not find tool")
	}

	if req.MemoryLimit == 0 {
		req.MemoryLimit = 209715200
	}
	if req.CPULimit == 0 {
		req.CPULimit = 1
	}

	t := &schemas.Task{
		Command: make([]string, 0),
		Env:     req.Env,
		Tool:    tool,
		Files:   make(map[string]string),
	}

	t.Command = append(t.Command, tool.Execute.Command)

	if req.Command != nil {
		t.Command = append(t.Command, req.Command...)
		return t, nil
	}

	// handle files
	for modifierName, inputs := range req.Inputs {
		modifier, err := tool.Execute.FindModifier(modifierName)
		if err != nil {
			return t, err
		}

		for _, variable := range modifier.Variables {
			if variable.Type == "file" {
				content, found := inputs[variable.Name]
				if !found {
					continue
				}
				t.Files[variable.Name] = content
				inputs[variable.Name] = FILES_PREFIX + variable.Name
			}
		}
	}

	// handle profiles and modifiers
	var format []string
	var err error

	if req.Profile != "" {
		format, err = tool.Execute.ProfileFormat(req.Profile, req.Inputs)
		if err != nil {
			return t, err
		}

	} else if len(req.Modifiers) != 0 {
		format, err = tool.Execute.Format(req.Modifiers, req.Inputs)
		if err != nil {
			return t, err
		}
	}

	t.Command = append(t.Command, format...)

	return t, nil
}

func StartTask(t *schemas.Task) (string, error) {
	c, err := podman.CreateContainer(
		t.Tool.Header.Image,
		t.Command,
		t.Env,
		t.MemoryLimit,
		t.CPULimit,
	)
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
