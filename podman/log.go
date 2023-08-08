package podman

import (
	"errors"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func GetContainerLog(id string, stderr bool, output chan string) error {
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
		Stderr: &stderr,
	}, output, output)

	if err != nil {
		return err
	}

	return nil
}

func StreamContainerLog(id string, stderr bool, output chan string) error {
	True := true
	err := containers.Logs(Connection, id, &containers.LogOptions{
		Follow: &True,
		Stderr: &stderr,
	}, output, output)

	if err != nil {
		return err
	}

	return nil
}
