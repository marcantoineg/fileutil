package fileutil

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFilePath = "./.test.json"

func Test_SaveToFile(t *testing.T) {
	testRuns := []struct {
		testName         string
		dataToSave       any
		expectedFileData string
		expectedError    error
	}{
		{
			"with nil",
			nil,
			"null",
			nil,
		},
		{
			"with struct",
			struct {
				Data string
			}{"hello, file!"},
			`{"Data": "hello, file!"}`,
			nil,
		},
		{
			"with complex struct",
			struct {
				Data                  string
				notExportedField      string
				NotTheActualFieldName string `json:"custom-named-field"`
			}{"hello, file!", "not exported", "custom-named-field-data"},
			`{"Data": "hello, file!","custom-named-field": "custom-named-field-data"}`,
			nil,
		},
		{
			"with string",
			"some-data",
			`"some-data"`,
			nil,
		},
	}

	for _, tr := range testRuns {
		t.Run(tr.testName, func(tt *testing.T) {
			os.Remove(testFilePath)
			tt.Cleanup(func() { os.Remove(testFilePath) })

			err := SaveToFile(tr.dataToSave, testFilePath)
			bytes, _ := os.ReadFile(testFilePath)
			actualFileData := strings.ReplaceAll(string(bytes), "\n", "")
			actualFileData = strings.ReplaceAll(actualFileData, "  ", "")

			assert.Equal(tt, tr.expectedError, err)
			assert.Equal(tt, tr.expectedFileData, actualFileData)
		})
	}
}

func Test_CreateEmptyFile(t *testing.T) {
	os.Remove(testFilePath)
	t.Cleanup(func() { os.Remove(testFilePath) })

	err := CreateEmptyFile(testFilePath)

	exists := Exists(testFilePath)

	assert.Nil(t, err)
	assert.True(t, exists)
}

func Test_CreateEmptyListFile(t *testing.T) {
	os.Remove(testFilePath)
	t.Cleanup(func() { os.Remove(testFilePath) })

	err := CreateEmptyListFile(testFilePath)

	bytes, err := os.ReadFile(testFilePath)

	assert.Nil(t, err)
	assert.Equal(t, "[]", string(bytes))
}

func Test_ReplaceTilde(t *testing.T) {
	testRuns := []struct {
		testName     string
		givenPath    string
		expectedPath string
	}{
		{
			"without tilde in path",
			"/this/is/a/path",
			"/this/is/a/path",
		},
		{
			"with only tilde in path",
			"~",
			"$HOME",
		},
		{
			"with tilde in path",
			"~/a/path",
			"$HOME/a/path",
		},
		{
			"with tilde not as first char",
			"/something/~",
			"/something/~",
		},
	}

	for _, tr := range testRuns {
		t.Run(tr.testName, func(tt *testing.T) {
			actual := ReplaceTilde(tr.givenPath)
			homeDir, _ := os.UserHomeDir()

			assert.Equal(tt, strings.ReplaceAll(tr.expectedPath, "$HOME", homeDir), actual)
		})
	}
}

func Test_Exists(t *testing.T) {
	os.Create(testFilePath)

	assert.True(t, Exists(testFilePath))

	os.Remove(testFilePath)
	assert.False(t, Exists(testFilePath))
}
