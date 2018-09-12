package commands

import (
	"fmt"

	tm "github.com/buger/goterm"
	u "github.com/sindbach/stitch-cli/user"

	"github.com/mitchellh/cli"
)

// NewAtlasClusterCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewAtlasClusterCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {

		return &AtlasClusterCommand{
			BaseCommand: &BaseCommand{
				Name: "cluster",
				UI:   ui,
			},
		}, nil
	}
}

// AtlasClusterCommand ...
type AtlasClusterCommand struct {
	*BaseCommand

	flagClusterList bool
	flagProjectID   string
}

// Help returns long-form help information for this command
func (ec *AtlasClusterCommand) Help() string {
	return `Atlas Clusters

REQUIRED: 
  --project-id [string] 
    Get an Atlas project using a specific ID.


OPTIONS:
  --list
	Get all Atlas organizations the authenticated user has access to.

` +
		ec.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ec *AtlasClusterCommand) Synopsis() string {
	return `Access Atlas Cluster.`
}

// Run executes the command
func (ec *AtlasClusterCommand) Run(args []string) int {
	set := ec.NewFlagSet()

	set.BoolVar(&ec.flagClusterList, "list", false, "")
	set.StringVar(&ec.flagProjectID, "project-id", "", "")

	if err := ec.BaseCommand.run(args); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}
	if ec.flagProjectID == "" {
		ec.UI.Error("Project ID is required. See --help for more information")
		return 1
	}
	if !ec.flagClusterList && ec.flagProjectID == "" {
		ec.UI.Error("see --help for more information")
		return 1
	}

	if err := ec.run(ec.flagProjectID); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}
	return 0
}

func (ec *AtlasClusterCommand) run(flagProjectID string) error {

	user, err := ec.User()
	if err != nil {
		return err
	}

	if !user.LoggedIn() {
		return u.ErrNotLoggedIn
	}

	ac, err := ec.AtlasClient()
	if err != nil {
		return err
	}

	cs, err := ac.ClustersByProjectID(flagProjectID)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ID\tName\tState\tPaused\n")

	for _, c := range cs {
		fmt.Fprintf(result, "%s\t%s\t%s\t%t\n", c.ID, c.Name, c.StateName, c.Paused)
	}
	tm.Println(result)
	tm.Flush()
	return nil
}
