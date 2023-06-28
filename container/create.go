package container

import (
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
)

type ContainerSpec struct {
	ID         string
	ImageName  string
	Command    []string
	EnvVarList map[string]string
	Stdin      string
}

func CreateContainer(t ContainerSpec) (string, error) {
	err := PullImage(t.ImageName)
	if err != nil {
		return "", err
	}
	s := specgen.NewSpecGenerator(t.ImageName, false)
	s.Command = t.Command
	s.Env = t.EnvVarList

	createResponse, err := containers.CreateWithSpec(Connection, s, nil)
	if err != nil {
		return "", err
	}
	if err := containers.Start(Connection, createResponse.ID, nil); err != nil {
		return "", err
	}

	return createResponse.ID, nil
}
