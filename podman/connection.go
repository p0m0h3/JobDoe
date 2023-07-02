package podman

import (
	"context"

	"github.com/containers/podman/v4/pkg/bindings"
)

var Connection context.Context

func OpenConnection(socket string) error {
	var err error
	Connection, err = bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return err
	}
	return nil
}
