package podman

import (
	"github.com/containers/podman/v4/pkg/bindings/images"
)

func PullImage(name string, regauth string) error {
	var (
		PULLPOLICY   = "newer"
		PULLQUIET    = true
		PULLAUTHFILE = regauth
	)

	_, err := images.Pull(Connection, name, &images.PullOptions{
		Policy:   &PULLPOLICY,
		Quiet:    &PULLQUIET,
		Authfile: &PULLAUTHFILE,
	})

	return err
}
