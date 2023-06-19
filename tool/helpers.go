package tool

import (
	"os"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
)

var Tools map[string]*tsf.Tool = make(map[string]*tsf.Tool)

func ReadTools() error {
	entries, err := os.ReadDir(os.Getenv("TOOLS_DIRECTORY"))
	if err != nil {
		return err
	}

	for _, toolFile := range entries {
		if strings.HasSuffix(toolFile.Name(), ".toml") {
			data, err := os.ReadFile(os.Getenv("TOOLS_DIRECTORY") + "/" + toolFile.Name())
			if err != nil {
				panic(err.Error() + " (Error opening file " + toolFile.Name() + ")")
			}

			tool, err := tsf.Parse(data)
			if err != nil {
				panic(err.Error() + " (Error parsing tsf file " + toolFile.Name() + ")")
			}
			Tools[strings.TrimSuffix(toolFile.Name(), ".toml")] = tool
		}
	}

	return nil
}
