package state

import (
	"git.fuzz.codes/fuzzercloud/tsf"
	"git.fuzz.codes/fuzzercloud/workerengine/schemas"
)

func init() {
	Tasks = make(map[string]*schemas.Task)
	Tools = make(map[string]*tsf.Spec)
}
