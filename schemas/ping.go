package schemas

type PingResponse struct {
	Version string `json:"version"`
	Spec    string `json:"spec"`
	Mode    string `json:"mode"`
	Tools   int    `json:"tools"`
	Tasks   int    `json:"tasks"`
}
