package command

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/derision-test/go-mockgen/internal/genlib/paths"
)

type Options struct {
	ImportPaths      []string
	PkgName          string
	Interfaces       []string
	OutputFilename   string
	OutputDir        string
	OutputImportPath string
	Prefix           string
	Force            bool
}

var GoIdentifierPattern = regexp.MustCompile("^[A-Za-z]([A-Za-z0-9_]*[A-Za-z])?$")

func parseArgs(
	name string,
	description string,
	version string,
	argHook ArgHookFunc,
	argValidator ArgValidatorFunc,
) (*Options, error) {
	app := kingpin.New(name, description).Version(version)

	opts := &Options{
		ImportPaths: []string{},
		Interfaces:  []string{},
	}

	app.Arg("path", "The import paths used to search for eligible interfaces").Required().StringsVar(&opts.ImportPaths)
	app.Flag("package", "The name of the generated package. It will be inferred from the output options by default.").Short('p').StringVar(&opts.PkgName)
	app.Flag("interfaces", "A whitelist of interfaces to generate given the import paths.").Short('i').StringsVar(&opts.Interfaces)
	app.Flag("dirname", "The target output directory. Each mock will be written to a unique file.").Short('d').StringVar(&opts.OutputDir)
	app.Flag("filename", "The target output file. All mocks are written to this file.").Short('o').StringVar(&opts.OutputFilename)
	app.Flag("import-path", "The import path of the generated package. It will be inferred from the tarrget directory by default.").StringVar(&opts.PkgName)
	app.Flag("prefix", "A prefix used in the name of each mock struct. Should be TitleCase by convention.").StringVar(&opts.Prefix)
	app.Flag("force", "Do not abort if a write to disk would overwrite an existing file.").Short('f').BoolVar(&opts.Force)
	argHook(app)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	validators := []ArgValidatorFunc{
		validateOutputPaths,
		validateOptions,
		argValidator,
	}

	for _, f := range validators {
		if fatal, err := f(opts); err != nil {
			if !fatal {
				kingpin.Fatalf("%s, try --help", err.Error())
			}

			return nil, err
		}
	}

	return opts, nil
}

func validateOptions(opts *Options) (bool, error) {
	if opts.PkgName != "" && opts.OutputImportPath != "" {
		return false, fmt.Errorf("package name and output import path are mutually exclusive")
	}

	if opts.OutputImportPath == "" {
		path, ok := paths.InferImportPath(opts.OutputDir)
		if !ok {
			return false, fmt.Errorf("could not infer output import path")
		}

		opts.OutputImportPath = path
	}

	if opts.PkgName == "" {
		opts.PkgName = opts.OutputImportPath[strings.LastIndex(opts.OutputImportPath, "/")+1:]
	}

	if !GoIdentifierPattern.Match([]byte(opts.PkgName)) {
		return false, fmt.Errorf("package name `%s` is illegal", opts.PkgName)
	}

	if opts.Prefix != "" && !GoIdentifierPattern.Match([]byte(opts.Prefix)) {
		return false, fmt.Errorf("prefix `%s` is illegal", opts.Prefix)
	}

	return false, nil
}

func validateOutputPaths(opts *Options) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return true, fmt.Errorf("failed to get current directory")
	}

	if opts.OutputFilename == "" && opts.OutputDir == "" {
		opts.OutputDir = wd
	}

	if opts.OutputFilename != "" && opts.OutputDir != "" {
		return false, fmt.Errorf("dirname and filename are mutually exclusive")
	}

	if opts.OutputFilename != "" {
		opts.OutputDir = path.Dir(opts.OutputFilename)
		opts.OutputFilename = path.Base(opts.OutputFilename)
	}

	if err := paths.EnsureDirExists(opts.OutputDir); err != nil {
		return true, fmt.Errorf(
			"failed to make output directory %s: %s",
			opts.OutputDir,
			err.Error(),
		)
	}

	if opts.OutputDir, err = cleanPath(opts.OutputDir); err != nil {
		return true, err
	}

	return false, nil
}

func cleanPath(path string) (cleaned string, err error) {
	cleaned = path
	for _, f := range []func(string) (string, error){filepath.Abs, filepath.EvalSymlinks} {
		if cleaned, err = f(cleaned); err != nil {
			break
		}
	}

	return
}
