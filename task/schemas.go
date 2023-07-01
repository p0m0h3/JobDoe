package task

type CreateTaskRequest struct {
	ToolName   string            `json:"name" validate:"required,alphanum"`
	Modifier   string            `json:"modifier" validate:"required,alphanum"`
	InputList  map[string]string `json:"inputs"`
	EnvVarList map[string]string `json:"env"`
	Stdin      string            `json:"stdin"`
}

type GetTaskResponse struct {
	ID        string
	ImageName string
	Status    string
	Command   []string
}

type CreateTaskResponse struct {
	ID   string `json:"id"`
	Tool string `json:"tool"`
}

type ErrorResponse struct {
	Code       int      `json:"code"`
	Validation []string `json:"validation,omitempty"`
	Message    string   `json:"message"`
}
