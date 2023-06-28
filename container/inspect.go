package container

import (
	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func InspectContainer(id string) (*define.InspectContainerData, error) {
	data, err := containers.Inspect(Connection, id, new(containers.InspectOptions).WithSize(true))
	if err != nil {
		return nil, err
	}

	return data, nil
}
