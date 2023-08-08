package schemas

import "time"

type CreateTaskRequest struct {
	ToolID   string            `json:"tool" validate:"required,ascii"`
	Modifier string            `json:"modifier" validate:"required,ascii"`
	Inputs   map[string]string `json:"inputs"`
	Env      map[string]string `json:"env"`
	Stdin    string            `json:"stdin"`
	Files    string            `json:"files"`
}

type Task struct {
	ID      string            `json:"id"`
	Command []string          `json:"cmd"`
	Env     map[string]string `json:"env"`
	Stdin   string            `json:"-"`
	Status  string            `json:"status"`
	Tool    Tool              `json:"-"`
	Files   string            `json:"-"`
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
