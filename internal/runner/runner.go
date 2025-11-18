package runner

type Runner interface {
	Run(code string, env map[string]string) string
	GetOutput(jobuuid string) (string, error)
}

func NewDefaultRunner() (Runner, error) {
	return NewK8sRunner()
}
