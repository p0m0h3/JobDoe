package state

import (
	"fuzz.codes/fuzzercloud/tsf"
	"fuzz.codes/fuzzercloud/workerengine/schemas"
)

func init() {
	Tasks = make(map[string]*schemas.Task)
	Tools = make(map[string]*tsf.Tool)
}
