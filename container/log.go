package container

import (
	"errors"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func GetContainerLog(id string) (chan string, error) {
	con, err := InspectContainer(id)
	if err != nil {
		return nil, err
	}

	status, err := define.StringToContainerStatus(con.State.Status)
	if err != nil {
		return nil, err
	}

	if status == define.ContainerStateStopped {
		return nil, errors.New("container is not stopped")
	}

	out := make(chan string, 1024)

	err = containers.Logs(Connection, id, nil, out, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
