package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
)

const (
	OperatorAPIPath  = "/api/v1"
	SchedulerAPIPath = "/api/v1/scheduler"
)

func defaults() *Profile {
	return &Profile{
		Master: "http://localhost:5050",
	}
}

// Profile contains environment specific options
type Profile struct {
	Master string `json:"master"`
	Scheme string `json:"scheme"`
}

// Options are functional profile options
type Option func(*Profile)

func Master(m string) Option {
	return func(p *Profile) {
		p.Master = m
	}
}

func (p *Profile) With(opts ...Option) *Profile {
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p Profile) Endpoint() *url.URL {
	u, _ := url.Parse(p.Master)
	u.Path = OperatorAPIPath
	return u
}

type Config struct {
	Profiles map[string]*Profile `json:profiles`
}

// LoadProfile loads a user configuration
// from ~/.mesos-cli.json creating a
// JSON file with defaults if it does
// not exist.
func LoadProfile(path, name string) (*Profile, error) {
	config := &Config{Profiles: map[string]*Profile{}}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.Profiles["default"] = defaults()
		raw, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, err
		}
		return config.Profiles["default"], ioutil.WriteFile(path, raw, os.FileMode(0755))
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(raw, &config); err != nil {
		return nil, err
	}
	if _, ok := config.Profiles[name]; !ok {
		return nil, fmt.Errorf("Cannot load profile: %s", name)
	}
	return config.Profiles[name], nil
}

func HomeDir() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.HomeDir
}
