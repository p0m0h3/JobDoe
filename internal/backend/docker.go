package backend

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	docker "github.com/moby/moby/client"
)

type DockerRunner struct {
	Client *docker.Client
}

func NewDockerRunner() (*DockerRunner, error) {
	c, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		return nil, err
	}
	return &DockerRunner{
		Client: c,
	}, nil
}

func (d *DockerRunner) Run(code string, env map[string]string) (string, error) {
	ctx := context.Background()
	containerName := uuid.NewString()
	envVars := make([]string, 0, len(env))

	for k, v := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	resp, err := d.Client.ContainerCreate(
		ctx,
		docker.ContainerCreateOptions{
			Name: containerName,
			Config: &container.Config{
				Image: "python:3.12-alpine",
				Cmd:   []string{"python", "-c", code},
				Env:   envVars,
			},
		},
	)
	if err != nil {
		return "", err
	}

	if _, err := d.Client.ContainerStart(ctx, resp.ID, docker.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (d *DockerRunner) GetOutput(id string) (string, error) {
	ctx := context.Background()

	inspect, err := d.Client.ContainerInspect(ctx, id, docker.ContainerInspectOptions{})
	if err != nil {
		return "", err
	}

	if !inspect.Container.State.Running && inspect.Container.State.ExitCode != 0 {
		return "", fmt.Errorf("container not succeeded")
	}

	// 2. Request logs
	logs, err := d.Client.ContainerLogs(ctx, id, docker.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		return "", err
	}
	defer logs.Close()

	// 3. Read logs into buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
