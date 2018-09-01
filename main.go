// Stitch is a tool for command-line administration of MongoDB Stitch applications.
package main

import (
	"os"
	"path/filepath"

	"github.com/sindbach/stitch-cli/commands"
	"github.com/sindbach/stitch-cli/utils"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI(filepath.Base(os.Args[0]), utils.CLIVersion)
	c.Args = os.Args[1:]

	var ui cli.Ui = &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c.Commands = map[string]cli.CommandFactory{
		"whoami":        commands.NewWhoamiCommandFactory(ui),
		"login":         commands.NewLoginCommandFactory(ui),
		"logout":        commands.NewLogoutCommandFactory(ui),
		"stitch":        commands.NewStitchCommandFactory(ui),
		"stitch export": commands.NewExportCommandFactory(ui),
		"stitch import": commands.NewImportCommandFactory(ui),
		"atlas":         commands.NewAtlasCommandFactory(ui),
		"atlas cluster": commands.NewAtlasClusterCommandFactory(ui),
		"atlas org":     commands.NewAtlasOrgCommandFactory(ui),
		"atlas project": commands.NewAtlasProjectCommandFactory(ui),
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
