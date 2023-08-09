package podman

import (
	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func GetContainerLog(id string, stderr bool, output chan string) error {
	False := false
	True := true
	err := containers.Logs(Connection, id, &containers.LogOptions{
		Follow: &False,
		Stderr: &stderr,
		Stdout: &True,
	}, output, output)
	if err != nil {
		return err
	}

	return nil
}

func StreamContainerLog(id string, stderr bool, output chan string) error {
	True := true
	return containers.Logs(Connection, id, &containers.LogOptions{
		Follow: &True,
		Stderr: &stderr,
		Stdout: &True,
	}, output, output)
}
