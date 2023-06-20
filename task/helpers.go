package task

import (
	"context"
	"errors"
	"os"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/go-playground/validator/v10"
)

var Connection context.Context

func ValidateCreateTaskRequest(r CreateTaskRequest) ([]string, error) {

	var validate = validator.New()

	var badFields []string
	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			badFields = append(badFields, err.StructNamespace())
		}
	}
	return badFields, err
}

type TaskSpec struct {
	ID        string
	ImageName string
	Command   []string
	Stdin     string
}

func InjectVariable(c []string, p string, v string) {
	for i, slice := range c {
		if slice == p {
			c[i] = v
			return
		}
	}
}

func NewTaskSpec(tool tsf.Tool, req CreateTaskRequest) (TaskSpec, error) {
	spec := TaskSpec{}
	spec.ImageName = tool.Name
	spec.Command = make([]string, 0)

	modifier, ok := tool.Exe.Modifiers[req.Modifier]
	if !ok {
		return spec, errors.New("could not find modifier")
	}

	spec.Command = append(spec.Command, tool.Exe.Command, modifier.String)

	for _, varPlaceholder := range modifier.Variables {
		variableName := strings.Trim(varPlaceholder, "{}")
		found := false
		for k, v := range req.InputList {
			if variableName == k {
				InjectVariable(spec.Command, varPlaceholder, v)
				found = true
				break
			}
		}
		if !found {
			return spec, errors.New("could not satisfy command variables")
		}
	}

	spec.Stdin = req.Stdin

	return spec, nil
}

func NewContainerTask(t TaskSpec) (string, error) {

	_, err := images.Pull(Connection, t.ImageName, nil)
	if err != nil {
		return "", err
	}
	s := specgen.NewSpecGenerator(t.ImageName, false)
	s.Command = t.Command

	createResponse, err := containers.CreateWithSpec(Connection, s, nil)
	if err != nil {
		return "", err
	}
	if err := containers.Start(Connection, createResponse.ID, nil); err != nil {
		return "", err
	}

	return createResponse.ID, nil
}

func InitConnection() {
	var err error
	Connection, err = bindings.NewConnection(context.Background(), os.Getenv("PODMAN_SOCKET_ADDRESS"))
	if err != nil {
		panic(err)
	}
}
