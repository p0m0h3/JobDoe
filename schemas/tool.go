package schemas

import "fuzz.codes/fuzzercloud/tsf"

type GetToolResponse struct {
	ID   string   `json:"id"`
	Spec tsf.Tool `json:"spec"`
}
