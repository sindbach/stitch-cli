package commands

import (
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
	u "github.com/sindbach/stitch-cli/user"

	"github.com/mitchellh/cli"
)

// NewAtlasOrgCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewAtlasOrgCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {

		return &AtlasOrgCommand{
			BaseCommand: &BaseCommand{
				Name: "org",
				UI:   ui,
			},
		}, nil
	}
}

// AtlasOrgCommand ...
type AtlasOrgCommand struct {
	*BaseCommand

	flagOrgList   bool
	flagProjectID string
}

// Help returns long-form help information for this command
func (ec *AtlasOrgCommand) Help() string {
	return `Atlas Organizations

OPTIONS:
  --list
	Get all Atlas organizations the authenticated user has access to.
  --project-id [string]
    Get an Atlas organization for a specific ID.
` +
		ec.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ec *AtlasOrgCommand) Synopsis() string {
	return `Access Atlas Organizations.`
}

// Run executes the command
func (ec *AtlasOrgCommand) Run(args []string) int {
	set := ec.NewFlagSet()

	set.BoolVar(&ec.flagOrgList, "list", false, "")
	set.StringVar(&ec.flagProjectID, "project-id", "", "")

	if err := ec.BaseCommand.run(args); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}

	if !ec.flagOrgList && ec.flagProjectID == "" {
		ec.UI.Error("see --help for more information")
		return 1
	}

	if err := ec.run(ec.flagOrgList, ec.flagProjectID); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}
	return 0
}

func (ec *AtlasOrgCommand) run(flagList bool, flagProjectID string) error {

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

	if flagProjectID != "" {
		ps, err := ac.ProjectByOrgID(flagProjectID)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		result := tm.NewTable(0, 5, 5, ' ', 0)
		fmt.Fprintf(result, "ID\tName\tReplicaSetCount\tShardCount\n")
		for _, p := range ps {
			fmt.Fprintf(result, "%s\t%s\t%d\t%d\n", p.ID, p.Name, p.ReplicaSetCount, p.ShardCount)
		}
		tm.Println(result)
		tm.Flush()
		return nil
	}

	u, err := ac.UserByName(user.Username)
	if err != nil {
		return fmt.Errorf("failed to list User info: %s", err)
	}

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ID\tName\tRole\n")

	for _, role := range u.Roles {
		if strings.HasPrefix(role.Name, "ORG_") {
			org, err := ac.OrgByID(role.OrgID)
			if err != nil {
				fmt.Println("Warning:", err)
				continue
			}
			fmt.Fprintf(result, "%s\t%s\t%s\n", org.ID, org.Name, role.Name)
		}
	}
	tm.Println(result)
	tm.Flush()
	return nil
}
