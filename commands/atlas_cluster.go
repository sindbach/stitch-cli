package commands

import (
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

// AtlasClusterCommand is used to export a Stitch App
type AtlasClusterCommand struct {
	*BaseCommand

	flagProjectID  string
	flagAppID      string
	flagOutput     string
	flagAsTemplate bool
}

// Help returns long-form help information for this command
func (ec *AtlasClusterCommand) Help() string {
	return `Export a stitch application to a local directory.

REQUIRED:
  --app-id [string]
	The App ID for your app (i.e. the name of your app followed by a unique suffix, like "my-app-nysja")

OPTIONS:
  --project-id [string]
	Lookup apps associated with this project id, as opposed to ids associated with the current user profile.

  --as-template
	Indicate that the application should be exported as a template.` +
		ec.BaseCommand.Help()
}

// Synopsis returns a one-liner description for this command
func (ec *AtlasClusterCommand) Synopsis() string {
	return `Export a stitch application to a local directory.`
}

// Run executes the command
func (ec *AtlasClusterCommand) Run(args []string) int {
	set := ec.NewFlagSet()

	set.StringVar(&ec.flagOutput, "output", "", "")
	set.StringVar(&ec.flagOutput, "o", "", "")
	set.BoolVar(&ec.flagAsTemplate, "as-template", false, "")

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

func (ec *AtlasClusterCommand) run() error {
	if ec.flagAppID == "" {
		return errAppIDRequired
	}

	user, err := ec.User()
	if err != nil {
		return err
	}

	if !user.LoggedIn() {
		return u.ErrNotLoggedIn
	}

	//stitchClient, err := ec.AtlasClient()
	//if err != nil {
	//	return err
	//}

	return nil
}
