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

func GetContainerStats(id string) (entities.ContainerStatsReport, error) {
	c, err := containers.Stats(Connection, []string{id}, nil)
	if err != nil {
		return entities.ContainerStatsReport{}, err
	}

	return <-c, nil
}

func WaitOnContainer(id string) error {
	options := &containers.WaitOptions{
		Condition: []define.ContainerStatus{define.ContainerStateExited},
	}
	_, err := containers.Wait(Connection, id, options)
	return err
}
