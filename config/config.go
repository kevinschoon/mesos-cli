package config

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
	"github.com/satori/go.uuid"
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
		TaskInfo: &mesos.TaskInfo{
			TaskID: mesos.TaskID{
				Value: uuid.NewV4().String(),
			},
			Command: &mesos.CommandInfo{
				Environment: &mesos.Environment{
					Variables: []mesos.Environment_Variable{},
				},
			},
			Container: &mesos.ContainerInfo{
				Volumes: []mesos.Volume{},
			},
			Resources: []mesos.Resource{
				mesos.Resource{
					Name:   "cpus",
					Type:   mesos.SCALAR.Enum(),
					Role:   proto.String("*"),
					Scalar: &mesos.Value_Scalar{Value: 0.1},
				},
			},
			Labels: &mesos.Labels{},
		},
	}
}

// Profile contains environment specific options
type Profile struct {
	Master   string `json:"master"`
	TaskInfo *mesos.TaskInfo
}

// Options are functional profile options
type Option func(*Profile)

func Master(master *string) Option {
	return func(p *Profile) {
		if ptrStrSet(master) {
			p.Master = *master
		}
	}
}

func Command(cmd *string) Option {
	return func(p *Profile) {
		if ptrStrSet(cmd) {
			p.TaskInfo.Command.Value = cmd
		}
	}
}

func TaskID(id *string) Option {
	return func(p *Profile) {
		if ptrStrSet(id) {
			p.TaskInfo.TaskID.Value = *id
		}
	}
}

func User(user *string) Option {
	return func(p *Profile) {
		if ptrStrSet(user) {
			p.TaskInfo.Command.User = user
		}
	}
}

func ptrStrSet(s *string) bool {
	if s == nil {
		return false
	}
	if *s == "" {
		return false
	}
	return true
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
