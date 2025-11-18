package server

import (
	"github.com/p0m0h3/jobdoe/internal/config"
	"github.com/p0m0h3/jobdoe/internal/runner"
	"github.com/p0m0h3/jobdoe/internal/volume"
)

type Env struct {
	Runner runner.RunnerService
	Volume volume.VolumeService
	Config config.Config
}
