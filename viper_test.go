package main

import (
	"github.com/nate-johnston/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitViperWillReturnErrorOnReadError(t *testing.T) {
	err := InitViper("fixtures/viper/not_there.yml", nil)
	assert.NotNil(t, err)
	assert.Equal(t, "open fixtures/viper/not_there.yml: no such file or directory", err.Error())
}

func TestInitViperEnsureItCallsValidateConfigWithConfigErrors(t *testing.T) {
	// By calling InitViper and expecting an error case we prove that IniViper is
	// calling ValidateConfig() and its working as expected.
	viperConfigs := []ViperConfig{
		ViperConfig{Key: "missing_required_string", Description: "Required String"},
	}
	err := InitViper("fixtures/viper/test.yml", viperConfigs)
	assert.NotNil(t, err)
	assert.Equal(t, "Required Configuration Missing: [Key: missing_required_string, Description: Required String]",
		err.Error())
}

func TestValidateConfig(t *testing.T) {
	viperConfigs := []ViperConfig{
		ViperConfig{Key: "required_string", Description: "Required String"},
		ViperConfig{Key: "nested.one", Description: "Required Nested One"},
		ViperConfig{Key: "test_default", Default: "Optional Value", Description: "Optional Value"},
		ViperConfig{Key: "test_alias", Alias: []string{"bubba", "forest"}, Description: "Optional Value"},
	}
	InitViper("fixtures/viper/test.yml", viperConfigs)

	err := ValidateConfig(viperConfigs)
	assert.Nil(t, err)
	assert.Equal(t, "Optional Value", viper.GetString("test_default"))
	assert.Equal(t, "test_alias value", viper.GetString("test_alias"))
	assert.Equal(t, "test_alias value", viper.GetString("bubba"))
	assert.Equal(t, "test_alias value", viper.GetString("forest"))
}
