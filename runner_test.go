package goverter

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/jmattheis/goverter/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestScenario(t *testing.T) {
	curPath := getCurrentPath()
	scenarios := filepath.Join(curPath, "scenario")
	execDir := filepath.Join(curPath, "execution")
	files, err := ioutil.ReadDir(scenarios)
	require.NoError(t, err)

	require.NoError(t, os.MkdirAll(execDir, os.ModePerm))
	require.NoError(t, clearDir(execDir))

	for _, file := range files {
		require.False(t, file.IsDir(), "should not be a directory")

		t.Run(file.Name(), func(t *testing.T) {
			scenarioFileName := filepath.Join(scenarios, file.Name())
			scenarioBytes, err := ioutil.ReadFile(scenarioFileName)
			require.NoError(t, err)

			scenario := Scenario{}
			err = yaml.Unmarshal(scenarioBytes, &scenario)
			require.NoError(t, err)

			for name, content := range scenario.Input {
				inPath := filepath.Join(execDir, name)
				err = os.MkdirAll(filepath.Dir(inPath), os.ModePerm)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(execDir, name), []byte(content), os.ModePerm)
				require.NoError(t, err)
			}
			genPkgName := "generated"

			global := append([]string{"outputPackage github.com/jmattheis/goverter/execution/" + genPkgName}, scenario.Global...)

			patterns := scenario.Patterns
			if len(patterns) == 0 {
				patterns = append(patterns, "github.com/jmattheis/goverter/execution")
			}

			files, err := generateConvertersRaw(
				&GenerateConfig{
					WorkingDir:      execDir,
					PackagePatterns: patterns,
					Global: config.RawLines{
						Lines:    global,
						Location: "scenario global",
					},
				})

			actualOutputFiles := toOutputFiles(execDir, files)

			if os.Getenv("UPDATE_SCENARIO") == "true" && scenario.ErrorStartsWith == "" {
				if err != nil {
					scenario.Success = []*OutputFile{}
					scenario.Error = replaceAbsolutePath(curPath, fmt.Sprint(err))
				} else {
					scenario.Success = toOutputFiles(execDir, files)
					scenario.Error = ""
				}
				newBytes, err := yaml.Marshal(&scenario)
				if assert.NoError(t, err) {
					os.WriteFile(scenarioFileName, newBytes, os.ModePerm)
				}
			}

			if scenario.ErrorStartsWith != "" {
				require.Error(t, err)
				strErr := replaceAbsolutePath(curPath, fmt.Sprint(err))
				require.Equal(t, scenario.ErrorStartsWith, strErr[:len(scenario.ErrorStartsWith)])
				return
			}

			if scenario.Error != "" {
				require.Error(t, err)
				require.Equal(t, scenario.Error, replaceAbsolutePath(curPath, fmt.Sprint(err)))
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, scenario.Success, "scenario.Success may not be empty")
			require.Equal(t, scenario.Success, actualOutputFiles)

			err = writeFiles(files)
			require.NoError(t, err)
			require.NoError(t, compile(execDir), "generated converter doesn't build")
		})
		if os.Getenv("SKIP_CLEAN") != "true" {
			clearDir(execDir)
		}
	}
}

func replaceAbsolutePath(curPath, body string) string {
	return strings.ReplaceAll(body, curPath, "/ABSOLUTE")
}

func compile(dir string) error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = dir
	_, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("Process exited with %d:\n%s", exit.ExitCode(), string(exit.Stderr))
		}
	}
	return err
}

func toOutputFiles(execDir string, files map[string][]byte) []*OutputFile {
	output := []*OutputFile{}
	for fileName, content := range files {
		rel, err := filepath.Rel(execDir, fileName)
		if err != nil {
			panic("could not create relpath")
		}
		output = append(output, &OutputFile{Name: rel, Content: string(content)})
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Name < output[j].Name
	})
	return output
}

type Scenario struct {
	Input  map[string]string `yaml:"input"`
	Global []string          `yaml:"global,omitempty"`

	Patterns []string      `yaml:"patterns,omitempty"`
	Success  []*OutputFile `yaml:"success,omitempty"`
	// for error cases, use either Error or ErrorStartsWith, not both
	Error           string `yaml:"error,omitempty"`
	ErrorStartsWith string `yaml:"error_starts_with,omitempty"`
}

type OutputFile struct {
	Name    string
	Content string
}

func (f *OutputFile) MarshalYAML() (interface{}, error) {
	return map[string]string{f.Name: f.Content}, nil
}

func (f *OutputFile) UnmarshalYAML(value *yaml.Node) error {
	v := map[string]string{}
	err := value.Decode(&v)

	for name, content := range v {
		f.Name = name
		f.Content = content
	}

	return err
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
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
