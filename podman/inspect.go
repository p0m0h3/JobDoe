package podman

import (
	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

func InspectContainer(id string) (*define.InspectContainerData, error) {
	data, err := containers.Inspect(Connection, id, new(containers.InspectOptions).WithSize(true))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAllContainers() ([]entities.ListContainer, error) {
	all := true
	return containers.List(Connection, &containers.ListOptions{
		All: &all,
	})
}
