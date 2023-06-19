package tool

import (
	"os"
	"strings"
)

var Tools map[string]string = make(map[string]string)

func ReadTools() error {
	entries, err := os.ReadDir(os.Getenv("TOOLS_DIRECTORY"))
	if err != nil {
		return err
	}

	for _, toolFile := range entries {
		if strings.HasSuffix(toolFile.Name(), ".toml") {
			Tools[strings.TrimSuffix(toolFile.Name(), ".toml")] = toolFile.Name()
		}
	}

	return nil
}
