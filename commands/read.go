package commands

import (
	"github.com/jawher/mow.cli"
	//"github.com/vektorlab/mesos-cli/client"
	"github.com/vektorlab/mesos-cli/config"
)

// TODO
type Read struct {
	*command
}

func NewRead() Command {
	return Read{
		&command{"read", "Read the contents of a file"},
	}
}

func (read Read) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		/*
				defaults := DefaultProfile()
				cmd.Spec = "[OPTIONS] FILE"
				var (
					master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
						//lines    = cmd.IntOpt("n lines", 0, "Output the last N lines")
						//tail     = cmd.BoolOpt("t tail", false, "Tail output")
						filename = cmd.StringArg("FILE", "", "Filename to retrieve")

						all         = cmd.BoolOpt("all", false, "Show all tasks")
						frameworkID = cmd.StringOpt("framework", "", "Filter FrameworkID")
						fuzzy       = cmd.BoolOpt("fuzzy", true, "Fuzzy match Task name or Task ID prefix")
						name        = cmd.StringOpt("name", "", "Filter Task name")
						id          = cmd.StringOpt("id", "", "Filter Task ID")
						state       = cmd.StringsOpt("state", []string{"TASK_FINISHED"}, "Filter based on Task state")
				)

				cmd.Before = func() {
					//*state = trimFlaged(*state, "--state")
				}

				//fmt.Println(*state)

				cmd.Action = func() {
					client := httpcli.New(httpcli.Endpoint(config.Profile(WithMaster(*master)).Master))
					fmt.Println(client)
					/*
						client := &Client{
							Hostname: config.Profile(WithMaster(*master)).Master,
						}
						filters, err := NewTaskFilters(&TaskFilterOptions{
							All:         *all,
							FrameworkID: *frameworkID,
							Fuzzy:       *fuzzy,
							ID:          *id,
							Name:        *name,
							States:      *state,
						})
						failOnErr(err)
						task, err := client.Task(filters...)
						failOnErr(err)
						agent, err := client.Agent(AgentFilterId(task.GetTaskId().GetValue()))
						failOnErr(err)
						// Lookup executor information in agent state
						client = &Client{Hostname: FQDN(agent)}
						executor, err := client.Executor(ExecutorFilterId(task.GetExecutorId().GetValue()))
						failOnErr(err)
						fmt.Println(executor, filename)
						/*
							files, err := client.Browse(executor.Directory)
							for _, file := range files {
								if file.Relative() == *filename {
									target = file
								}
							}
							if target == nil {
								failOnErr(fmt.Errorf("cannot find file %s", *filename))
							}
							fp := &FilePaginator{
								data:   make(chan *fileData),
								cancel: make(chan bool),
								path:   target.Path,
								tail:   *tail,
							}
							failOnErr(Monitor(client, os.Stdout, *lines, fp))
				}
			}
		*/
	}
}
