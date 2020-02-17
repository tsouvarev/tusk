package appcli

import (
	"fmt"
	"sort"

	"github.com/urfave/cli"

	"github.com/rliebz/tusk/runner"
)

// addTasks adds a series of tasks to a cli.App using a command creator.
func addTasks(app *cli.App, cfg *runner.Config, create commandCreator) error {
	for _, t := range cfg.Tasks {
		if err := addTask(app, cfg, t, create); err != nil {
			return fmt.Errorf("could not add task %q: %w", t.Name, err)
		}
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	return nil
}

func addTask(app *cli.App, cfg *runner.Config, t *runner.Task, create commandCreator) error {
	if t.Private {
		return nil
	}

	command, err := create(app, t)
	if err != nil {
		return fmt.Errorf(`could not create command %q: %w`, t.Name, err)
	}

	if err := addAllFlagsUsed(cfg, command, t); err != nil {
		return fmt.Errorf("could not add flags: %w", err)
	}

	app.Commands = append(app.Commands, *command)

	return nil
}
