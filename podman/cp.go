package podman

import (
	"io"

	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func CopyIntoContainer(id string, data io.Reader, path string) error {
	copy, err := containers.CopyFromArchive(Connection, id, path, data)
	if err != nil {
		return err
	}

	return copy()
}

func CopyFromContainer(id string, output io.Writer, path string) error {
	copy, err := containers.CopyToArchive(Connection, id, path, output)
	if err != nil {
		return err
	}

	return copy()
}
