package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"io/ioutil"
	"os"
	"os/user"
)

var (
	ConfigPath  = ""
	ErrNoConfig = errors.New("No Config File Present")
)

type Profile struct {
	Master        string               `json:"master"`
	FrameworkInfo *mesos.FrameworkInfo `json:"framework_info"`
}

type Config struct {
	profile  string
	Profiles map[string]*Profile `json:"profiles"`
}

func (c Config) Profile() *Profile {
	if profile, ok := c.Profiles[c.profile]; ok {
		return profile
	}
	panic(fmt.Errorf("No Profile %s", c.profile))
}

func loadConfig(config *Config) error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		return ErrNoConfig
	}
	raw, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, config)
}

// LoadConfig loads a user configuration
// from ~/.mesos-cli.json creating a
// JSON file with defaults if it does
// not exist.
func LoadConfig(profile, master string) (*Config, error) {
	// Default config
	config := &Config{
		profile: profile,
		Profiles: map[string]*Profile{
			"default": &Profile{
				Master: master,
				FrameworkInfo: &mesos.FrameworkInfo{
					Name: proto.String("mesos-cli"),
				},
			},
		},
	}
	err := loadConfig(config)
	if err != nil && err != ErrNoConfig {
		return nil, err
	}
	// If there is no configuration file
	// save the default one above.
	if err == ErrNoConfig {
		return config, SaveConfig(config)
	}
	return config, nil
}

func SaveConfig(config *Config) error {
	raw, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigPath, raw, os.FileMode(0755))
}

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	ConfigPath = fmt.Sprintf("%s/.mesos-cli.json", u.HomeDir)
}
