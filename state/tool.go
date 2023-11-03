package state

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.fuzz.codes/fuzzercloud/tsf"
)

var Tools map[string]*tsf.Spec

func ReadTools() error {
	entries, err := os.ReadDir(os.Getenv("TOOLS_DIRECTORY"))
	if err != nil {
		return err
	}

	for _, toolFile := range entries {
		if strings.HasSuffix(toolFile.Name(), ".toml") {
			data, err := os.ReadFile(os.Getenv("TOOLS_DIRECTORY") + "/" + toolFile.Name())
			if err != nil {
				log.Println(fmt.Errorf("%v (Error reading file %s)", err, toolFile.Name()))
			}

			tool, err := tsf.Parse(data)
			if err != nil {
				log.Println(fmt.Errorf("%v (Error parsing TSF file %s)", err, toolFile.Name()))
			}
			if err == nil {
				Tools[strings.TrimSuffix(toolFile.Name(), ".toml")] = &tool
			}
		}
	}

	return nil
}
