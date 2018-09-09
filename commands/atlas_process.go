package commands

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/mitchellh/cli"
	u "github.com/sindbach/stitch-cli/user"
)

// NewAtlasProcessCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewAtlasProcessCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {

		return &AtlasProcessCommand{
			BaseCommand: &BaseCommand{
				Name: "process",
				UI:   ui,
			},
		}, nil
	}
}

// AtlasProcessCommand ...
type AtlasProcessCommand struct {
	*BaseCommand

	flagProcessList bool
	flagProjectID   string
	flagProcessID   string
	flagLogType     string
}

// Help returns long-form help information for this command
func (ac *AtlasProcessCommand) Help() string {
	return `Atlas Processes
REQUIRED: 
  --project-id [string]
	Get an Atlas project using a specific ID.

OPTIONS:
  --list
	Get all Atlas process for a project.
  --process-id [string]
    Specify a Process ID
  --log [string]
	Get log file given a Process ID. Choices: mongodb, mongos, mongosqld, mongodb-audit-log, mongos-audit-log

` +
		ac.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ac *AtlasProcessCommand) Synopsis() string {
	return `Access Atlas Process.`
}

// Run executes the command
func (ac *AtlasProcessCommand) Run(args []string) int {
	set := ac.NewFlagSet()

	set.BoolVar(&ac.flagProcessList, "list", false, "")
	set.StringVar(&ac.flagProjectID, "project-id", "", "")
	set.StringVar(&ac.flagProcessID, "process-id", "", "")
	set.StringVar(&ac.flagLogType, "log", "", "")

	if err := ac.BaseCommand.run(args); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}
	if ac.flagProjectID == "" {
		ac.UI.Error("Project ID is required. See --help for more information")
		return 1
	}
	if err := ac.run(ac.flagProjectID, ac.flagProcessID, ac.flagLogType); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (ac *AtlasProcessCommand) run(flagProjectID string, flagProcessID string, flagLogType string) error {

	user, err := ac.User()
	if err != nil {
		return err
	}

	if !user.LoggedIn() {
		return u.ErrNotLoggedIn
	}

	client, err := ac.AtlasClient()
	if err != nil {
		return err
	}

	if flagProcessID != "" && flagLogType != "" {
		err = client.LogByProcessID(flagProjectID, flagProcessID, fmt.Sprintf("%s.gz", flagLogType))
		if err != nil {
			return fmt.Errorf("failed to fetch Log: %s", err)
		}
		return nil
	}

	ps, err := client.ProcessByProjectID(flagProjectID)
	if err != nil {
		return fmt.Errorf("failed to list Processes: %s", err)
	}

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ID\tReplicaSet\tVersion\tType\tLastping\n")

	for _, p := range ps {
		fmt.Fprintf(result, "%s\t%s\t%s\t%s\t%s\n", p.ID, p.ReplicasetName, p.Version, p.Lastping, p.Created)
	}
	tm.Println(result)

	tm.Flush()

	return nil

}
