package schemas

type CreateTaskRequest struct {
	ToolID   string            `json:"toolId" validate:"required,alphanum"`
	Modifier string            `json:"modifier" validate:"required,alphanum"`
	Inputs   map[string]string `json:"inputs"`
	Env      map[string]string `json:"env"`
	Stdin    string            `json:"stdin"`
}
