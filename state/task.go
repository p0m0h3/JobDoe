package state

import (
	"archive/tar"
	"bytes"
	b64 "encoding/base64"
	"errors"
	"io"
	"os"
	"strings"

	"git.fuzz.codes/fuzzercloud/tsf"
	"git.fuzz.codes/fuzzercloud/workerengine/config"
	"git.fuzz.codes/fuzzercloud/workerengine/podman"
	"git.fuzz.codes/fuzzercloud/workerengine/schemas"
)

const (
	FILES_PREFIX   = "/files/"
	OUTPUTS_PREFIX = "/outputs/"
	MIN_MEMORY     = 209715200
	MIN_CPU        = 1
	CPU_FREQ       = 5 // 500 MHz
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
	// TODO: serious refactor needed for formatting the command and output files
	tool, found := Tools[req.Tool]
	if !found {
		return nil, errors.New("could not find tool")
	}

	if req.Inputs == nil {
		req.Inputs = make(map[string]map[string]string)
	}

	if req.Memory == 0 {
		req.Memory = MIN_MEMORY
	}
	if req.CPU == 0 {
		req.CPU = MIN_CPU
	}

	t := &schemas.Task{
		Command: make([]string, 0),
		Env:     req.Env,
		Tool:    tool,
		Files:   make(map[string]string),
		Memory:  req.Memory,
		CPU:     req.CPU,
	}

	t.Command = append(t.Command, strings.Split(tool.Execute.Command, " ")...)

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
	var modifiers []*tsf.Modifier
	var err error

	modifiers, err = tool.Execute.FindModifiers(req.Modifiers)
	if err != nil {
		return t, err
	}

	if req.Profile != "" {
		// find profile
		profile, err := tool.Execute.FindProfile(req.Profile)
		if err != nil {
			return t, err
		}
		req.Inputs = tool.Execute.SetProfileDefaults(profile, req.Inputs)

		profileModifiers, err := tool.Execute.FindModifiers(profile.Modifiers)
		if err != nil {
			return t, err
		}

		modifiers = append(modifiers, profileModifiers...)
	}

	// handle output files
	for _, modifier := range modifiers {
		for _, variable := range modifier.Variables {
			if variable.Type == "output" {
				req.Inputs[modifier.Name] = make(map[string]string)
				req.Inputs[modifier.Name][variable.Name] = OUTPUTS_PREFIX + variable.Name
			}
		}
	}

	format, err = tool.Execute.Format(modifiers, req.Inputs)
	if err != nil {
		return t, err
	}

	t.Command = append(t.Command, format...)

	return t, nil
}

func StartTask(t *schemas.Task) (string, error) {
	config, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	err = podman.PullImage(t.Tool.Header.Image, config.RegistryAuth)
	if err != nil {
		return "", err
	}

	c, err := podman.CreateContainer(
		t.Tool.Header.Image,
		t.Command,
		t.Env,
		t.Memory,
		t.CPU*CPU_FREQ,
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

	// make the empty output directory
	outputArchive := makeArchiveFromFiles(make(map[string]string))
	err = podman.CopyIntoContainer(t.ID, outputArchive, OUTPUTS_PREFIX)
	if err != nil {
		return t.ID, err
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
