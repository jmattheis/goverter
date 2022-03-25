package goverter

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestScenario(t *testing.T) {
	scenarios := path.Join(getCurrentPath(), "scenario")
	execDir := path.Join(getCurrentPath(), "execution")
	files, err := ioutil.ReadDir(scenarios)
	require.NoError(t, err)

	require.NoError(t, os.MkdirAll(execDir, os.ModePerm))
	require.NoError(t, clearDir(execDir))

	for _, file := range files {
		require.False(t, file.IsDir(), "should not be a directory")

		t.Run(file.Name(), func(t *testing.T) {
			scenarioFileName := path.Join(scenarios, file.Name())
			scenarioBytes, err := ioutil.ReadFile(scenarioFileName)
			require.NoError(t, err)

			scenario := Scenario{}
			err = yaml.Unmarshal(scenarioBytes, &scenario)
			require.NoError(t, err)

			for name, content := range scenario.Input {
				err = ioutil.WriteFile(path.Join(execDir, name), []byte(content), os.ModePerm)
				require.NoError(t, err)
			}
			genFile := path.Join(execDir, "generated", "generated.go")
			err = GenerateConverterFile(
				genFile,
				GenerateConfig{
					PackageName:   "generated",
					ScanDir:       "github.com/jmattheis/goverter/execution",
					ExtendMethods: scenario.Extends,
				})

			body, _ := ioutil.ReadFile(genFile)

			if os.Getenv("UPDATE_SCENARIO") == "true" && scenario.ErrorStartsWith == "" {
				if err != nil {
					scenario.Success = ""
					scenario.Error = replaceAbsolutePath(fmt.Sprint(err))
				} else {
					scenario.Success = string(body)
					scenario.Error = ""
				}
				newBytes, err := yaml.Marshal(&scenario)
				if assert.NoError(t, err) {
					ioutil.WriteFile(scenarioFileName, newBytes, os.ModePerm)
				}
			}

			if scenario.ErrorStartsWith != "" {
				require.Error(t, err)
				strErr := replaceAbsolutePath(fmt.Sprint(err))
				require.Equal(t, scenario.ErrorStartsWith, strErr[:len(scenario.ErrorStartsWith)])
				return
			}

			if scenario.Error != "" {
				require.Error(t, err)
				require.Equal(t, scenario.Error, replaceAbsolutePath(fmt.Sprint(err)))
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, scenario.Success, "scenario.Success may not be empty")
			require.Equal(t, scenario.Success, string(body))
			require.NoError(t, compile(genFile), "generated converter doesn't build")
		})
		clearDir(execDir)
	}
}

func replaceAbsolutePath(body string) string {
	return strings.ReplaceAll(body, getCurrentPath(), "/ABSOLUTE")
}

func compile(file string) error {
	cmd := exec.Command("go", "build", "")
	cmd.Dir = filepath.Dir(file)
	_, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("Process exited with %d:\n%s", exit.ExitCode(), string(exit.Stderr))
		}
	}
	return err
}

type Scenario struct {
	Input   map[string]string `yaml:"input"`
	Extends []string          `yaml:"extends,omitempty"`
	Success string            `yaml:"success,omitempty"`
	// for error cases, use either Error or ErrorStartsWith, not both
	Error           string `yaml:"error,omitempty"`
	ErrorStartsWith string `yaml:"error_starts_with,omitempty"`
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
