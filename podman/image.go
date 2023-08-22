package podman

import "github.com/containers/podman/v4/pkg/bindings/images"

var (
	PULLPOLICY = "newer"
	PULLQUIET  = true
)

func PullImage(name string) error {
	_, err := images.Pull(Connection, name, &images.PullOptions{
		Policy: &PULLPOLICY,
		Quiet:  &PULLQUIET,
	})

	return err
}
