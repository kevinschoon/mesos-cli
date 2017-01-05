package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var testConfig = []byte(`
{
	"profiles": {
		"my_profile": {
			"master": "my-host:5050"
		}
	}
}
`)

// Configuration should be loaded
// in the following order:
// 1) Default options
// 2) Profile options (~/.mesos-cli.json)
// 3) Commandline flags
// With the highest number taking
// precedence
func TestLoadConfig(t *testing.T) {
	ioutil.WriteFile("/tmp/mesos-cli_test_config.json", testConfig, os.FileMode(0755))

	config, _ := LoadConfig("/tmp/mesos-cli_test_config.json", "default")
	assert.Equal(t, config.profile, "default")
	assert.Equal(t, "127.0.0.1:5050", config.Profile().Master)
	assert.Equal(t, "override:5050", config.Profile(WithMaster("override:5050")).Master)

	config, _ = LoadConfig("/tmp/mesos-cli_test_config.json", "my_profile")
	assert.Equal(t, config.profile, "my_profile")
	assert.Equal(t, "my-host:5050", config.Profile().Master)

	config, _ = LoadConfig("/tmp/mesos-cli_test_config.json", "my_profile")
	assert.Equal(t, "override:5050", config.Profile(WithMaster("override:5050")).Master)

}
