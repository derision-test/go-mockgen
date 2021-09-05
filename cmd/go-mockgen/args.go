package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/derision-test/go-mockgen/internal/mockgen/consts"
	"github.com/derision-test/go-mockgen/internal/mockgen/generation"
	"github.com/derision-test/go-mockgen/internal/mockgen/paths"
)

func parseArgs() (*generation.Options, error) {
	opts := &generation.Options{
		ImportPaths: []string{},
		Interfaces:  []string{},
	}

	app := kingpin.New(consts.Name, consts.Description).Version(consts.Version)
	app.UsageWriter(os.Stdout)

	app.Arg("path", "The import paths used to search for eligible interfaces").Required().StringsVar(&opts.ImportPaths)
	app.Flag("package", "The name of the generated package. It will be inferred from the output options by default.").Short('p').StringVar(&opts.PkgName)
	app.Flag("interfaces", "A list of target interfaces to generate defined in the given the import paths.").Short('i').StringsVar(&opts.Interfaces)
	app.Flag("exclude", "A list of interfaces to exclude from generation. Mocks for all other exported interfaces defined in the given import paths are generated.").Short('e').StringsVar(&opts.Exclude)
	app.Flag("dirname", "The target output directory. Each mock will be written to a unique file.").Short('d').StringVar(&opts.OutputDir)
	app.Flag("filename", "The target output file. All mocks are written to this file.").Short('o').StringVar(&opts.OutputFilename)
	app.Flag("import-path", "The import path of the generated package. It will be inferred from the target directory by default.").StringVar(&opts.PkgName)
	app.Flag("prefix", "A prefix used in the name of each mock struct. Should be TitleCase by convention.").StringVar(&opts.Prefix)
	app.Flag("force", "Do not abort if a write to disk would overwrite an existing file.").Short('f').BoolVar(&opts.Force)
	app.Flag("disable-formatting", "Do not run goimports over the rendered files.").BoolVar(&opts.DisableFormatting)
	app.Flag("goimports", "Path to the goimports binary.").Default("goimports").StringVar(&opts.GoImportsBinary)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	validators := []func(opts *generation.Options) (bool, error){
		validateOutputPaths,
		validateOptions,
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

func validateOutputPaths(opts *generation.Options) (bool, error) {
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

var goIdentifierPattern = regexp.MustCompile("^[A-Za-z]([A-Za-z0-9_]*[A-Za-z])?$")

func validateOptions(opts *generation.Options) (bool, error) {
	if opts.PkgName != "" && opts.OutputImportPath != "" {
		return false, fmt.Errorf("package name and output import path are mutually exclusive")
	}

	if len(opts.Interfaces) != 0 && len(opts.Exclude) != 0 {
		return false, fmt.Errorf("interface lists and exclude lists are mutually exclusive")
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

	if !goIdentifierPattern.Match([]byte(opts.PkgName)) {
		return false, fmt.Errorf("package name `%s` is illegal", opts.PkgName)
	}

	if opts.Prefix != "" && !goIdentifierPattern.Match([]byte(opts.Prefix)) {
		return false, fmt.Errorf("prefix `%s` is illegal", opts.Prefix)
	}

	return false, nil
}

func cleanPath(path string) (cleaned string, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return "", err
	}

	if path, err = filepath.EvalSymlinks(path); err != nil {
		return "", err
	}

	return path, nil
}
