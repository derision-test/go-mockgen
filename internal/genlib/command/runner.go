package command

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kingpin"

	"github.com/derision-test/go-mockgen/internal/genlib/types"
)

type (
	commandConfig struct {
		argHook      ArgHookFunc
		argValidator ArgValidatorFunc
	}

	Generator        func(ifaces []*types.Interface, opts *Options) error
	ArgHookFunc      func(app *kingpin.Application)
	ArgValidatorFunc func(opts *Options) (bool, error)
)

func Run(
	name string,
	description string,
	version string,
	typeGetter types.TypeGetter,
	generator Generator,
	configs ...ConfigFunc,
) error {
	config := &commandConfig{
		argHook:      func(_ *kingpin.Application) {},
		argValidator: func(_ *Options) (bool, error) { return false, nil },
	}

	for _, f := range configs {
		f(config)
	}

	opts, err := parseArgs(
		name,
		description,
		version,
		config.argHook,
		config.argValidator,
	)

	if err != nil {
		return err
	}

	ifaces, err := Extract(
		typeGetter,
		opts.ImportPaths,
		opts.Interfaces,
	)

	if err != nil {
		return err
	}

	nameMap := map[string]struct{}{}
	for _, t := range ifaces {
		nameMap[strings.ToLower(t.Name)] = struct{}{}
	}

	for _, name := range opts.Interfaces {
		if _, ok := nameMap[strings.ToLower(name)]; !ok {
			return fmt.Errorf("type '%s' not found in supplied import paths", name)
		}
	}

	return generator(ifaces, opts)
}
