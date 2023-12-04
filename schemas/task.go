package schemas

import (
	"time"

	"git.fuzz.codes/fuzzercloud/tsf"
)

type CreateTaskRequest struct {
	Tool      string                       `json:"tool" validate:"required,printascii"`
	Modifiers []string                     `json:"modifiers" validate:"printascii"`
	Profile   string                       `json:"profile" validate:"printascii"`
	Command   []string                     `json:"command" validate:"printascii"`
	Inputs    map[string]map[string]string `json:"inputs"`
	Env       map[string]string            `json:"env"`
	Memory    int64                        `json:"memory"`
	CPU       uint64                       `json:"cpu"`
}

type Task struct {
	ID      string            `json:"id"`
	Command []string          `json:"command"`
	Env     map[string]string `json:"env"`
	Status  string            `json:"status"`
	Files   map[string]string `json:"-"`
	Tool    *tsf.Spec         `json:"-"`
	Memory  int64             `json:"memory"`
	CPU     uint64            `json:"cpu"`
}

type GetTaskStatsResponse struct {
	ID          string        `json:"id"`
	AvgCPU      float64       `json:"avgcpu"`
	CPU         float64       `json:"cpu"`
	MemUsage    uint64        `json:"memory"`
	MemLimit    uint64        `json:"memoryLimit"`
	MemPerc     float64       `json:"memoryPercent"`
	NetInput    uint64        `json:"netInput"`
	NetOutput   uint64        `json:"netOutput"`
	BlockInput  uint64        `json:"blockInput"`
	BlockOutput uint64        `json:"blockOutput"`
	UpTime      time.Duration `json:"uptime"`
	Duration    uint64        `json:"duration"`
}
