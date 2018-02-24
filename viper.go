package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// ViperConfig holds information per config item.  If an Default
// is not set then it is assumed that the value is Required.
// NOTE: You can not set a default on a nested value.  i.e. a value
//       within a has in a json or yaml file. (nested.value) you can
//       set nested values as required.
type ViperConfig struct {
	Key         string      // The config key that is required.
	Default     interface{} // Default Value to set.
	Alias       []string    // Any key Aliases that should be registered
	Description string      // Description of the config.
}

// InitViper with the passed path and config.
func InitViper(path string, viperConfigs []ViperConfig) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := ValidateConfig(viperConfigs); err != nil {
		return err
	}
	return nil
}

// ValidateConfig will check the defined var ViperConfigs []ViperConfig and validate
// the existances of the required keys, and set defaults for all keys where defaults are
// defined.
func ValidateConfig(viperConfigs []ViperConfig) error {
	var errs []error
	for _, rc := range viperConfigs {
		if rc.Default == nil && viper.Get(rc.Key) == nil {
			errs = append(errs, fmt.Errorf("Key: %s, Description: %s", rc.Key, rc.Description))
		} else {
			viper.SetDefault(rc.Key, rc.Default)
		}

		if len(rc.Alias) > 0 {
			for _, a := range rc.Alias {
				viper.RegisterAlias(a, rc.Key)
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("Required Configuration Missing: %v", errs)
	}
	return nil
}
