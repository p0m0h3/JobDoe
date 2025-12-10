package server

import (
	"github.com/p0m0h3/jobdoe/internal/backend"
	"github.com/p0m0h3/jobdoe/internal/config"
	"github.com/p0m0h3/jobdoe/internal/volume"
)

type Env struct {
	Backend backend.BackendService
	Volume  volume.VolumeService
	Config  config.Config
}
