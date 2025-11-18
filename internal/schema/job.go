package schema

type CreateJobRequest struct {
	Code string            `json:"code"`
	Env  map[string]string `json:"env"`
}

type CreateJobResponse struct {
	Code string            `json:"code"`
	Env  map[string]string `json:"env"`
}
