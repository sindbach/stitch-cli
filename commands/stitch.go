package commands

import (
	"github.com/mitchellh/cli"
)

// NewStitchCommandFactory returns a new cli.CommandFactory given a cli.Ui
func NewStitchCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &StitchCommand{
			BaseCommand: &BaseCommand{
				Name: "stitch",
				UI:   ui,
			},
		}, nil
	}
}

// StitchCommand is a group category
type StitchCommand struct {
	*BaseCommand
}

// Help returns long-form help information for this command
func (ec *StitchCommand) Help() string {
	return `Command group for MongoDB Stitch.`
}

// Synopsis returns a one-liner description for this command
func (ec *StitchCommand) Synopsis() string {
	return `Command group for MongoDB Stitch.`
}

// Run executes the command
func (ec *StitchCommand) Run(args []string) int {

	if err := ec.BaseCommand.run(args); err != nil {
		ec.UI.Error(err.Error())
		return 1
	}
	return 0
}
