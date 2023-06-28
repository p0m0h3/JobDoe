package task

import (
	"errors"
	"strings"

	"fuzz.codes/fuzzercloud/workerengine/container"
	"fuzz.codes/fuzzercloud/workerengine/tool"
	"github.com/go-playground/validator/v10"
)

func ValidateRequest[Request any](r Request) ([]string, error) {

	var validate = validator.New()

	var badFields []string
	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			badFields = append(badFields, err.StructNamespace())
		}
	}
	return badFields, err
}

func InjectVariables(c []string, p string, v string) {
	for i, slice := range c {
		if slice == p {
			c[i] = v
			return
		}
	}
}

func NewTaskSpec(req CreateTaskRequest) (container.ContainerSpec, error) {
	tool := *tool.Tools[req.ToolName]
	spec := container.ContainerSpec{}
	spec.ImageName = tool.Name
	spec.Command = make([]string, 0)

	modifier, ok := tool.Exe.Modifiers[req.Modifier]
	if !ok {
		return spec, errors.New("could not find modifier")
	}

	spec.Command = append(spec.Command, tool.Exe.Command, modifier.String)

	for _, varPlaceholder := range modifier.Variables {
		variableName := strings.Trim(varPlaceholder, "{}")
		found := false
		for k, v := range req.InputList {
			if variableName == k {
				InjectVariables(spec.Command, varPlaceholder, v)
				found = true
				break
			}
		}
		if !found {
			return spec, errors.New("could not satisfy command variables")
		}
	}

	spec.Stdin = req.Stdin
	spec.EnvVarList = req.EnvVarList

	return spec, nil
}
