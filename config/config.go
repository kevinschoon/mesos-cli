package config

import (
	"encoding/json"
	"errors"
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

var ErrNoConfig = errors.New("No Config File Present")

func DefaultProfile() *Profile {
	return &Profile{
		Master: "localhost:5050",
		Scheme: "http",
	}
}

// Profile contains environment specific options
type Profile struct {
	Master string `json:"master"`
	Scheme string `json:"scheme"`
}

func (p Profile) Endpoint() url.URL {
	return url.URL{
		Scheme: p.Scheme,
		Host:   p.Master,
	}
}

// ProfileOptions are functional profile options
type ProfileOption func(*Profile)

func WithMaster(m string) ProfileOption {
	return func(p *Profile) {
		// Checks to be sure we do not override with a default value from CLI
		if m != DefaultProfile().Master {
			p.Master = m
		}
	}
}

func WithSchema(m string) ProfileOption {
	return func(p *Profile) {
		if m != DefaultProfile().Scheme {
			p.Scheme = m
		}
	}
}

// Merge the options from another profile
func (p *Profile) Merge(other *Profile) {
	if other.Master != "" {
		p.Master = other.Master
	}
	if other.Scheme != "" {
		p.Scheme = other.Scheme
	}
}

func (p *Profile) With(opts ...ProfileOption) {
	for _, opt := range opts {
		opt(p)
	}
}

type CfgFn func() *Config

// Config is a global configuration file usually stored
// in the user's home (~/.mesos-cli.json).
type Config struct {
	profile  string
	Profiles map[string]*Profile `json:"profiles"`
}

// Profile returns a profile loaded from disks with optional
// commandline overrides.
func (c Config) Profile(opts ...ProfileOption) *Profile {
	// Start with all default options
	profile := DefaultProfile()
	if other, ok := c.Profiles[c.profile]; ok {
		// Merge (override) any options included
		// in the user profile
		profile.Merge(other)
	} else {
		panic(fmt.Sprintf("unknown profile %s", c.profile))
	}
	// Finally any command line flags
	// take precedence
	for _, opt := range opts {
		opt(profile)
	}
	return profile
}

func loadConfig(path string, config *Config) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrNoConfig
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, config)
}

// LoadConfig loads a user configuration
// from ~/.mesos-cli.json creating a
// JSON file with defaults if it does
// not exist.
func LoadConfig(path, profile string) (*Config, error) {
	// Default config
	config := &Config{
		profile: profile,
		Profiles: map[string]*Profile{
			"default": DefaultProfile(),
		},
	}
	err := loadConfig(path, config)
	if err != nil && err != ErrNoConfig {
		return nil, err
	}
	// If there is no configuration file
	// save the default one above.
	if err == ErrNoConfig {
		return config, SaveConfig(path, config)
	}
	return config, nil
}

func SaveConfig(path string, config *Config) error {
	raw, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, raw, os.FileMode(0755))
}

func HomeDir() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.HomeDir
}
