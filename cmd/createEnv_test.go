package cmd

import (
	"os"
	"testing"
)

func TestBuildEnvFileString(t *testing.T) {
	secretMap := make(map[string]interface{})
	secretMap["testKey"] = "testVal"
	secretMap["testKey2"] = "testVal2"

	result := BuildEnvFileString(secretMap)

	want := "testKey=testVal\ntestKey2=testVal2\n"
	if result != want {
		t.Fatalf(`BuildEnvFileString() result: %q, want: %q`, result, want)
	}
}

func TestInitEnvFile(t *testing.T) {
	envFileNameResult := ""

	osGetWd = func() (string, error) {
		return "", nil
	}

	osCreate = func(fileName string) (*os.File, error) {
		envFileNameResult = fileName

		return nil, nil
	}

	t.Run("Creates file correctly where path ends with /", func(t *testing.T) {
		path := "./"
		envFileName := "testFileName1"

		_, err := InitEnvFile(path, envFileName)

		if err != nil {
			t.Errorf("InitEnvFile(%q,%q) Error: %v", path, envFileName, err)
		}

		expectedFileName := "/./testFileName1"

		if envFileNameResult != expectedFileName {
			t.Errorf("InitEnvFile() incorrect file name result: %q, want %q", envFileNameResult, expectedFileName)
		}
	})

	t.Run("Creates file correclty where path ends without /", func(t *testing.T) {
		path := ".."
		envFileName := "testFileName2"

		_, err := InitEnvFile(path, envFileName)

		if err != nil {
			t.Errorf("InitEnvFile(%q,%q) Error: %v", path, envFileName, err)
		}

		expectedFileName := "/../testFileName2"

		if envFileNameResult != expectedFileName {
			t.Errorf("InitEnvFile() incorrect file name result: %q, want %q", envFileNameResult, expectedFileName)
		}
	})

}
