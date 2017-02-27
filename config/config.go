package config

import (
	"encoding/json"
	"fmt"
	"github.com/mesos/mesos-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
)

const (
	OperatorAPIPath  = "/api/v1"
	SchedulerAPIPath = "/api/v1/scheduler"
)

// Options are functional profile options
type Option func(*Profile)

// Profile contains environment specific options
type Profile struct {
	Master  string `json:"master"`
	Debug   bool   `json:"debug"`
	Restart bool   `json:"restart"`
	log     *zap.Logger
}

func (p *Profile) Log() *zap.Logger {
	if p.log == nil {
		cfg := zapConfig(p.Debug)
		if p.Debug {
			cfg.Level.SetLevel(zap.DebugLevel)
			cfg.EncoderConfig.CallerKey = "caller"
		}
		logger, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		p.log = logger
	}
	return p.log
}

func (p *Profile) With(opts ...Option) *Profile {
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p Profile) Framework() *mesos.FrameworkInfo {
	return &mesos.FrameworkInfo{
		ID:   &mesos.FrameworkID{Value: ""},
		User: "root",
		Name: "mesos-cli",
	}
}

func (p Profile) Scheduler() *url.URL {
	u, _ := url.Parse(p.Master)
	u.Path = SchedulerAPIPath
	return u
}

func (p Profile) Endpoint() *url.URL {
	u, _ := url.Parse(p.Master)
	u.Path = OperatorAPIPath
	return u
}

type Config struct {
	Profiles map[string]*Profile `json:"profiles"`
}

// LoadProfile loads a user configuration
// from ~/.mesos-cli.json creating a
// JSON file with defaults if it does
// not exist.
func LoadProfile(path, name string) (profile *Profile, err error) {
	config := &Config{Profiles: map[string]*Profile{}}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		config.Profiles["default"] = defaults()
		raw, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(path, raw, os.FileMode(0755))
		if err != nil {
			return nil, err
		}
		profile = config.Profiles["default"]
	} else {
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
		profile = config.Profiles[name]
	}
	return profile, nil
}

func zapConfig(debug bool) *zap.Config {
	return &zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

}

func defaults() *Profile {

	return &Profile{
		Master: "http://localhost:5050",
	}
}

// Master specifies the Mesos master hostname
func Master(master string) Option {
	return func(p *Profile) {
		if master != "" {
			p.Master = master
		}
	}
}

func Debug(debug bool) Option {
	return func(p *Profile) {
		p.Debug = debug
	}
}

func Restart(restart bool) Option {
	return func(p *Profile) {
		p.Restart = restart
	}
}

func HomeDir() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.HomeDir
}
