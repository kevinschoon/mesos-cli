package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/pailer"
	"os"
)

type Read struct {
	AgentID  *string
	Path     *string
	Follow   *bool
	Lines    *int
	Hostname *string
	profile  Profile
}

func (_ Read) Name() string { return "read" }
func (_ Read) Desc() string { return "Read the contents of a file" }
func (r *Read) SetProfile(p Profile) {
	r.profile = func() *config.Profile {
		profile := p()
		if *r.Hostname != "" {
			profile = profile.With(
				config.Master(*r.Hostname),
			)
		}
		return profile
	}
}

func (r *Read) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID PATH"
		r.AgentID = cmd.StringArg("ID", "", "AgentID")
		r.Path = cmd.StringArg("PATH", "", "path to read")
		r.Follow = cmd.BoolOpt("f follow", false, "follow the content")
		r.Lines = cmd.IntOpt("n nlines", 0, "number of lines to read")
		r.Hostname = cmd.StringOpt("m master", "", "mesos master")
		cmd.Action = func() {
			caller, err := NewAgentCaller(r.profile(), *r.AgentID)
			failOnErr(err)
			pg := &pailer.FilePaginator{
				Path:   *r.Path,
				Follow: *r.Follow,
				Max:    uint64(*r.Lines),
			}
			failOnErr(pailer.Monitor(caller, os.Stdout, -1, pg))
		}
	}
}
