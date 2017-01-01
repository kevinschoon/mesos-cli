package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

var ErrNoConfig = errors.New("No Config File Present")

type Profile struct {
	Master string
}

type Config struct {
	profile  string
	Profiles map[string]*Profile
}

func (c Config) Profile() *Profile {
	if profile, ok := c.Profiles[c.profile]; ok {
		return profile
	}
	panic(fmt.Errorf("No Profile %s", c.profile))
}

// Loads a user configuration from
// ~/.mesos-exec.json creating
// an empty JSON file if it does
// not exist.
func loadConfig(config *Config) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/.mesos-exec.json", u.HomeDir)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return ErrNoConfig
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(raw, config); err != nil {
		return err
	}
	return nil
}

func LoadConfig(profile, master string) (*Config, error) {
	config := &Config{
		profile: profile,
		Profiles: map[string]*Profile{
			"default": &Profile{
				Master: master,
			},
		},
	}
	err := loadConfig(config)
	if err != nil && err != ErrNoConfig {
		return nil, err
	}
	return config, nil
}
