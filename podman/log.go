package podman

import (
	"errors"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func GetContainerLog(id string, output chan string) error {
	con, err := InspectContainer(id)
	if err != nil {
		return err
	}

	status, err := define.StringToContainerStatus(con.State.Status)
	if err != nil {
		return err
	}

	if status == define.ContainerStateStopped {
		return errors.New("container is not stopped")
	}

	follow := false
	err = containers.Logs(Connection, id, &containers.LogOptions{
		Follow: &follow,
	}, output, nil)

	if err != nil {
		return err
	}

	return nil
}
