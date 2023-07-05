package podman

import (
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
)

func CreateContainer(
	image string,
	command []string,
	env map[string]string,
) (*entities.ContainerCreateResponse, error) {
	err := PullImage(image)
	if err != nil {
		return nil, err
	}
	s := specgen.NewSpecGenerator(image, false)
	s.Command = command
	s.Env = env

	createResponse, err := containers.CreateWithSpec(Connection, s, nil)
	if err != nil {
		return nil, err
	}
	if err := containers.Start(Connection, createResponse.ID, nil); err != nil {
		return nil, err
	}

	return &createResponse, nil
}
