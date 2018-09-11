package commands

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/mitchellh/cli"
	u "github.com/sindbach/stitch-cli/user"
)

// NewAtlasProjectCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewAtlasProjectCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {

		return &AtlasProjectCommand{
			BaseCommand: &BaseCommand{
				Name: "project",
				UI:   ui,
			},
		}, nil
	}
}

// AtlasProjectCommand ...
type AtlasProjectCommand struct {
	*BaseCommand

	flagProjectList bool
	flagProjectID   string
	flagOrgID       string
}

// Help returns long-form help information for this command
func (ac *AtlasProjectCommand) Help() string {
	return `Atlas Projects

OPTIONS:
  --list
	Get all Atlas projects the authenticated user has access to.
  --project-id [string]
	Get an Atlas project using a specific ID.
  --org-id [string]
	Get an Atlas project using a specific ID.
` +
		ac.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ac *AtlasProjectCommand) Synopsis() string {
	return `Access Atlas Projects.`
}

// Run executes the command
func (ac *AtlasProjectCommand) Run(args []string) int {
	set := ac.NewFlagSet()

	set.BoolVar(&ac.flagProjectList, "list", false, "")
	set.StringVar(&ac.flagProjectID, "project-id", "", "")
	set.StringVar(&ac.flagOrgID, "org-id", "", "")

	if err := ac.BaseCommand.run(args); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}
	if !ac.flagProjectList && ac.flagProjectID == "" && ac.flagOrgID == "" {
		ac.UI.Error("see --help for more information")
		return 1
	}
	if err := ac.run(ac.flagProjectList, ac.flagOrgID, ac.flagProjectID); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (ac *AtlasProjectCommand) run(flagList bool, flagOrgID string, flagProjectID string) error {

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

	if flagOrgID != "" {
		ps, err := client.ProjectByOrgID(flagOrgID)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		result := tm.NewTable(0, 5, 5, ' ', 0)
		fmt.Fprintf(result, "ProjectID\tName\tReplicaSetCount\tShardCount\n")
		for _, p := range ps {
			fmt.Fprintf(result, "%s\t%s\t%d\t%d\n", p.ID, p.Name, p.ReplicaSetCount, p.ShardCount)
		}
		tm.Println(result)
		tm.Flush()
		return nil
	}

	if flagProjectID != "" {
		p, err := client.ProjectByID(flagProjectID)
		if err != nil {
			return fmt.Errorf("failed to list Project info: %s", err)
		}

		result := tm.NewTable(0, 5, 5, ' ', 0)
		fmt.Fprintf(result, "ProjectID\tName\tOrgID\tReplicaSet\tShard\n")
		fmt.Fprintf(result, "%s\t%s\t%s\t%d\t%d\n", p.ID, p.Name, p.OrgID, p.ReplicaSetCount, p.ShardCount)
		tm.Println(result)
		tm.Flush()
		return nil
	}

	ps, err := client.Projects()
	if err != nil {
		return fmt.Errorf("failed to list Projects: %s", err)
	}

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ProjectID\tName\tOrgID\tReplicaSet\tShard\n")

	for _, p := range ps {
		fmt.Fprintf(result, "%s\t%s\t%s\t%d\t%d\n", p.ID, p.Name, p.OrgID, p.ReplicaSetCount, p.ShardCount)
	}
	tm.Println(result)
	tm.Flush()
	return nil
}
