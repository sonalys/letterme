package utils

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type invalidStruct struct {
	C string `json:"c"`
}

func (c invalidStruct) Validate() error {
	return errors.New("mock")
}

func Test_LoadFromEnv(t *testing.T) {
	type testStruct struct {
		FieldA string `json:"a"`
		FieldB int64  `json:"b"`
	}

	testEnvKey := "TEST_ENV_KEY"

	expConfig := &testStruct{
		FieldA: "foo",
		FieldB: -54,
	}

	buf, err := json.Marshal(expConfig)
	require.NoError(t, err)

	err = os.Setenv(testEnvKey, string(buf))
	require.NoError(t, err)

	t.Run("all ok", func(t *testing.T) {
		resultConfig := new(testStruct)
		err = LoadFromEnv(testEnvKey, resultConfig)
		require.NoError(t, err, "should parse env to struct")

		require.Equal(t, expConfig, resultConfig)
	})

	t.Run("invalid struct", func(t *testing.T) {
		resultConfig := new(invalidStruct)
		err = LoadFromEnv(testEnvKey, resultConfig)
		require.Error(t, err, "should not parse wrong struct")
	})

	t.Run("without validatable interface should pass", func(t *testing.T) {
		resultConfig := new(struct{})
		err = LoadFromEnv(testEnvKey, resultConfig)
		require.NoError(t, err, "we can't validate configs without validatable interface")
	})

	t.Run("missing env key", func(t *testing.T) {
		resultConfig := new(struct{})
		err = LoadFromEnv("NON_EXISTENT_ENV_KEY", resultConfig)
		require.Error(t, err, "should give error for not found keys")
	})

	t.Run("invalid env value", func(t *testing.T) {
		err = os.Setenv(testEnvKey, "{ invalidJson: 456 }")
		require.NoError(t, err)

		resultConfig := new(struct{})
		err = LoadFromEnv(testEnvKey, resultConfig)
		require.Error(t, err, "should give error for invalid json")
	})
}
