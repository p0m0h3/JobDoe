package schemas

type CreateTaskRequest struct {
	ToolID   string            `json:"toolId" validate:"required"`
	Modifier string            `json:"modifier" validate:"required"`
	Inputs   map[string]string `json:"inputs"`
	Env      map[string]string `json:"env"`
	Stdin    string            `json:"stdin"`
	Files    map[string]string `json:"files"`
}

type Task struct {
	ID      string            `json:"id"`
	Command []string          `json:"cmd"`
	Env     map[string]string `json:"env"`
	Stdin   string            `json:"-"`
	Status  string            `json:"status"`
	Tool    Tool              `json:"-"`
	Files   map[string]string `json:"files"`
}
