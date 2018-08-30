package commands

import (
	"github.com/mitchellh/cli"
)

// NewAtlasCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewAtlasCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &AtlasCommand{
			BaseCommand: &BaseCommand{
				Name: "atlas",
				UI:   ui,
			},
		}, nil
	}
}

// AtlasCommand is a group category
type AtlasCommand struct {
	*BaseCommand
}

// Help returns long-form help information for this command
func (ec *AtlasCommand) Help() string {
	return `Command group for MongoDB Atlas.`
}

// Synopsis returns a one-liner description for this command
func (ec *AtlasCommand) Synopsis() string {
	return `Command group for MongoDB Atlas.`
}

// Run executes the command
func (ec *AtlasCommand) Run(args []string) int {

	if err := ec.BaseCommand.run(args); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}
	return 0
}
