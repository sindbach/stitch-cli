package commands

import (
	"fmt"

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

	flagList bool
}

// Help returns long-form help information for this command
func (ec *AtlasOrgCommand) Help() string {
	return `Atlas Organizations

OPTIONS:
  --list
	Get all Atlas organizations the authenticated user has access to.
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

	set.BoolVar(&ec.flagList, "list", false, "")

	if err := ec.BaseCommand.run(args); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}

	if err := ec.run(); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (ec *AtlasOrgCommand) run() error {

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

	orgs, err := ac.Orgs()
	if err != nil {
		return fmt.Errorf("failed to list Organizations: %s", err)
	}

	for _, org := range orgs {
		fmt.Println(org.Name)
		fmt.Println(org.ID)

	}
	return nil
}
