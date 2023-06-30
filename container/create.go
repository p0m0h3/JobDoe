package container

import (
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
)

func CreateContainer(
	image string,
	command []string,
	env map[string]string,
) (string, error) {
	err := PullImage(image)
	if err != nil {
		return "", err
	}
	s := specgen.NewSpecGenerator(image, false)
	s.Command = command
	s.Env = env

	createResponse, err := containers.CreateWithSpec(Connection, s, nil)
	if err != nil {
		return "", err
	}
	if err := containers.Start(Connection, createResponse.ID, nil); err != nil {
		return "", err
	}

	return createResponse.ID, nil
}
