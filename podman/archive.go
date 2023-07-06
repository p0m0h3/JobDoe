package podman

import (
	"io"

	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func CopyFileIntoContainer(id string, data io.Reader, path string) error {
	copy, err := containers.CopyFromArchive(Connection, id, path, data)
	if err != nil {
		return err
	}

	return copy()
}
