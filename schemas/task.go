package schemas

type CreateTaskRequest struct {
	ToolID   string            `json:"toolId" validate:"required,alphanum"`
	Modifier string            `json:"modifier" validate:"required,alphanum"`
	Inputs   map[string]string `json:"inputs"`
	Env      map[string]string `json:"env"`
	Stdin    string            `json:"stdin"`
}

type Task struct {
	ID      string            `json:"id"`
	Command []string          `json:"cmd"`
	Env     map[string]string `json:"env"`
	Stdin   string            `json:"stdin"`
	Status  string            `json:"status"`
	Tool    Tool              `json:"tool"`
}
