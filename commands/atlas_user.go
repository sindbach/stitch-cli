package commands

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/mitchellh/cli"
	u "github.com/sindbach/stitch-cli/user"
)

const (
	flagUserList = "list"
	flagUserID   = "user-id"
)

// NewAtlasUserCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewAtlasUserCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {

		return &AtlasUserCommand{
			BaseCommand: &BaseCommand{
				Name: "user",
				UI:   ui,
			},
		}, nil
	}
}

// AtlasUserCommand ...
type AtlasUserCommand struct {
	*BaseCommand

	flagUserList  bool
	flagProjectID string
}

// Help returns long-form help information for this command
func (ac *AtlasUserCommand) Help() string {
	return `Atlas Users

OPTIONS:
  --list
	Get all Atlas users the authenticated user has access to.
  --user-id [string]
	The Atlas User ID.
` +
		ac.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ac *AtlasUserCommand) Synopsis() string {
	return `Access Atlas Users.`
}

// Run executes the command
func (ac *AtlasUserCommand) Run(args []string) int {
	set := ac.NewFlagSet()

	set.BoolVar(&ac.flagUserList, flagUserList, false, "")
	set.StringVar(&ac.flagProjectID, flagProjectID, "", "")

	if err := ac.BaseCommand.run(args); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}
	if !ac.flagUserList && ac.flagProjectID == "" {
		ac.UI.Error("see --help for more information")
		return 1
	}
	if err := ac.run(ac.flagUserList, ac.flagProjectID); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (ac *AtlasUserCommand) printOutput(id string, name string, orgid string, replset int, shard int) {

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ID\tName\tOrgID\tReplicaSet\tShard\n")
	fmt.Fprintf(result, "%s\t%s\t%s\t%d\t%d\n", id, name, orgid, replset, shard)
	tm.Println(result)
	tm.Flush()
}

func (ac *AtlasProjectCommand) printOutput(id string, name string, orgid string, replset int, shard int) {

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ID\tName\tOrgID\tReplicaSet\tShard\n")
	fmt.Fprintf(result, "%s\t%s\t%s\t%d\t%d\n", id, name, orgid, replset, shard)
	tm.Println(result)
	tm.Flush()
}

func (ac *AtlasUserCommand) run(flagUserList bool, flagProjectID string) error {

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

	u, err := client.UserByName(user.Username)
	if err != nil {
		return fmt.Errorf("failed to list User info: %s", err)
	}

	result := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(result, "ID\tUsername\tEmail\n")
	fmt.Fprintf(result, "%s\t%s\t%s\n", u.ID, u.Username, u.Email)
	tm.Println(result)

	roleTable := tm.NewTable(0, 5, 5, ' ', 0)
	fmt.Fprintf(roleTable, "Name\tOrgID\tProjectID\n")
	for _, role := range u.Roles {
		fmt.Fprintf(roleTable, "%s\t%s\t%s\n", role.Name, role.OrgID, role.ProjectID)
	}
	tm.Println(roleTable)
	tm.Flush()
	return nil
}
