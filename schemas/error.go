package schemas

type ErrorResponse struct {
	Code    int      `json:"code"`
	Details []string `json:"details,omitempty"`
	Message string   `json:"message"`
}
