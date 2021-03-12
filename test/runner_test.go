package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jmattheis/go-genconv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestScenario(t *testing.T) {
	scenarios := path.Join(getCurrentPath(), "scenario")
	execDir := path.Join(getCurrentPath(), "execution")
	files, err := ioutil.ReadDir(scenarios)
	require.NoError(t, err)

	require.NoError(t, os.MkdirAll(execDir, 0755))
	require.NoError(t, clearDir(execDir))

	for _, file := range files {
		clearDir(execDir)

		require.False(t, file.IsDir(), "should not be a directory")
		t.Run(file.Name(), func(t *testing.T) {
			scenarioFileName := path.Join(scenarios, file.Name())
			scenarioBytes, err := ioutil.ReadFile(scenarioFileName)

			scenario := Scenario{}
			err = yaml.Unmarshal(scenarioBytes, &scenario)
			require.NoError(t, err)

			for name, content := range scenario.Input {
				err = ioutil.WriteFile(path.Join(execDir, name), []byte(content), 0644)
				require.NoError(t, err)
			}

			body, err := genconv.Generate(genconv.GenerateConfig{
				PackageName: "generated",
				ScanDir:     "github.com/jmattheis/go-genconv/test/execution",
			})

			if os.Getenv("UPDATE_SCENARIO") == "true" {
				fmt.Println("UPDATING ", file.Name())
				if err != nil {
					scenario.Success = ""
					scenario.Error = err.Error()
				} else {
					scenario.Success = string(body)
					scenario.Error = ""
				}
				newBytes, err := yaml.Marshal(&scenario)
				if assert.NoError(t, err) {
					ioutil.WriteFile(scenarioFileName, newBytes, 0644)
				}
			}

			if scenario.Error != "" {
				require.Error(t, err)
				require.EqualError(t, err, scenario.Error)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, scenario.Success, "scenario.Success may not be empty")
				require.Equal(t, scenario.Success, string(body))
			}
		})
		clearDir(execDir)
	}
}

type Scenario struct {
	Input   map[string]string `yaml:"input"`
	Error   string            `yaml:"error,omitempty"`
	Success string            `yaml:"success,omitempty"`
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}

func clearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}
