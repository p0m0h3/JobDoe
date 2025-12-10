package backend

import "fmt"

type BackendService interface {
	Run(code string, env map[string]string) (string, error)
	GetOutput(jobuuid string) (string, error)
}

func NewBackendService(t string) (BackendService, error) {
	switch t {
	case "k8s":
		return NewK8sRunner()
	case "docker":
		return NewDockerRunner()
	default:
		return nil, fmt.Errorf("unsupported backend: %s", t)
	}
}
