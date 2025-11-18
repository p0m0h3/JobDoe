package runner

type RunnerService interface {
	Run(code string, env map[string]string) string
	GetOutput(jobuuid string) (string, error)
}

func New() (RunnerService, error) {
	return NewK8sRunner()
}
