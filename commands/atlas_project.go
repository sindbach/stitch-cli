package commands

import (
	"fmt"

	u "github.com/sindbach/stitch-cli/user"

	"github.com/mitchellh/cli"
)

const (
	flagList      = "list"
	flagProjectID = "project-id"
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

	flagList      bool
	flagProjectID string
}

// Help returns long-form help information for this command
func (ac *AtlasProjectCommand) Help() string {
	return `Atlas Projects

OPTIONS:
  --list
	Get all Atlas projects the authenticated user has access to.
  --project-id [string]
	The Atlas Project ID.
` +
		ac.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ac *AtlasProjectCommand) Synopsis() string {
	return `Access Atlas Organizations.`
}

// Run executes the command
func (ac *AtlasProjectCommand) Run(args []string) int {
	set := ac.NewFlagSet()

	set.BoolVar(&ac.flagList, flagList, false, "")
	set.StringVar(&ac.flagProjectID, flagProjectID, "", "")

	if err := ac.BaseCommand.run(args); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}
	if !ac.flagList || ac.flagProjectID == "" {
		ac.UI.Error("see --help for more information")
		return 1
	}
	if err := ac.run(ac.flagList, ac.flagProjectID); err != nil {
		ac.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (ac *AtlasProjectCommand) run(flagList bool, flagProjectID string) error {

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
	if flagProjectID != "" {
		group, err := client.GroupByID(flagProjectID)
		if err != nil {
			return fmt.Errorf("failed to list Project info: %s", err)
		}
		fmt.Println(group.Name)
		fmt.Println(group.ID)
		return nil
	}

	groups, err := client.Groups()
	if err != nil {
		return fmt.Errorf("failed to list Projects: %s", err)
	}

	for _, group := range groups {
		fmt.Println(group.Name)
		fmt.Println(group.ID)
	}
	return nil
}
