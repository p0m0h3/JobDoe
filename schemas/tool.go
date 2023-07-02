package schemas

import "fuzz.codes/fuzzercloud/tsf"

type Tool struct {
	ID   string    `json:"id"`
	Spec *tsf.Tool `json:"spec"`
}
