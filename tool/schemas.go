package tool

import "fuzz.codes/fuzzercloud/tsf"

type GetToolResponse struct {
	Name string   `json:"name"`
	Spec tsf.Tool `json:"spec"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
