package podman

import (
	"os"

	"github.com/containers/podman/v4/pkg/bindings/images"
)

func PullImage(name string) error {
	var (
		PULLPOLICY   = "newer"
		PULLQUIET    = true
		PULLAUTHFILE = os.Getenv("REGISTRY_AUTH_FILE")
	)

	_, err := images.Pull(Connection, name, &images.PullOptions{
		Policy:   &PULLPOLICY,
		Quiet:    &PULLQUIET,
		Authfile: &PULLAUTHFILE,
	})

	return err
}
