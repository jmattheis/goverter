package goverter

import (
	"fmt"
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
	rootDir := getCurrentPath()
	scenarioDir := filepath.Join(rootDir, "scenario")
	workDir := filepath.Join(rootDir, "execution")
	scenarioFiles, err := os.ReadDir(scenarioDir)
	require.NoError(t, err)
	require.NoError(t, clearDir(workDir))

	for _, file := range scenarioFiles {
		require.False(t, file.IsDir(), "should not be a directory")
		file := file

		testName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			testWorkDir := filepath.Join(workDir, testName)
			require.NoError(t, os.MkdirAll(testWorkDir, os.ModePerm))
			require.NoError(t, clearDir(testWorkDir))
			scenarioFilePath := filepath.Join(scenarioDir, file.Name())
			scenarioFileBytes, err := os.ReadFile(scenarioFilePath)
			require.NoError(t, err)

			scenario := Scenario{}
			err = yaml.Unmarshal(scenarioFileBytes, &scenario)
			require.NoError(t, err)
			err = os.WriteFile(filepath.Join(testWorkDir, "go.mod"), []byte("module github.com/jmattheis/goverter/execution\ngo 1.16"), os.ModePerm)
			require.NoError(t, err)

			for name, content := range scenario.Input {
				inPath := filepath.Join(testWorkDir, name)
				err = os.MkdirAll(filepath.Dir(inPath), os.ModePerm)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(testWorkDir, name), []byte(content), os.ModePerm)
				require.NoError(t, err)
			}
			genPkgName := "generated"

			global := append([]string{"output:package github.com/jmattheis/goverter/execution/" + genPkgName}, scenario.Global...)

			patterns := scenario.Patterns
			if len(patterns) == 0 {
				patterns = append(patterns, "github.com/jmattheis/goverter/execution")
			}

			files, err := generateConvertersRaw(
				&GenerateConfig{
					WorkingDir:      testWorkDir,
					PackagePatterns: patterns,
					Global: config.RawLines{
						Lines:    global,
						Location: "scenario global",
					},
				})

			actualOutputFiles := toOutputFiles(testWorkDir, files)

			if scenario.ErrorStartsWith != "" {
				require.Error(t, err)
				strErr := replaceAbsolutePath(testWorkDir, fmt.Sprint(err))
				require.Equal(t, scenario.ErrorStartsWith, strErr[:len(scenario.ErrorStartsWith)])
				return
			}

			if os.Getenv("UPDATE_SCENARIO") == "true" {
				if err != nil {
					scenario.Success = []*OutputFile{}
					scenario.Error = replaceAbsolutePath(testWorkDir, fmt.Sprint(err))
				} else {
					scenario.Success = toOutputFiles(testWorkDir, files)
					scenario.Error = ""
				}
				newBytes, err := yaml.Marshal(&scenario)
				if assert.NoError(t, err) {
					os.WriteFile(scenarioFilePath, newBytes, os.ModePerm)
				}
			}

			if scenario.Error != "" {
				require.Error(t, err)
				require.Equal(t, scenario.Error, replaceAbsolutePath(testWorkDir, fmt.Sprint(err)))
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, scenario.Success, "scenario.Success may not be empty")
			require.Equal(t, scenario.Success, actualOutputFiles)

			err = writeFiles(files)
			require.NoError(t, err)
			require.NoError(t, compile(testWorkDir), "generated converter doesn't build")
		})
	}
}

func replaceAbsolutePath(curPath, body string) string {
	return strings.ReplaceAll(body, curPath, "@workdir")
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
