package appcli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/rliebz/tusk/runner"
)

// copyFlags copies all command flags from one cli.App to another.
func copyFlags(target, source *cli.App) {
	for i := range target.Commands {
		targetCommand := &target.Commands[i]
		for j := range source.Commands {
			sourceCommand := source.Commands[j]
			if targetCommand.Name == sourceCommand.Name {
				targetCommand.Flags = sourceCommand.Flags
			}
		}
	}
}

// addAllFlagsUsed adds the top-level flags to tasks where interpolation is used.
func addAllFlagsUsed(cfg *runner.Config, cmd *cli.Command, t *runner.Task) error {
	dependencies, err := runner.FindAllOptions(t, cfg)
	if err != nil {
		return err
	}

	for _, opt := range dependencies {
		if opt.Private {
			continue
		}

		if err := addFlag(cmd, opt); err != nil {
			return errors.Wrapf(
				err,
				`could not add flag "%s" to command "%s"`,
				opt.Name,
				t.Name,
			)
		}
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))
	return nil
}

func addFlag(command *cli.Command, opt *runner.Option) error {
	newFlag, err := createCLIFlag(opt)
	if err != nil {
		return err
	}

	for _, oldFlag := range command.Flags {
		if oldFlag.GetName() == newFlag.GetName() {
			return nil
		}
	}

	command.Flags = append(command.Flags, newFlag)

	return nil
}

// createCLIFlag converts an Option into a cli.Flag.
func createCLIFlag(opt *runner.Option) (cli.Flag, error) {
	name := opt.Name
	if opt.Short != "" {
		name = fmt.Sprintf("%s, %s", name, opt.Short)
	}

	opt.Type = strings.ToLower(opt.Type)
	switch opt.Type {
	case "int", "integer":
		return cli.IntFlag{
			Name:  name,
			Usage: opt.Usage,
		}, nil
	case "float", "float64", "double":
		return cli.Float64Flag{
			Name:  name,
			Usage: opt.Usage,
		}, nil
	case "bool", "boolean":
		return cli.BoolFlag{
			Name:  name,
			Usage: opt.Usage,
		}, nil
	case "string", "":
		return cli.StringFlag{
			Name:  name,
			Usage: opt.Usage,
		}, nil
	default:
		return nil, fmt.Errorf(`unsupported flag type "%s"`, opt.Type)
	}
}
