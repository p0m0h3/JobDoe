package podman

import (
	"math/rand"

	"git.fuzz.codes/fuzzercloud/workerengine/config"
	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func appendTaskProxyEnv(env map[string]string) map[string]string {
	c, err := config.GetConfig()
	if err != nil {
		return env
	}

	if env == nil {
		env = make(map[string]string)
	}

	proxy := "http://" + c.Proxies[rand.Intn(len(c.Proxies))]
	env["http_proxy"] = proxy
	env["HTTP_PROXY"] = proxy
	env["https_proxy"] = proxy
	env["HTTPS_PROXY"] = proxy

	return env
}

func CreateContainer(
	image string,
	command []string,
	env map[string]string,
	memory int64,
	CPU uint64,
) (*entities.ContainerCreateResponse, error) {
	s := specgen.NewSpecGenerator(image, false)
	s.Command = command
	s.Env = appendTaskProxyEnv(env)

	s.ResourceLimits = &specs.LinuxResources{}
	swap := 2 * memory
	s.ResourceLimits.Memory = &specs.LinuxMemory{
		Limit: &memory,
		Swap:  &swap,
	}
	s.ResourceLimits.CPU = &specs.LinuxCPU{
		Shares: &CPU,
	}

	createResponse, err := containers.CreateWithSpec(Connection, s, nil)
	if err != nil {
		return nil, err
	}

	return &createResponse, nil
}

func DeleteContainer(id string) error {
	if err := StopContainer(id); err != nil {
		return err
	}
	TRUE := true

	options := &containers.RemoveOptions{
		Force:   &TRUE,
		Volumes: &TRUE,
	}

	if _, err := containers.Remove(Connection, id, options); err != nil {
		return err
	}
	return nil
}

func PruneTasks() error {
	_, err := containers.Prune(Connection, nil)
	return err
}

func StartContainer(id string) error {
	if err := containers.Start(Connection, id, nil); err != nil {
		return err
	}
	return nil
}

func StopContainer(id string) error {
	if err := containers.Stop(Connection, id, nil); err != nil {
		return err
	}

	return nil
}

func InspectContainer(id string) (*define.InspectContainerData, error) {
	data, err := containers.Inspect(Connection, id, new(containers.InspectOptions).WithSize(true))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAllContainers() ([]entities.ListContainer, error) {
	all := true
	return containers.List(Connection, &containers.ListOptions{
		All: &all,
	})
}

func GetContainerStats(id string) (entities.ContainerStatsReport, error) {
	c, err := containers.Stats(Connection, []string{id}, nil)
	if err != nil {
		return entities.ContainerStatsReport{}, err
	}

	return <-c, nil
}

func WaitOnContainer(id string) error {
	options := &containers.WaitOptions{
		Condition: []define.ContainerStatus{define.ContainerStateExited},
	}
	_, err := containers.Wait(Connection, id, options)
	return err
}
